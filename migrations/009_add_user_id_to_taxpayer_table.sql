-- Add user_id to taxpayers for one-to-one link
ALTER TABLE taxpayers ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;

-- Indexes for the new relationship
CREATE INDEX IF NOT EXISTS idx_taxpayers_user_id ON taxpayers(user_id);
CREATE INDEX IF NOT EXISTS idx_taxpayers_national_id ON taxpayers(national_id);

-- Sync updated_at trigger (as before)
CREATE OR REPLACE FUNCTION sync_updated_at() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_users_updated_at ON users;
CREATE TRIGGER trigger_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION sync_updated_at();

DROP TRIGGER IF EXISTS trigger_taxpayers_updated_at ON taxpayers;
CREATE TRIGGER trigger_taxpayers_updated_at BEFORE UPDATE ON taxpayers FOR EACH ROW EXECUTE FUNCTION sync_updated_at();