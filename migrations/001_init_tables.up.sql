

-- Создание таблицы wallets
CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_wallets_wallet_id ON wallets(id);
CREATE INDEX IF NOT EXISTS idx_wallets_status ON wallets(status);

ALTER TABLE wallets 
ADD CONSTRAINT wallets_balance_check 
CHECK (balance >= 0);

ALTER TABLE wallets 
ADD CONSTRAINT wallets_status_check 
CHECK (status IN ('ACTIVE', 'FROZEN', 'CLOSED'));

-- Создание функции для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';


CREATE TRIGGER update_wallets_updated_at 
    BEFORE UPDATE ON wallets 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
