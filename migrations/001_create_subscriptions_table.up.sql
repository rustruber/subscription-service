-- Создание таблицы подписок
CREATE TABLE IF NOT EXISTS subscriptions (
                                             id UUID PRIMARY KEY,
                                             service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    user_id UUID NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_service_name ON subscriptions(service_name);
CREATE INDEX IF NOT EXISTS idx_subscriptions_start_date ON subscriptions(start_date);

-- Авто-обновление updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_subscriptions_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();