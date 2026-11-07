package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNodeNotFound = errors.New("node tidak ditemukan")
	ErrUserNotFound = errors.New("user tidak ditemukan")
)

type User struct {
	ID             int64
	Username       string
	PasswordHash   string
	TelegramChatID *string
}

type NodeSummary struct {
	NodeID   string      `json:"node_id"`
	Name     string      `json:"name"`
	Online   bool        `json:"online"`
	Mode     core.Mode   `json:"mode"`
	LastSeen *time.Time  `json:"last_seen"`
	Temp     *float64    `json:"temp"`
	NTU      *float64    `json:"ntu"`
	Status   *core.Class `json:"status"`
}

type NodeDetail struct {
	NodeID         string      `json:"node_id"`
	Name           string      `json:"name"`
	Online         bool        `json:"online"`
	Mode           core.Mode   `json:"mode"`
	LastSeen       *time.Time  `json:"last_seen"`
	Temp           *float64    `json:"temp"`
	NTU            *float64    `json:"ntu"`
	Status         *core.Class `json:"status"`
	ClassTemp      *core.Class `json:"class_temp"`
	ClassTurbidity *core.Class `json:"class_turbidity"`
	ClassOverall   *core.Class `json:"class_overall"`
	Heater         *int        `json:"heater"`
	Fan            *int        `json:"fan"`
	Pump           *int        `json:"pump"`
	Failsafe       bool        `json:"failsafe"`
}

type TelemetryPoint struct {
	TS     time.Time  `json:"ts"`
	Temp   float64    `json:"temp"`
	NTU    float64    `json:"ntu"`
	Status core.Class `json:"status"`
}

type ActuatorLog struct {
	ID     int64       `json:"id"`
	TS     time.Time   `json:"ts"`
	Heater int         `json:"heater"`
	Fan    int         `json:"fan"`
	Pump   int         `json:"pump"`
	Source core.Source `json:"source"`
}

type AlertLog struct {
	ID      int64      `json:"id"`
	TS      time.Time  `json:"ts"`
	Level   core.Class `json:"level"`
	Message string     `json:"message"`
}

// EnsureAdminUser membuat akun awal jika belum ada.
func (s *Store) EnsureAdminUser(ctx context.Context, username, passwordHash string) error {
	const query = `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		ON CONFLICT (username) DO NOTHING;
	`
	_, err := s.pool.Exec(ctx, query, username, passwordHash)
	if err != nil {
		return fmt.Errorf("db: EnsureAdminUser gagal: %w", err)
	}
	return nil
}

// GetUserByUsername mengambil akun user berdasarkan username.
func (s *Store) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	const query = `
		SELECT id, username, password_hash, telegram_chat_id
		FROM users
		WHERE username = $1;
	`
	var u User
	err := s.pool.QueryRow(ctx, query, username).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.TelegramChatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("db: GetUserByUsername gagal: %w", err)
	}
	return &u, nil
}

