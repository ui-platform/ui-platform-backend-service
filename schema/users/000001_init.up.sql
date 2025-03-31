CREATE TABLE IF NOT EXISTS users (
                                     id UUID NOT NULL DEFAULT gen_random_uuid(),
                                     email VARCHAR(100) NOT NULL UNIQUE,
                                     password VARCHAR(255) NOT NULL,
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                     updated_at TIMESTAMP DEFAULT NOW(),
                                     deleted_at TIMESTAMP DEFAULT NULL,
                                     PRIMARY KEY (id)
);