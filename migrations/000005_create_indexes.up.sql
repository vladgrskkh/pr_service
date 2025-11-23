CREATE INDEX IF NOT EXISTS idx_users_team_name
ON users(team_name);

CREATE INDEX IF NOT EXISTS idx_pull_requests_author_id
ON pull_requests(author_id);
