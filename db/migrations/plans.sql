-- Tabla para almacenar los planes de suscripción
CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    currency_id UUID NOT NULL REFERENCES currencies(id),
    interval VARCHAR(50) NOT NULL CHECK (interval IN ('monthly', 'yearly')),
    features JSONB NOT NULL DEFAULT '[]'::JSONB,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Índices para optimizar búsquedas
CREATE INDEX IF NOT EXISTS idx_plans_is_active ON plans(is_active);
CREATE INDEX IF NOT EXISTS idx_plans_is_public ON plans(is_public);
CREATE INDEX IF NOT EXISTS idx_plans_sort_order ON plans(sort_order);

-- Insertar planes iniciales
INSERT INTO plans (id, name, description, price, currency_id, interval, features, is_active, is_public, sort_order, created_at, updated_at)
VALUES 
    (
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380b11', 
        'Gratis', 
        'Plan básico con funcionalidades limitadas', 
        0, 
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', -- USD
        'monthly',
        '[
            {"name": "Procesamiento de texto", "description": "Generación de texto mediante IA", "value": "Ilimitado", "included": true},
            {"name": "Procesamiento de imágenes", "description": "Generación y edición de imágenes mediante IA", "value": "No disponible", "included": false},
            {"name": "Procesamiento de audio", "description": "Transcripción y generación de audio mediante IA", "value": "No disponible", "included": false}
        ]'::JSONB,
        TRUE,
        TRUE,
        1,
        NOW(),
        NOW()
    ),
    (
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380b12', 
        'Pro', 
        'Plan profesional con todas las funcionalidades', 
        19.99, 
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', -- USD
        'monthly',
        '[
            {"name": "Procesamiento de texto", "description": "Generación de texto mediante IA", "value": "Ilimitado", "included": true},
            {"name": "Procesamiento de imágenes", "description": "Generación y edición de imágenes mediante IA", "value": "Ilimitado", "included": true},
            {"name": "Procesamiento de audio", "description": "Transcripción y generación de audio mediante IA", "value": "Ilimitado", "included": true}
        ]'::JSONB,
        TRUE,
        TRUE,
        2,
        NOW(),
        NOW()
    ),
    (
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380b13', 
        'Pro Anual', 
        'Plan profesional con todas las funcionalidades - Facturación anual', 
        199.90, 
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', -- USD
        'yearly',
        '[
            {"name": "Procesamiento de texto", "description": "Generación de texto mediante IA", "value": "Ilimitado", "included": true},
            {"name": "Procesamiento de imágenes", "description": "Generación y edición de imágenes mediante IA", "value": "Ilimitado", "included": true},
            {"name": "Procesamiento de audio", "description": "Transcripción y generación de audio mediante IA", "value": "Ilimitado", "included": true}
        ]'::JSONB,
        TRUE,
        TRUE,
        3,
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO UPDATE 
SET name = EXCLUDED.name,
    description = EXCLUDED.description,
    price = EXCLUDED.price,
    currency_id = EXCLUDED.currency_id,
    interval = EXCLUDED.interval,
    features = EXCLUDED.features,
    is_active = EXCLUDED.is_active,
    is_public = EXCLUDED.is_public,
    sort_order = EXCLUDED.sort_order,
    updated_at = NOW(); 