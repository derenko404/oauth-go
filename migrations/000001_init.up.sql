BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  -- User
  id BIGSERIAL PRIMARY KEY,
  name TEXT,
  email TEXT,
  avatar_url TEXT,
  is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,

  provider TEXT NOT NULL,
  provider_user_id TEXT NOT NULL,
  UNIQUE (provider, provider_user_id),

  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_users_username ON users (name);
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_not_deleted ON users(id) WHERE deleted_at IS NULL;

CREATE TABLE user_sessions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
  -- Session
  ip_address VARCHAR(255) NOT NULL,
  user_agent TEXT NOT NULL,
  location VARCHAR(255) NOT NULL,
  device_id VARCHAR(255) NOT NULL,
  last_active_at TIMESTAMP NOT NULL DEFAULT NOW(),

  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

COMMIT;