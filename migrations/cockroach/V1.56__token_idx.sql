CREATE INDEX IF NOT EXISTS user_user_agent_idx ON auth.tokens (user_id, user_agent_id);