\c "TravelerTrack"

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);

CREATE TABLE IF NOT EXISTS places (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    location JSONB,
    address VARCHAR(255),
    rating INTEGER,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_places_deleted_at ON places (deleted_at);

CREATE TABLE IF NOT EXISTS user_places (
    user_id UUID NOT NULL,
    place_id UUID NOT NULL,

    PRIMARY KEY (user_id, place_id),

    CONSTRAINT fk_user_places_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE,

    CONSTRAINT fk_user_places_place
        FOREIGN KEY (place_id)
        REFERENCES places (id)
        ON DELETE CASCADE
);

SELECT 'Database schema created successfully (or already existed).' AS status;