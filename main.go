package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Employee struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	BirthDay  string `json:"birthDay"`
	Gender    string `json:"gender"`
}

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}
type event struct {
	ID   int
	name string
	date string
}

var employees = []Employee{
	{ID: 1, FirstName: "Son", LastName: "Nong", BirthDay: "19.05.1999", Gender: "M"},
	{ID: 1, FirstName: "Max", LastName: "Mustermann", BirthDay: "18.04.1998", Gender: "M"},
	{ID: 1, FirstName: "Maria", LastName: "Musterfrau", BirthDay: "17.03.1997", Gender: "W"},
}
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	for _, album := range albums {
		if album.ID == id {
			c.IndentedJSON(http.StatusOK, album)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func getEmployees(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, employees)
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	err := c.BindJSON(&newAlbum)
	checkErr(err)

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func putAlbum(c *gin.Context) {
	id := c.Param("id")
	idNum, _ := strconv.Atoi(id)
	var copyAlbum album = albums[idNum]

	err := c.BindJSON(&copyAlbum)
	checkErr(err)
	albums[idNum] = copyAlbum
	c.IndentedJSON(http.StatusCreated, copyAlbum)

}

func deleteAlbum(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id > len(albums) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumById)
	router.GET("/employees", getEmployees)
	router.POST("/albums", postAlbums)
	router.PUT("albums/:id", putAlbum)

	router.Run("localhost:8080")
}
