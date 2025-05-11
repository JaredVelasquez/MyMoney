-- Tabla de categorías
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    color VARCHAR(20),
    icon VARCHAR(50),
    user_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_categories_user_id ON categories(user_id);
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);

-- Insertar algunos datos de ejemplo
INSERT INTO categories (id, name, description, color, icon, user_id, created_at, updated_at)
VALUES 
    ('11111111-1111-1111-1111-111111111101', 'Salario', 'Ingresos por trabajo', '#4CAF50', 'money', 'd15ab58a-4689-4745-bb27-46ec4757731f', NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222202', 'Inversiones', 'Ingresos por inversiones', '#2196F3', 'trending_up', 'd15ab58a-4689-4745-bb27-46ec4757731f', NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333303', 'Alimentación', 'Gastos en comida', '#F44336', 'restaurant', 'd15ab58a-4689-4745-bb27-46ec4757731f', NOW(), NOW()),
    ('44444444-4444-4444-4444-444444444404', 'Transporte', 'Gastos en transporte', '#FF9800', 'directions_car', 'd15ab58a-4689-4745-bb27-46ec4757731f', NOW(), NOW()),
    ('55555555-5555-5555-5555-555555555505', 'Entretenimiento', 'Gastos en ocio', '#9C27B0', 'movie', 'd15ab58a-4689-4745-bb27-46ec4757731f', NOW(), NOW()),
    ('66666666-6666-6666-6666-666666666606', 'Freelance', 'Ingresos por trabajos freelance', '#4CAF50', 'work', 'd15ab58a-4689-4745-bb27-46ec4757731f', NOW(), NOW())
ON CONFLICT (id) DO NOTHING; 