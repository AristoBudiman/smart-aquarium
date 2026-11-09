package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pastikan Store mengimplementasikan core.Store.
var _ core.Store = (*Store)(nil)

// Store mengelola operasi basis data PostgreSQL menggunakan connection pool pgxpool.
type Store struct {
	pool *pgxpool.Pool
}

// New membuat instance Store baru dari *pgxpool.Pool yang sudah ada.
func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

// Connect menginisialisasi connection pool dan mencoba menyambung ulang ke PostgreSQL.
// Koneksi dicoba ulang dengan jeda karena PostgreSQL dapat siap lebih lambat saat startup.
func Connect(ctx context.Context, databaseURL string) (*Store, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("db: format database URL tidak valid: %w", err)
	}

	var pool *pgxpool.Pool
	var lastErr error

	const maxAttempts = 15
	const retryDelay = 2 * time.Second

	for i := 1; i <= maxAttempts; i++ {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("db: koneksi dibatalkan: %w", ctx.Err())
		default:
		}

		pool, err = pgxpool.NewWithConfig(ctx, cfg)
		if err == nil {
			pingCtx, pingCancel := context.WithTimeout(ctx, 3*time.Second)
			pingErr := pool.Ping(pingCtx)
			pingCancel()

			if pingErr == nil {
				return New(pool), nil
			}

			lastErr = pingErr
			pool.Close()
		} else {
			lastErr = err
		}

		if i < maxAttempts {
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("db: koneksi dibatalkan: %w", ctx.Err())
			case <-time.After(retryDelay):
			}
		}
	}

	return nil, fmt.Errorf("db: gagal terhubung ke PostgreSQL setelah %d percobaan: %w", maxAttempts, lastErr)
}

// Close menutup semua koneksi di dalam pool.
func (s *Store) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

// Pool mengembalikan underlying *pgxpool.Pool.
func (s *Store) Pool() *pgxpool.Pool {
	return s.pool
}

// EnsureNode mendaftarkan node baru jika belum ada dengan nama default dan mode manual.
// Mengembalikan created = true jika baris baru berhasil dibuat.
func (s *Store) EnsureNode(ctx context.Context, nodeID string) (bool, error) {
	const query = `
		INSERT INTO nodes (node_id, name, mode)
		VALUES ($1, $2, 'manual')
		ON CONFLICT (node_id) DO NOTHING;
	`
	defaultName := "Akuarium " + nodeID
	cmdTag, err := s.pool.Exec(ctx, query, nodeID, defaultName)
	if err != nil {
		return false, fmt.Errorf("db EnsureNode (%s): %w", nodeID, err)
	}
	return cmdTag.RowsAffected() > 0, nil
}

// Snapshot mengambil kondisi node saat ini beserta status prediksi terakhir.
func (s *Store) Snapshot(ctx context.Context, nodeID string) (core.NodeSnapshot, error) {
	const query = `
		SELECT n.mode, n.online, n.heater, n.fan, n.pump, n.failsafe,
		       COALESCE(p.status, 'NORMAL') AS last_status
		FROM nodes n
		LEFT JOIN LATERAL (
			SELECT status
			FROM predictions
			WHERE node_id = n.node_id
			ORDER BY ts DESC
			LIMIT 1
		) p ON true
		WHERE n.node_id = $1;
	`

	var modeStr string
	var online bool
	var heater, fan, pump *int16
	var failsafe bool
	var lastStatusStr string

	err := s.pool.QueryRow(ctx, query, nodeID).Scan(
		&modeStr,
		&online,
		&heater,
		&fan,
		&pump,
		&failsafe,
		&lastStatusStr,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.NodeSnapshot{}, fmt.Errorf("db Snapshot: node %s tidak ditemukan: %w", nodeID, err)
		}
		return core.NodeSnapshot{}, fmt.Errorf("db Snapshot (%s): %w", nodeID, err)
	}

	hasState := heater != nil && fan != nil && pump != nil
	var state core.ActuatorState
	if hasState {
		state.Heater = int(*heater)
		state.Fan = int(*fan)
		state.Pump = int(*pump)
		state.Failsafe = failsafe
	}

	return core.NodeSnapshot{
		NodeID:     nodeID,
		Mode:       core.Mode(modeStr),
		Online:     online,
		HasState:   hasState,
		State:      state,
		LastStatus: core.Class(lastStatusStr),
	}, nil
}

