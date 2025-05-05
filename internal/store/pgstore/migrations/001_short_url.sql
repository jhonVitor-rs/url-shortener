-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS short_urls (
  "id" uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  "slug" TEXT UNIQUE NOT NULL,
  "original_url" TEXT NOT NULL,
  "created_at" TIMESTAMP DEFAULT NOW(),
  "expires_at" TIMESTAMP,
  "access_count" INTEGER DEFAULT 0
);
---- create above / drop below ----
DROP TABLE IF EXISTS short_urls;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.