CREATE TABLE IF NOT EXISTS movies (
  id bigserial PRIMARY KEY,
  created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),
  title text NOT NULL,
  year integer NOT NULL,
  duration integer NOT NULL,
  genres text [] NOT NULL,
  version integer NOT NULL DEFAULT 1
)