package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/mtp721/micobo-assignment/pkg/db"
	"github.com/mtp721/micobo-assignment/pkg/handlers"
)

func main() {
	//Init db
	db := db.Init()
	defer db.Close()

	//handler object with handler methods
	h := handlers.New(db)

	// API Endpoints
	router := gin.Default()
	router.GET("/employees", h.GetEmployees)          //get all employees
	router.POST("/employees", h.PostEmployee)         //registers new employee
	router.PUT("/employees/:id", h.PutEmployee)       //update employees info
	router.DELETE("/employees/:id", h.DeleteEmployee) //delete specified employee

	router.GET("/events", h.GetEvents)    //get all upcoming events
	router.GET("/events/:id", h.GetEvent) //get specific event

	/*returns the list of the employees that are assisting to the event,
	should accept query parameters for filtering if they need or don't need accommodation*/
	router.GET("/events/:id/employees", h.GetEmployeesForEvent)

	router.Run("localhost:8080")
}
