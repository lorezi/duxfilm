CREATE TABLE IF NOT EXISTS permissions (
  id bigserial PRIMARY KEY,
  code text NOT NULL
);
CREATE TABLE IF NOT EXISTS users_permissions (
  user_id BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
  permission_id BIGINT NOT NULL REFERENCES permissions ON DELETE CASCADE,
  PRIMARY KEY (user_id, permission_id)
) -- Add the two permissions to the table
INSERT INTO
  permissions (code)
VALUES
  ('movies:read'),
  ('movies:write');