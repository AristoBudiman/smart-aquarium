ALTER TABLE users ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE nodes ADD COLUMN owner_id BIGINT REFERENCES users (id) ON DELETE SET NULL;
ALTER TABLE nodes ADD COLUMN claim_code TEXT NOT NULL DEFAULT upper(substr(md5(random()::text || clock_timestamp()::text), 1, 8));
CREATE INDEX nodes_owner ON nodes (owner_id);
