-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS users (
  "id" uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  "name" VARCHAR(50) NOT NULL,
  "email" VARCHAR(100) UNIQUE NOT NULL,
  "created_at" TIMESTAMP DEFAULT NOW()
);
---- create above / drop below ----
DROP TABLE IF EXISTS users;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.