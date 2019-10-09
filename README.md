# Go REST + GraphQL client test

REST APIs + GraphQL client for retrieving and creating new bike stations

## Requirements

- Go 1.12+ ([Download](https://golang.org/dl/))

## Run the server

- Windows: `go build -o bin/go-rest-graphql.exe && bin\go-rest-graphql.exe` at root
- Others: `go build -o bin/go-rest-graphql && bin/go-rest-graphql` at root

## Endpoints

- GET `/stations`: List stations
- POST `/stations`: Create a new station with signature
```
{
  "stationId": string,
  "name": string,
  "bikesAvailable": int,
  "spacesAvailable": int
}
```

## Notes

- A deployed version is available [here](https://go-rest-graphql.herokuapp.com/stations)

## Troubleshooting

- `gcc` related error when building on Windows: Download and install MinGW-w64 [here](https://sourceforge.net/projects/mingw-w64/files/latest/download)
- Error getting packages when building: Try running `go get {package name here}` for failing packages
