CREATE OR REPLACE FUNCTION set_merged_at()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'merged' AND OLD.status IS DISTINCT FROM NEW.status THEN
        NEW.merged_at = NOW();
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_set_merged_at
BEFORE UPDATE OF status ON pull_requests
FOR EACH ROW
EXECUTE FUNCTION set_merged_at();