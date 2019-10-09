package main

import (
	"net/http"
	"database/sql"
	"context"
	_ "github.com/mattn/go-sqlite3"
	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
	"os"
)

type Station struct {
	ID              string `json:"stationId"`
	Name            string `json:"name"`
	BikesAvailable  int    `json:"bikesAvailable"`
	SpacesAvailable int    `json:"spacesAvailable"`
}

func main() {
	database := getDB()
	router := gin.Default()
	client := graphql.NewClient("https://api.digitransit.fi/routing/v1/routers/finland/index/graphql", nil)
	controller := &Controller{database: database, client: client}

	router.GET("/stations", controller.GetStations)
	router.POST("/stations", controller.CreateStation)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

func getDB() *sql.DB {
	database, err := sql.Open("sqlite3", "./stations.db")
	if err != nil {
		panic(err)
	}
	_, err = database.Exec(`CREATE TABLE IF NOT EXISTS stations (
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT,
		bikes_available INTEGER,
		spaces_available INTEGER
	)`)
	if err != nil {
		panic(err)
	}
	return database
}

type Controller struct {
	database *sql.DB
	client *graphql.Client
}

func (controller *Controller) GetStations(c *gin.Context) {
	// Get stations from local database
	rows, err := controller.database.Query(`SELECT * FROM stations`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var (
		id string
		name string
		bikesAvailable int
		spacesAvailable int
	)
	var stations = make([]Station, 0)
	for rows.Next() {
		err = rows.Scan(&id, &name, &bikesAvailable, &spacesAvailable)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		stations = append(stations, Station{id, name, bikesAvailable, spacesAvailable})
	}

	// Get stations from GraphQL
	var query struct {
		BikeRentalStations []struct {
			StationId graphql.String
			Name graphql.String
			BikesAvailable graphql.Int
			SpacesAvailable graphql.Int
		}
	}
	err = controller.client.Query(context.Background(), &query, nil)
	if err != nil {
		// Return only the ones from local database if can't query
		c.JSON(http.StatusOK, stations)
		return
	}
	for _, station := range query.BikeRentalStations {
		stations = append(stations, Station{
			string(station.StationId),
			string(station.Name),
			int(station.BikesAvailable),
			int(station.SpacesAvailable),
		})
	}
	c.JSON(http.StatusOK, stations)
}

func (controller *Controller) CreateStation(c *gin.Context) {
	var station Station
	err := c.ShouldBindJSON(&station)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	statement, err := controller.database.Prepare(`INSERT INTO stations(id, name, bikes_available, spaces_available) values(?, ?, ?, ?)`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = statement.Exec(station.ID, station.Name, station.BikesAvailable, station.SpacesAvailable)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusCreated, "success")
}