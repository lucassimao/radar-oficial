-- Add Diários dos Municípios do Piauí institution
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
    2,
    'Diário Oficial dos Municípios do Piauí',
    'municipios-pi',
    'municipal',
    'PI',
    NULL,
    'https://www.diarioficialdosmunicipios.org/',
    TRUE,
    NOW(),
    NOW()
)
ON CONFLICT (id) DO NOTHING;