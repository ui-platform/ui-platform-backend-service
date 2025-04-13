-- screens
CREATE TYPE screens_status AS ENUM ('Unpublished', 'Published', 'Archived');
CREATE TABLE IF NOT EXISTS screens
(
    id          UUID PRIMARY KEY        DEFAULT gen_random_uuid(),
    project_id  UUID           NOT NULL,
    name        VARCHAR(100)   NOT NULL,
    description TEXT,
    status      screens_status NOT NULL DEFAULT 'Unpublished',
    widgets     JSONB          NOT NULL,
    settings    JSONB          NOT NULL,
    created_at  TIMESTAMP      NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP               DEFAULT NOW(),
    deleted_at  TIMESTAMP               DEFAULT NULL
);
CREATE INDEX IF NOT EXISTS idx_screens_project_id ON screens (project_id);

CREATE TABLE IF NOT EXISTS screens_branches
(
    id         UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    screen_id  UUID         NOT NULL,
    name       VARCHAR(100) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    UNIQUE (screen_id, name)
);
CREATE INDEX IF NOT EXISTS idx_branches_screen_id ON screens_branches (screen_id);


CREATE TABLE IF NOT EXISTS screens_widgets
(
    id         UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    branch_id  UUID      NOT NULL,
    widgets    JSONB     NOT NULL,
    settings   JSONB     NOT NULL,
    version    INT       NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE (branch_id, version)
);
CREATE INDEX IF NOT EXISTS idx_widgets_branch_id ON screens_widgets (branch_id);