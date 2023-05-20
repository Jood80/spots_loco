/**
  * Returns spots which have a domain with count greater than 1 
  * Changes the website field, so it contains the domain only 
  * Counts how many spots have the same domain
*/

SELECT
  name,
  SUBSTRING(website FROM '^(?:https?://)?(?:www\.)?([^/#?]+)') AS domain,
  COUNT(*) AS count,
  STRING_AGG(name, ', ') AS spots
FROM "MY_TABLE"
WHERE website IS NOT NULL
GROUP BY name, domain
HAVING COUNT(*) > 1
ORDER BY count DESC;


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
