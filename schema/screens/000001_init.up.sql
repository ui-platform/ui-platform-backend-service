-- projects
CREATE TYPE screens_status AS ENUM ('Unpublished', 'Published', 'Archived');
CREATE TABLE IF NOT EXISTS screens (
                                        id UUID NOT NULL DEFAULT gen_random_uuid(),
                                        project_id UUID NOT NULL,
                                        name VARCHAR(100) NOT NULL,
                                        description VARCHAR(4000),
                                        status projects_status NOT NULL DEFAULT 'Unpublished',
                                        content JSONB NOT NULL,
                                        settings JSONB NOT NULL,
                                        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                        updated_at TIMESTAMP DEFAULT NOW(),
                                        deleted_at TIMESTAMP DEFAULT NULL,
                                        PRIMARY KEY (id)
);