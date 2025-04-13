DROP TABLE IF EXISTS projects;
DROP TYPE IF EXISTS projects_status;

DROP TABLE IF EXISTS projects_membership;
DROP INDEX IF EXISTS idx_projects_membership_user_id;
DROP INDEX IF EXISTS idx_projects_membership_project_id;
