INSERT INTO teams(name)
SELECT 'team_' || LPAD(i::text, 2, '0')
FROM generate_series(1, 20) AS s(i)
ON CONFLICT DO NOTHING;

INSERT INTO users(id, name, team_name, is_active)
SELECT
    'u' || LPAD(i::text, 3, '0') AS id,
    'User ' || LPAD(i::text, 3, '0') AS name,
    'team_' || LPAD(((i - 1) % 20 + 1)::text, 2, '0') AS team_name,
    (random() < 0.8) AS is_active
FROM generate_series(1, 200) AS s(i)
ON CONFLICT DO NOTHING;
