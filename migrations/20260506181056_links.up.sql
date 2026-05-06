CREATE TABLE IF NOT EXISTS "links" (
    "click_count"  INTEGER      NOT NULL,
    "created_at"   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    "expires_at"   TIMESTAMPTZ,
    "id"           VARCHAR(26)  PRIMARY KEY,
    "original_url" TEXT         NOT NULL,
    "slug"         VARCHAR(8)   NOT NULL UNIQUE
);
