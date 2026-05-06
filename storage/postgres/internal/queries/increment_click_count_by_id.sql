UPDATE "links"
SET
    click_count = "click_count" + 1
WHERE
    id = $1;
