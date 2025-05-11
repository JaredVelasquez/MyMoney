-- Esquema principal de la base de datos

-- Tabla de usuarios
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Índices para usuarios
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

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

-- Índices para categorías
CREATE INDEX IF NOT EXISTS idx_categories_user_id ON categories(user_id);
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);

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

-- Índices para métodos de pago
CREATE INDEX IF NOT EXISTS idx_payment_methods_user_id ON payment_methods(user_id);
CREATE INDEX IF NOT EXISTS idx_payment_methods_name ON payment_methods(name);

-- Tabla de transacciones
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
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
CREATE INDEX IF NOT EXISTS idx_transactions_currency_id ON transactions(currency_id);

-- Insertar datos de ejemplo para usuario
INSERT INTO users (id, email, name, password, created_at, updated_at)
VALUES 
    ('00000000-0000-0000-0000-000000000001', 'admin@example.com', 'Admin', '$2a$10$1qAz2wSx3eDc4rFv5tGb5edIUVkIuDGJHGNQ4qRntVkH9aHktvQaO', NOW(), NOW()),
    ('00000000-0000-0000-0000-000000000002', 'user@example.com', 'User', '$2a$10$1qAz2wSx3eDc4rFv5tGb5edIUVkIuDGJHGNQ4qRntVkH9aHktvQaO', NOW(), NOW())
ON CONFLICT (id) DO NOTHING; 