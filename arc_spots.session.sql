-- SELECT *
-- FROM (
-- SELECT
--     name,
--     SUBSTRING(website, '^(?:https?://)?(?:www\.)?([^/#?]+)') AS domain,
--     COUNT(*) AS count,
--     STRING_AGG(name, ', ') AS spots
--   FROM "MY_TABLE"
--   GROUP BY name, domain
--   HAVING COUNT(*) > 1) AS res
-- WHERE res.domain IS NOT NULL


SELECT
  name,
  SUBSTRING(website, '^(?:https?://)?(?:www\.)?([^/#?]+)') AS domain,
  COUNT(*) AS count,
  STRING_AGG(name, ', ') AS spots
FROM "MY_TABLE"
GROUP BY name, domain
HAVING COUNT(*) > 1 AND SUBSTRING(website, '^(?:https?://)?(?:www\.)?([^/#?]+)') IS NOT NULL
ORDER BY count DESC;



