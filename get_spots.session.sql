-- radius = 10 | lat = 1 | long= 2 | shape = Square

SELECT id, ST_X(coordinates::geometry) AS lat, ST_Y(coordinates::geometry)AS long
FROM "MY_TABLE"
WHERE ST_X(coordinates::geometry) >= -9 AND ST_X(coordinates::geometry) <= 11 AND ST_Y(coordinates::geometry)>= -8 and ST_Y(coordinates::geometry)<=12;

-- radius = 10 | lat = 1 | long= 2 | shape = CIRCLE
SELECT id, ST_X(coordinates::geometry) AS lat, ST_Y(coordinates::geometry)AS long
FROM "MY_TABLE"
WHERE SQRT(ABS((1- ST_X(coordinates::geometry))* (1-ST_X(coordinates::geometry))) + ABS((2-ST_Y(coordinates::geometry))*(2-ST_Y(coordinates::geometry)))) <= 10;


SELECT
	id,
	rating,
	ST_X(coordinates::geometry) AS lat,
	ST_Y(coordinates::geometry) AS long,
	SQRT(ABS((1- ST_X(coordinates::geometry))* (1-ST_X(coordinates::geometry))) + ABS((2-ST_Y(coordinates::geometry))*(2-ST_Y(coordinates::geometry)))) AS distance
FROM "MY_TABLE"
WHERE SQRT(ABS((1- ST_X(coordinates::geometry))* (1-ST_X(coordinates::geometry))) + ABS((2-ST_Y(coordinates::geometry))*(2-ST_Y(coordinates::geometry)))) <= 100
GROUP BY id, rating, lat, long
ORDER BY  
Case 
WHEN SQRT(ABS((1- ST_X(coordinates::geometry))* (1-ST_X(coordinates::geometry))) + ABS((2-ST_Y(coordinates::geometry))*(2-ST_Y(coordinates::geometry))))  <=50
Then 0
ELSE SQRT(ABS((1- ST_X(coordinates::geometry))* (1-ST_X(coordinates::geometry))) + ABS((2-ST_Y(coordinates::geometry))*(2-ST_Y(coordinates::geometry)))) 
END ASC, rating DESC;

