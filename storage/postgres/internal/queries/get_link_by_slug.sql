SELECT "click_count", "created_at", "expires_at", "id", "original_url", "slug"
FROM "links"
WHERE
    slug = $1
LIMIT 1;
