-- Tabla de categorías
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL CHECK (type IN ('INCOME', 'EXPENSE')),
    color VARCHAR(20),
    icon VARCHAR(50),
    user_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_categories_user_id ON categories(user_id);
CREATE INDEX IF NOT EXISTS idx_categories_type ON categories(type);
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);

-- Insertar algunos datos de ejemplo
INSERT INTO categories (id, name, description, type, color, icon, user_id, created_at, updated_at)
VALUES 
    ('11111111-1111-1111-1111-111111111101', 'Salario', 'Ingresos por trabajo', 'INCOME', '#4CAF50', 'money', '00000000-0000-0000-0000-000000000001', NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222202', 'Inversiones', 'Ingresos por inversiones', 'INCOME', '#2196F3', 'trending_up', '00000000-0000-0000-0000-000000000001', NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333303', 'Alimentación', 'Gastos en comida', 'EXPENSE', '#F44336', 'restaurant', '00000000-0000-0000-0000-000000000001', NOW(), NOW()),
    ('44444444-4444-4444-4444-444444444404', 'Transporte', 'Gastos en transporte', 'EXPENSE', '#FF9800', 'directions_car', '00000000-0000-0000-0000-000000000001', NOW(), NOW()),
    ('55555555-5555-5555-5555-555555555505', 'Entretenimiento', 'Gastos en ocio', 'EXPENSE', '#9C27B0', 'movie', '00000000-0000-0000-0000-000000000002', NOW(), NOW()),
    ('66666666-6666-6666-6666-666666666606', 'Freelance', 'Ingresos por trabajos freelance', 'INCOME', '#4CAF50', 'work', '00000000-0000-0000-0000-000000000002', NOW(), NOW())
ON CONFLICT (id) DO NOTHING; 