// SaveTelemetry menyimpan pembacaan sensor telemetry dan memperbarui last_seen pada nodes.
func (s *Store) SaveTelemetry(ctx context.Context, t core.Telemetry) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("db SaveTelemetry begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	const insertTelemetry = `
		INSERT INTO telemetry (node_id, ts, temp, ntu)
		VALUES ($1, $2, $3, $4);
	`
	if _, err := tx.Exec(ctx, insertTelemetry, t.NodeID, t.TS, t.Temp, t.NTU); err != nil {
		return fmt.Errorf("db SaveTelemetry insert (%s): %w", t.NodeID, err)
	}

	const updateLastSeen = `
		UPDATE nodes
		SET last_seen = $1
		WHERE node_id = $2;
	`
	if _, err := tx.Exec(ctx, updateLastSeen, t.TS, t.NodeID); err != nil {
		return fmt.Errorf("db SaveTelemetry update last_seen (%s): %w", t.NodeID, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("db SaveTelemetry commit tx: %w", err)
	}
	return nil
}

// SavePrediction menyimpan hasil klasifikasi model ke tabel predictions.
func (s *Store) SavePrediction(ctx context.Context, nodeID string, ts time.Time, c core.Classes) error {
	const query = `
		INSERT INTO predictions (node_id, ts, class_temp, class_turbidity, class_overall, status)
		VALUES ($1, $2, $3, $4, $5, $6);
	`
	_, err := s.pool.Exec(
		ctx,
		query,
		nodeID,
		ts,
		string(c.Temp),
		string(c.Turbidity),
		string(c.Overall),
		string(c.Status),
	)
	if err != nil {
		return fmt.Errorf("db SavePrediction (%s): %w", nodeID, err)
	}
	return nil
}

// SaveState menyimpan status aktuator aktual dan status failsafe ke tabel nodes.
func (s *Store) SaveState(ctx context.Context, nodeID string, st core.ActuatorState) error {
	const query = `
		UPDATE nodes
		SET heater = $1, fan = $2, pump = $3, failsafe = $4
		WHERE node_id = $5;
	`
	_, err := s.pool.Exec(ctx, query, st.Heater, st.Fan, st.Pump, st.Failsafe, nodeID)
	if err != nil {
		return fmt.Errorf("db SaveState (%s): %w", nodeID, err)
	}
	return nil
}

// SetOnline memperbarui status online dan timestamp last_seen node.
func (s *Store) SetOnline(ctx context.Context, nodeID string, online bool) error {
	const query = `
		UPDATE nodes
		SET online = $1,
		    last_seen = CASE WHEN $1 THEN now() ELSE last_seen END
		WHERE node_id = $2;
	`
	_, err := s.pool.Exec(ctx, query, online, nodeID)
	if err != nil {
		return fmt.Errorf("db SetOnline (%s): %w", nodeID, err)
	}
	return nil
}

// LogActuators mencatat riwayat perubahan aktuator ke tabel actuator_logs.
func (s *Store) LogActuators(ctx context.Context, nodeID string, a core.Actuators, source core.Source) error {
	const query = `
		INSERT INTO actuator_logs (node_id, ts, heater, fan, pump, source)
		VALUES ($1, now(), $2, $3, $4, $5);
	`
	_, err := s.pool.Exec(ctx, query, nodeID, a.Heater, a.Fan, a.Pump, string(source))
	if err != nil {
		return fmt.Errorf("db LogActuators (%s): %w", nodeID, err)
	}
	return nil
}

// SaveAlert mencatat alert baru ke tabel alerts dan mengembalikan ID baris yang dibuat.
func (s *Store) SaveAlert(ctx context.Context, a core.Alert) (int64, error) {
	const query = `
		INSERT INTO alerts (node_id, ts, level, message, sent)
		VALUES ($1, $2, $3, $4, FALSE)
		RETURNING id;
	`
	ts := a.TS
	if ts.IsZero() {
		ts = time.Now().UTC()
	}

	var id int64
	err := s.pool.QueryRow(ctx, query, a.NodeID, ts, string(a.Level), a.Message).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("db SaveAlert (%s): %w", a.NodeID, err)
	}
	return id, nil
}

// MarkAlertSent menandai bahwa notifikasi alert telah berhasil terkirim.
func (s *Store) MarkAlertSent(ctx context.Context, id int64) error {
	const query = `
		UPDATE alerts
		SET sent = TRUE
		WHERE id = $1;
	`
	_, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("db MarkAlertSent (%d): %w", id, err)
	}
	return nil
}

// TelegramChatIDs mengambil daftar telegram_chat_id pemilik node dari tabel nodes dan users.
// Mengembalikan slice kosong jika node belum diklaim atau pemilik belum mengisi chat ID.
func (s *Store) TelegramChatIDs(ctx context.Context, nodeID string) ([]string, error) {
	const query = `
		SELECT u.telegram_chat_id
		FROM nodes n
		JOIN users u ON u.id = n.owner_id
		WHERE n.node_id = $1
		  AND u.telegram_chat_id IS NOT NULL
		  AND u.telegram_chat_id != '';
	`
	rows, err := s.pool.Query(ctx, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("db TelegramChatIDs query (%s): %w", nodeID, err)
	}
	defer rows.Close()

	ids := make([]string, 0)
	for rows.Next() {
		var chatID string
		if err := rows.Scan(&chatID); err != nil {
			return nil, fmt.Errorf("db TelegramChatIDs scan (%s): %w", nodeID, err)
		}
		ids = append(ids, chatID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db TelegramChatIDs rows (%s): %w", nodeID, err)
	}

	return ids, nil
}
