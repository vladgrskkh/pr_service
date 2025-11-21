CREATE TABLE IF NOT EXISTS pull_requests (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    author_id integer NOT NULL,
    status text NOT NULL,
    assigned_reviewers integer[]
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    merged_at timestamp(0) with time zone DEFAULT NULL,
);
