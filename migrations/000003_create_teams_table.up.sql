CREATE TABLE IF NOT EXISTS teams (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    members integer[],
)