-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS short_urls (
  "id" uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "slug" TEXT UNIQUE NOT NULL,
  "original_url" TEXT NOT NULL,
  "created_at" TIMESTAMP DEFAULT NOW(),
  "expires_at" TIMESTAMP,
  "access_count" INTEGER DEFAULT 0,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
---- create above / drop below ----
DROP TABLE IF EXISTS short_urls;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.