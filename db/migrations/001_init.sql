CREATE TABLE users (
    id               BIGSERIAL PRIMARY KEY,
    username         TEXT NOT NULL UNIQUE,
    password_hash    TEXT NOT NULL,
    telegram_chat_id TEXT
);

CREATE TABLE nodes (
    node_id    TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    mode       TEXT NOT NULL DEFAULT 'manual' CHECK (mode IN ('auto', 'manual')),
    online     BOOLEAN NOT NULL DEFAULT FALSE,
    last_seen  TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    heater     SMALLINT CHECK (heater IN (0, 1)),
    fan        SMALLINT CHECK (fan IN (0, 1)),
    pump       SMALLINT CHECK (pump IN (0, 1)),
    failsafe   BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE telemetry (
    id      BIGSERIAL PRIMARY KEY,
    node_id TEXT NOT NULL REFERENCES nodes (node_id) ON DELETE CASCADE,
    ts      TIMESTAMPTZ NOT NULL,
    temp    DOUBLE PRECISION NOT NULL,
    ntu     DOUBLE PRECISION NOT NULL
);
CREATE INDEX telemetry_node_ts ON telemetry (node_id, ts DESC);

CREATE TABLE predictions (
    id              BIGSERIAL PRIMARY KEY,
    node_id         TEXT NOT NULL REFERENCES nodes (node_id) ON DELETE CASCADE,
    ts              TIMESTAMPTZ NOT NULL,
    class_temp      TEXT NOT NULL,
    class_turbidity TEXT NOT NULL,
    class_overall   TEXT NOT NULL,
    status          TEXT NOT NULL
);
CREATE INDEX predictions_node_ts ON predictions (node_id, ts DESC);

CREATE TABLE actuator_logs (
    id      BIGSERIAL PRIMARY KEY,
    node_id TEXT NOT NULL REFERENCES nodes (node_id) ON DELETE CASCADE,
    ts      TIMESTAMPTZ NOT NULL DEFAULT now(),
    heater  SMALLINT NOT NULL,
    fan     SMALLINT NOT NULL,
    pump    SMALLINT NOT NULL,
    source  TEXT NOT NULL CHECK (source IN ('auto', 'manual', 'failsafe'))
);
CREATE INDEX actuator_logs_node_ts ON actuator_logs (node_id, ts DESC);

CREATE TABLE alerts (
    id      BIGSERIAL PRIMARY KEY,
    node_id TEXT NOT NULL REFERENCES nodes (node_id) ON DELETE CASCADE,
    ts      TIMESTAMPTZ NOT NULL DEFAULT now(),
    level   TEXT NOT NULL,
    message TEXT NOT NULL,
    sent    BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX alerts_node_ts ON alerts (node_id, ts DESC);
