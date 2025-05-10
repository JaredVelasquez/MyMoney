-- Tabla de métodos de pago
CREATE TABLE IF NOT EXISTS payment_methods (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    user_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_payment_methods_user_id ON payment_methods(user_id);
CREATE INDEX IF NOT EXISTS idx_payment_methods_name ON payment_methods(name);

-- Insertar algunos datos de ejemplo
INSERT INTO payment_methods (id, name, description, is_active, user_id, created_at, updated_at)
VALUES 
    ('11111111-1111-1111-1111-111111111111', 'Tarjeta de Crédito', 'Visa terminada en 4242', TRUE, '00000000-0000-0000-0000-000000000001', NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222222', 'Efectivo', 'Pagos en efectivo', TRUE, '00000000-0000-0000-0000-000000000001', NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333333', 'Transferencia Bancaria', 'Cuenta corriente', TRUE, '00000000-0000-0000-0000-000000000002', NOW(), NOW())
ON CONFLICT (id) DO NOTHING; 