CREATE TABLE IF NOT EXISTS sessions (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expiration TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    PRIMARY KEY (id)
);