DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'pull_req_status'
    ) THEN
        EXECUTE $sql$
            CREATE TYPE pull_req_status AS ENUM ('OPEN', 'MERGED');
        $sql$;
    END IF;
END;
$$;

CREATE TABLE IF NOT EXISTS pull_requests (
    id text PRIMARY KEY,
    name text NOT NULL,
    author_id text NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status pull_req_status NOT NULL,
    assigned_reviewers text[], -- 2 reviewers max
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    merged_at timestamp(0) with time zone DEFAULT NULL
);