// UpdateUserTelegramChatID memperbarui telegram_chat_id milik pengguna.
func (s *Store) UpdateUserTelegramChatID(ctx context.Context, userID int64, chatID string) error {
	const query = `
		UPDATE users
		SET telegram_chat_id = $1
		WHERE id = $2;
	`
	var param any = chatID
	if chatID == "" {
		param = nil
	}
	cmd, err := s.pool.Exec(ctx, query, param, userID)
	if err != nil {
		return fmt.Errorf("db: UpdateUserTelegramChatID gagal: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

// ListNodeSummaries mengambil seluruh daftar node beserta data sensor dan status terbaru.
func (s *Store) ListNodeSummaries(ctx context.Context) ([]NodeSummary, error) {
	const query = `
		SELECT n.node_id, n.name, n.online, n.mode, n.last_seen,
		       t.temp, t.ntu,
		       p.status
		FROM nodes n
		LEFT JOIN LATERAL (
			SELECT temp, ntu
			FROM telemetry
			WHERE node_id = n.node_id
			ORDER BY ts DESC
			LIMIT 1
		) t ON true
		LEFT JOIN LATERAL (
			SELECT status
			FROM predictions
			WHERE node_id = n.node_id
			ORDER BY ts DESC
			LIMIT 1
		) p ON true
		ORDER BY n.node_id ASC;
	`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("db: ListNodeSummaries query gagal: %w", err)
	}
	defer rows.Close()

	list := make([]NodeSummary, 0)
	for rows.Next() {
		var item NodeSummary
		var modeStr string
		var lastSeen *time.Time
		var statusStr *string

		if err := rows.Scan(
			&item.NodeID,
			&item.Name,
			&item.Online,
			&modeStr,
			&lastSeen,
			&item.Temp,
			&item.NTU,
			&statusStr,
		); err != nil {
			return nil, fmt.Errorf("db: ListNodeSummaries scan gagal: %w", err)
		}

		item.Mode = core.Mode(modeStr)
		if lastSeen != nil {
			utc := lastSeen.UTC()
			item.LastSeen = &utc
		}
		if statusStr != nil {
			st := core.Class(*statusStr)
			item.Status = &st
		}

		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: ListNodeSummaries rows error: %w", err)
	}

	return list, nil
}

// GetNodeDetail mengambil informasi detail sebuah node beserta telemetry dan status prediksi terakhir.
func (s *Store) GetNodeDetail(ctx context.Context, nodeID string) (*NodeDetail, error) {
	const query = `
		SELECT n.node_id, n.name, n.online, n.mode, n.last_seen,
		       n.heater, n.fan, n.pump, n.failsafe,
		       t.temp, t.ntu,
		       p.class_temp, p.class_turbidity, p.class_overall, p.status
		FROM nodes n
		LEFT JOIN LATERAL (
			SELECT temp, ntu
			FROM telemetry
			WHERE node_id = n.node_id
			ORDER BY ts DESC
			LIMIT 1
		) t ON true
		LEFT JOIN LATERAL (
			SELECT class_temp, class_turbidity, class_overall, status
			FROM predictions
			WHERE node_id = n.node_id
			ORDER BY ts DESC
			LIMIT 1
		) p ON true
		WHERE n.node_id = $1;
	`

	var detail NodeDetail
	var modeStr string
	var lastSeen *time.Time
	var heater, fan, pump *int16
	var classTemp, classTurb, classOver, status *string

	err := s.pool.QueryRow(ctx, query, nodeID).Scan(
		&detail.NodeID,
		&detail.Name,
		&detail.Online,
		&modeStr,
		&lastSeen,
		&heater,
		&fan,
		&pump,
		&detail.Failsafe,
		&detail.Temp,
		&detail.NTU,
		&classTemp,
		&classTurb,
		&classOver,
		&status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNodeNotFound
		}
		return nil, fmt.Errorf("db: GetNodeDetail (%s) gagal: %w", nodeID, err)
	}

	detail.Mode = core.Mode(modeStr)
	if lastSeen != nil {
		utc := lastSeen.UTC()
		detail.LastSeen = &utc
	}
	if heater != nil {
		h := int(*heater)
		detail.Heater = &h
	}
	if fan != nil {
		f := int(*fan)
		detail.Fan = &f
	}
	if pump != nil {
		p := int(*pump)
		detail.Pump = &p
	}
	if classTemp != nil {
		ct := core.Class(*classTemp)
		detail.ClassTemp = &ct
	}
	if classTurb != nil {
		ct := core.Class(*classTurb)
		detail.ClassTurbidity = &ct
	}
	if classOver != nil {
		co := core.Class(*classOver)
		detail.ClassOverall = &co
	}
	if status != nil {
		st := core.Class(*status)
		detail.Status = &st
	}

	return &detail, nil
}

// UpdateNodeName memperbarui nama akuarium dan mengembalikan detail terbarunya.
func (s *Store) UpdateNodeName(ctx context.Context, nodeID, name string) (*NodeDetail, error) {
	const query = `
		UPDATE nodes
		SET name = $1
		WHERE node_id = $2;
	`
	cmd, err := s.pool.Exec(ctx, query, name, nodeID)
	if err != nil {
		return nil, fmt.Errorf("db: UpdateNodeName gagal: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return nil, ErrNodeNotFound
	}

	return s.GetNodeDetail(ctx, nodeID)
}

// UpdateNodeMode memperbarui mode pengoperasian akuarium ('auto'/'manual') dan mengembalikan detailnya.
func (s *Store) UpdateNodeMode(ctx context.Context, nodeID string, mode core.Mode) (*NodeDetail, error) {
	const query = `
		UPDATE nodes
		SET mode = $1
		WHERE node_id = $2;
	`
	cmd, err := s.pool.Exec(ctx, query, string(mode), nodeID)
	if err != nil {
		return nil, fmt.Errorf("db: UpdateNodeMode gagal: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return nil, ErrNodeNotFound
	}

	return s.GetNodeDetail(ctx, nodeID)
}

// GetTelemetryHistory mengembalikan histori telemetry. Jika rentang <= 48 jam mengembalikan data mentah,
// jika rentang > 48 jam mengembalikan snapshot rata-rata per jam menggunakan date_bin.
func (s *Store) GetTelemetryHistory(ctx context.Context, nodeID string, from, to time.Time) ([]TelemetryPoint, error) {
	delta := to.Sub(from)

	var query string
	if delta <= 48*time.Hour {
		query = `
			SELECT t.ts, t.temp, t.ntu, COALESCE(p.status, 'NORMAL') AS status
			FROM telemetry t
			LEFT JOIN predictions p ON p.node_id = t.node_id AND p.ts = t.ts
			WHERE t.node_id = $1 AND t.ts BETWEEN $2 AND $3
			ORDER BY t.ts ASC;
		`
	} else {
		query = `
			SELECT date_bin('1 hour', t.ts, TIMESTAMPTZ '2000-01-01') AS ts,
			       avg(t.temp) AS temp,
			       avg(t.ntu) AS ntu,
			       (ARRAY['NORMAL', 'CAUTION', 'WARNING', 'SEVERE'])[max(
			           CASE p.status WHEN 'SEVERE' THEN 4 WHEN 'WARNING' THEN 3 WHEN 'CAUTION' THEN 2 ELSE 1 END
			       )] AS status
			FROM telemetry t
			LEFT JOIN predictions p ON p.node_id = t.node_id AND p.ts = t.ts
			WHERE t.node_id = $1 AND t.ts BETWEEN $2 AND $3
			GROUP BY 1
			ORDER BY 1 ASC;
		`
	}

	rows, err := s.pool.Query(ctx, query, nodeID, from, to)
	if err != nil {
		return nil, fmt.Errorf("db: GetTelemetryHistory query gagal: %w", err)
	}
	defer rows.Close()

	points := make([]TelemetryPoint, 0)
	for rows.Next() {
		var pt TelemetryPoint
		var statusStr string
		if err := rows.Scan(&pt.TS, &pt.Temp, &pt.NTU, &statusStr); err != nil {
			return nil, fmt.Errorf("db: GetTelemetryHistory scan gagal: %w", err)
		}
		pt.TS = pt.TS.UTC()
		pt.Status = core.Class(statusStr)
		points = append(points, pt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: GetTelemetryHistory rows error: %w", err)
	}

	return points, nil
}

// GetNodeLogs mengambil log perubahan aktuator dan alert terbaru untuk suatu node.
func (s *Store) GetNodeLogs(ctx context.Context, nodeID string) ([]ActuatorLog, []AlertLog, error) {
	const queryActuators = `
		SELECT id, ts, heater, fan, pump, source
		FROM actuator_logs
		WHERE node_id = $1
		ORDER BY ts DESC
		LIMIT 50;
	`
	rowsAct, err := s.pool.Query(ctx, queryActuators, nodeID)
	if err != nil {
		return nil, nil, fmt.Errorf("db: GetNodeLogs actuator_logs query gagal: %w", err)
	}
	defer rowsAct.Close()

	actuators := make([]ActuatorLog, 0)
	for rowsAct.Next() {
		var item ActuatorLog
		var src string
		if err := rowsAct.Scan(&item.ID, &item.TS, &item.Heater, &item.Fan, &item.Pump, &src); err != nil {
			return nil, nil, fmt.Errorf("db: GetNodeLogs actuator_logs scan gagal: %w", err)
		}
		item.TS = item.TS.UTC()
		item.Source = core.Source(src)
		actuators = append(actuators, item)
	}
	if err := rowsAct.Err(); err != nil {
		return nil, nil, fmt.Errorf("db: GetNodeLogs actuator_logs rows error: %w", err)
	}

	const queryAlerts = `
		SELECT id, ts, level, message
		FROM alerts
		WHERE node_id = $1
		ORDER BY ts DESC
		LIMIT 50;
	`
	rowsAlert, err := s.pool.Query(ctx, queryAlerts, nodeID)
	if err != nil {
		return nil, nil, fmt.Errorf("db: GetNodeLogs alerts query gagal: %w", err)
	}
	defer rowsAlert.Close()

	alerts := make([]AlertLog, 0)
	for rowsAlert.Next() {
		var item AlertLog
		var lvl string
		if err := rowsAlert.Scan(&item.ID, &item.TS, &lvl, &item.Message); err != nil {
			return nil, nil, fmt.Errorf("db: GetNodeLogs alerts scan gagal: %w", err)
		}
		item.TS = item.TS.UTC()
		item.Level = core.Class(lvl)
		alerts = append(alerts, item)
	}
	if err := rowsAlert.Err(); err != nil {
		return nil, nil, fmt.Errorf("db: GetNodeLogs alerts rows error: %w", err)
	}

	return actuators, alerts, nil
}
