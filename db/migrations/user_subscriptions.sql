-- Tabla de suscripciones de usuarios
CREATE TABLE IF NOT EXISTS user_subscriptions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    plan_id VARCHAR(36) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('active', 'cancelled', 'expired', 'pending', 'failed')),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    renewal_date TIMESTAMP,
    cancellation_date TIMESTAMP,
    last_payment_date TIMESTAMP,
    next_payment_attempt TIMESTAMP,
    payment_method_id VARCHAR(100),
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_subscription_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_subscription_plan FOREIGN KEY (plan_id) REFERENCES plans(id)
);

-- Índices para optimizar las consultas más comunes
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_user_id ON user_subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_plan_id ON user_subscriptions(plan_id);
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_status ON user_subscriptions(status);
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_end_date ON user_subscriptions(end_date);
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_renewal_date ON user_subscriptions(renewal_date);

-- Índice compuesto para buscar suscripciones activas de un usuario (caso de uso común)
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_user_status ON user_subscriptions(user_id, status);

-- Función de actualización de timestamp
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger para actualizar automáticamente el campo updated_at
DROP TRIGGER IF EXISTS update_user_subscriptions_modtime ON user_subscriptions;
CREATE TRIGGER update_user_subscriptions_modtime
BEFORE UPDATE ON user_subscriptions
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- Inserts de ejemplo (opcional)
/*
INSERT INTO user_subscriptions (
    id, user_id, plan_id, status, start_date, end_date, 
    renewal_date, cancellation_date, last_payment_date, metadata, created_at, updated_at
) VALUES 
(
    'sub_01', 
    'user_01', -- Reemplazar con un ID de usuario real
    'plan_free', 
    'active', 
    NOW(), 
    NOW() + INTERVAL '1 year', 
    NULL, 
    NULL, 
    NULL, 
    '{"source": "registration", "notes": "Plan gratuito inicial"}',
    NOW(),
    NOW()
)
ON CONFLICT (id) DO UPDATE SET
    status = EXCLUDED.status,
    end_date = EXCLUDED.end_date,
    updated_at = NOW();
*/ 