-- Tabla de transacciones actualizada con currency_id como UUID
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    amount DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    description TEXT,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    category_id UUID NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('INCOME', 'EXPENSE')),
    payment_method_id UUID,
    currency_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT,
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (currency_id) REFERENCES currencies(id)
);

-- Índices para transacciones
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions(category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_payment_method_id ON transactions(payment_method_id);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
CREATE INDEX IF NOT EXISTS idx_transactions_currency_id ON transactions(currency_id);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);




-- Insertar algunos datos de ejemplo
INSERT INTO transactions (id, amount, description, date, category_id, type, payment_method_id, user_id, currency_id, created_at, updated_at)
VALUES 
    ('11111111-1111-1111-1111-111111111201', 1500.00, 'Salario mensual', '2023-05-01 12:00:00+00', '11111111-1111-1111-1111-111111111101', 'INCOME', '11111111-1111-1111-1111-111111111111', 'd15ab58a-4689-4745-bb27-46ec4757731f', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222202', 50.00, 'Dividendos', '2023-05-05 14:30:00+00', '22222222-2222-2222-2222-222222222202', 'INCOME', '11111111-1111-1111-1111-111111111111', 'd15ab58a-4689-4745-bb27-46ec4757731f', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333203', 120.50, 'Compra supermercado', '2023-05-10 18:45:00+00', '33333333-3333-3333-3333-333333333303', 'EXPENSE', '22222222-2222-2222-2222-222222222222', 'd15ab58a-4689-4745-bb27-46ec4757731f', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NOW(), NOW()),
    ('44444444-4444-4444-4444-444444444204', 35.00, 'Gasolina', '2023-05-12 10:15:00+00', '44444444-4444-4444-4444-444444444404', 'EXPENSE', '11111111-1111-1111-1111-111111111111', 'd15ab58a-4689-4745-bb27-46ec4757731f', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NOW(), NOW()),
    ('55555555-5555-5555-5555-555555555205', 80.00, 'Cine y cena', '2023-05-15 20:30:00+00', '55555555-5555-5555-5555-555555555505', 'EXPENSE', '11111111-1111-1111-1111-111111111111', 'd15ab58a-4689-4745-bb27-46ec4757731f', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NOW(), NOW()),
    ('66666666-6666-6666-6666-666666666206', 500.00, 'Proyecto freelance', '2023-05-20 09:00:00+00', '66666666-6666-6666-6666-666666666606', 'INCOME', '33333333-3333-3333-3333-333333333333', 'd15ab58a-4689-4745-bb27-46ec4757731f', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NOW(), NOW())
ON CONFLICT (id) DO NOTHING; 