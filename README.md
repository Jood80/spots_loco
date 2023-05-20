# spots_loco


### How to run;
1. Clone the repo and `cd spots_loco`
2. Create your database, and make sure to update `database.config.go` file with your valid credentials and your db name.
3. Run `go run main.go` || `nodemon --exec go run main.go --signal SIGTERM` 
4. Hit `http://localhost:8080/spots?lat=1&long=2&radius=10&type=circle` 
5. Your response may look like this:
```ts
[
    {
        "id": "0e53e421-e506-4c46-a7a1-9751d9581651",
        "name": "Vegan Yes",
        "website": {
            "String": "http://www.veganyes.co.uk/",
            "Valid": true
        },
        "coordinates": "0101000020E610000040D42247DF49B2BF1DFC1FBB66C24940",
        "rating": 9.997621583951322,
        "lat": -0.0714397,
        "long": 51.51876009999999,
        "distance": 49.53035013879962
    },
...
]
```
