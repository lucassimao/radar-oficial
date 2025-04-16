ALTER TABLE diarios
ADD CONSTRAINT diarios_institution_description_key UNIQUE (institution_id, description);
