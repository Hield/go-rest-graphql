package main

import (
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/gin-gonic/gin"
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
	controller := &Controller{database: database}

	router.GET("/stations", controller.GetStations)
	router.POST("/stations", controller.CreateStation)

	router.Run(":8080")
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
}

func (controller *Controller) GetStations(c *gin.Context) {
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