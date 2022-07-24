package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/mtp721/micobo-assignment/pkg/models"
)

type handler struct {
	DB *sql.DB
}

func New(db *sql.DB) handler {
	return handler{db}
}

// Returns a list of all employees
func (h handler) GetEmployees(c *gin.Context) {
	var employees []models.Employee
	rows, err := h.DB.Query(`SELECT * FROM employees ORDER BY id`)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var employee models.Employee
		if err := rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.BirthDay, &employee.Gender); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		employees = append(employees, employee)
	}
	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, employees)
}

// register a new employee
func (h handler) PostEmployee(c *gin.Context) {
	var employee models.Employee
	if err := c.BindJSON(&employee); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "binding error: " + err.Error()})
		return
	}
	fmt.Println(employee)
	row := h.DB.QueryRow(`INSERT INTO employees (first_name, last_name, birthday, gender) VALUES ($1, $2, $3, $4) RETURNING id`,
		employee.FirstName, employee.LastName, employee.BirthDay, employee.Gender)

	if err := row.Scan(&employee.ID); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, employee)

}

// modfiy employee
func (h handler) PutEmployee(c *gin.Context) {
	var employee models.Employee
	id := c.Param("id")

	//Query employee with the specified id and store old values
	row := h.DB.QueryRow("SELECT * FROM employees WHERE id = $1", id)
	if err := row.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.BirthDay, &employee.Gender); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Error querying db: " + err.Error()})
		return
	}
	//Override values of employee with updated values from request
	if err := c.BindJSON(&employee); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "binding error: " + err.Error()})
		return
	}
	//Update employee in db
	updRow := h.DB.QueryRow("UPDATE employees SET first_name = $1, last_name = $2, birthday = $3, gender = $4 WHERE id = $5 RETURNING *",
		employee.FirstName, employee.LastName, employee.BirthDay, employee.Gender, id)

	//Write returned values from db to employee to make sure values were updated correctly
	if err := updRow.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.BirthDay, &employee.Gender); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error updating employee: " + err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, employee)
}

// delete employee from db
func (h handler) DeleteEmployee(c *gin.Context) {
	var employee models.Employee
	id := c.Param("id")
	row := h.DB.QueryRow("DELETE FROM employees WHERE id = $1 RETURNING *", id)
	if err := row.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.BirthDay, &employee.Gender); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Error deleting employee: " + err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, employee)
}

// Returns a list of all events
func (h handler) GetEvents(c *gin.Context) {
	var events []models.Event
	rows, err := h.DB.Query(`SELECT * FROM events ORDER BY id`)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Date); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, events)
}

// get event specified by id
func (h handler) GetEvent(c *gin.Context) {
	var event models.Event
	id := c.Param("id")

	//Query event with the specified id
	row := h.DB.QueryRow("SELECT * FROM events WHERE id = $1", id)
	if err := row.Scan(&event.ID, &event.Name, &event.Date); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Error querying db: " + err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, event)
}

/*returns the list of the employees that are attending the event specified by event_id,
accepts query parameter for filtering if an employee need accommodation or not*/
func (h handler) GetEmployeesForEvent(c *gin.Context) {
	var employees []models.Employee
	eventId := c.Param("id")
	accommodation := c.Query("accommodation")

	var accommodationQuery string
	switch accommodation {
	case "true":
		accommodationQuery = "AND accommodation = true"
	case "false":
		accommodationQuery = "AND accommodation = false"
	default:
		accommodationQuery = ""
	}

	query := fmt.Sprintf(`SELECT id, first_name, last_name, birthday, gender FROM employees JOIN attendances 
		ON attendances.employee_id = employees.id WHERE attendances.event_id = $1 %s ORDER BY id`, accommodationQuery)

	rows, err := h.DB.Query(query, eventId)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var employee models.Employee
		if err := rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.BirthDay, &employee.Gender); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		employees = append(employees, employee)
	}
	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, employees)

}
