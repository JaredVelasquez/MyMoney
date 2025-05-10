-- Tabla para almacenar las monedas
CREATE TABLE IF NOT EXISTS currencies (
    id UUID PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Índices para optimizar búsquedas
CREATE INDEX IF NOT EXISTS idx_currencies_code ON currencies(code);
CREATE INDEX IF NOT EXISTS idx_currencies_is_active ON currencies(is_active);

-- Insertar monedas iniciales
INSERT INTO currencies (id, code, name, symbol, is_active, created_at, updated_at)
VALUES 
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'USD', 'Dólar estadounidense', '$', TRUE, NOW(), NOW()),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'EUR', 'Euro', '€', TRUE, NOW(), NOW()),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'GBP', 'Libra esterlina', '£', TRUE, NOW(), NOW()),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'JPY', 'Yen japonés', '¥', TRUE, NOW(), NOW()),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'MXN', 'Peso mexicano', '$', TRUE, NOW(), NOW()),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', 'BTC', 'Bitcoin', '₿', TRUE, NOW(), NOW())
ON CONFLICT (code) DO UPDATE 
SET name = EXCLUDED.name,
    symbol = EXCLUDED.symbol,
    is_active = EXCLUDED.is_active,
    updated_at = NOW(); 