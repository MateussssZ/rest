DROP TRIGGER IF EXISTS update_wallets_updated_at ON wallets;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS wallets CASCADE;