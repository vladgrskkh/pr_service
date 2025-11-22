CREATE TABLE IF NOT EXISTS users (
    id text PRIMARY KEY,
    name text NOT NULL,
    team_name text NOT NULL REFERENCES teams(name) ON DELETE SET NULL,
    is_active bool NOT NULL
);