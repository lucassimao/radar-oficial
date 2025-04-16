INSERT INTO institutions (
    id,
    name,
    slug,
    type,
    estate,
    city,
    source_url,
    active,
    created_at,
    updated_at
)
VALUES (
    DEFAULT,
    'Governo do Estado do Piauí',
    'governo-pi',
    'state',
    'PI',
    NULL,
    'https://www.diario.pi.gov.br/doe/',
    TRUE,
    NOW(),
    NOW()
);
