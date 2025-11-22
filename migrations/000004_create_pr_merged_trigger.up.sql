CREATE OR REPLACE FUNCTION set_merged_at()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'MERGED' AND OLD.status IS DISTINCT FROM NEW.status THEN
        NEW.merged_at = NOW();
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_trigger
        WHERE tgname = 'trg_set_merged_at'
    ) THEN
        EXECUTE $sql$
            CREATE TRIGGER trg_set_merged_at
            BEFORE UPDATE OF status ON pull_requests
            FOR EACH ROW
            EXECUTE FUNCTION set_merged_at();
        $sql$;
    END IF;
END;
$$;