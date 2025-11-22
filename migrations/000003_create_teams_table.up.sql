CREATE TABLE IF NOT EXISTS teams (
    id bigserial PRIMARY KEY,
    name text UNIQUE NOT NULL,
    members integer[]
)