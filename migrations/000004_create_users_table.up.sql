CREATE Table If NOT EXISTS users (
  id bigserial PRIMARY KEY,
  created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
  username text NOT NULL,
  email citext UNIQUE NOT NULL,
  password_hash bytea NOT NULL,
  activated bool NOT NULL,
  version integer NOT NULL DEFAULT 1
);
CREATE EXTENSION citext