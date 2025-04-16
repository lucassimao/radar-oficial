ALTER TABLE diarios
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS last_modified_at,
    ALTER COLUMN published_at TYPE DATE,
    ALTER COLUMN published_at SET NOT NULL;
