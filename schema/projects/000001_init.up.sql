-- projects
CREATE TYPE projects_status AS ENUM ('Unpublished', 'Published', 'Archived');
CREATE TABLE IF NOT EXISTS projects (
                                        id UUID NOT NULL DEFAULT gen_random_uuid(),
                                        name VARCHAR(100) NOT NULL,
                                        description VARCHAR(4000),
                                        status projects_status NOT NULL DEFAULT 'Unpublished',
                                        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                        updated_at TIMESTAMP DEFAULT NOW(),
                                        deleted_at TIMESTAMP DEFAULT NULL,
                                        PRIMARY KEY (id)
);

-- projects_membership
CREATE TABLE IF NOT EXISTS projects_membership (
                                                   project_id UUID NOT NULL,
                                                   user_id UUID NOT NULL,
                                                   is_owner BOOLEAN NOT NULL DEFAULT FALSE,
                                                   added_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                                   deleted_at TIMESTAMP DEFAULT NULL,
                                                   PRIMARY KEY (project_id, user_id)
);
CREATE INDEX idx_projects_membership_user_id ON projects_membership (user_id);
CREATE INDEX idx_projects_membership_project_id ON projects_membership (project_id);