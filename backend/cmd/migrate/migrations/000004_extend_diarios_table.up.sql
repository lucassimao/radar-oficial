ALTER TABLE diarios
    ADD COLUMN description TEXT,
    ADD COLUMN last_modified_at TIMESTAMP WITHOUT TIME ZONE,
    ALTER COLUMN published_at TYPE TIMESTAMP WITHOUT TIME ZONE,
    ALTER COLUMN published_at DROP NOT NULL;


