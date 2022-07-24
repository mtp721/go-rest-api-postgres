package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/mtp721/micobo-assignment/pkg/handlers"
	"github.com/mtp721/micobo-assignment/pkg/models"
	"github.com/stretchr/testify/assert"
)

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}
func TestGetEmployees(t *testing.T) {
	//Init mock db
	db, mock := newMock()

	h := handlers.New(db)
	//Init router
	router := gin.Default()
	router.GET("/employees", h.GetEmployees) //get all employees

	//http get request
	req, _ := http.NewRequest("GET", "/employees", nil)

	//mock db should return this on specified query
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "birthday", "gender"}).AddRow(
		1, "Son", "Nong", "1999-05-19", "m").AddRow(2, "Max", "Mustermann", "1998-04-18", "m")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM employees ORDER BY id")).WillReturnRows(rows)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expectedResp := `[
		{
			"id": 1,
			"firstName": "Son",
			"lastName": "Nong",
			"birthDay": "1999-05-19",
			"gender": "m"
		},
		{
			"id": 2,
			"firstName": "Max",
			"lastName": "Mustermann",
			"birthDay": "1998-04-18",
			"gender": "m"
		}
	]`
	assert := assert.New(t)
	assert.Equal(http.StatusOK, w.Code, "expected http Code 200")
	assert.JSONEq(expectedResp, w.Body.String(), "Response body doesn't match")
	/*
		t.Logf("status: %d", w.Code)
		t.Logf("response: %s", w.Body.String())
	*/
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestPostEmployee(t *testing.T) {
	//Init mock db
	db, mock := newMock()

	h := handlers.New(db)
	//Init router
	router := gin.Default()
	router.POST("/employees", h.PostEmployee) //registers new employee

	//http request
	req, _ := http.NewRequest("POST", "/employees", strings.NewReader(
		`{"firstName": "Joe", "lastName": "Jones", "birthday": "1997-09-12", "gender": "m"}`))

	//mock db should return this on specified query
	rows := sqlmock.NewRows([]string{"id"}).AddRow(3)

	emp := models.Employee{ID: 3, FirstName: "Joe", LastName: "Jones", BirthDay: "1997-09-12", Gender: "m"}

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO employees (first_name, last_name, birthday, gender) VALUES ($1, $2, $3, $4) RETURNING id`)).WithArgs(
		emp.FirstName, emp.LastName, emp.BirthDay, emp.Gender).WillReturnRows(rows)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expectedResp := `{
			"id": 3,
			"firstName": "Joe",
			"lastName": "Jones",
			"birthDay": "1997-09-12",
			"gender": "m"
		}`
	assert := assert.New(t)
	assert.Equal(http.StatusCreated, w.Code, "expected http Code 201")
	assert.JSONEq(expectedResp, w.Body.String(), "Response body doesn't match")

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
func TestPutEmployee(t *testing.T) {
	//Init mock db
	db, mock := newMock()

	h := handlers.New(db)
	//Init router
	router := gin.Default()
	router.PUT("/employees/:id", h.PutEmployee) //update employees info

	//http request
	req, _ := http.NewRequest("PUT", "/employees/3", strings.NewReader(
		`{"firstName": "Geo", "lastName": "Dude"}`))

	//employee to test
	emp := models.Employee{ID: 3, FirstName: "Joe", LastName: "Jones", BirthDay: "1997-09-12", Gender: "m"}
	//employee after update
	updEmp := models.Employee{ID: 3, FirstName: "Geo", LastName: "Dude", BirthDay: "1997-09-12", Gender: "m"}

	//row for select query
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "birthday", "gender"}).AddRow(
		emp.ID, emp.FirstName, emp.LastName, emp.BirthDay, emp.Gender)
	//row after specified employee is updated in db
	updRows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "birthday", "gender"}).AddRow(
		updEmp.ID, updEmp.FirstName, updEmp.LastName, updEmp.BirthDay, updEmp.Gender)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM employees WHERE id = $1")).WithArgs(strconv.Itoa(emp.ID)).WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(
		"UPDATE employees SET first_name = $1, last_name = $2, birthday = $3, gender = $4 WHERE id = $5 RETURNING *")).WithArgs(
		updEmp.FirstName, updEmp.LastName, updEmp.BirthDay, updEmp.Gender, strconv.Itoa(updEmp.ID)).WillReturnRows(updRows)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expectedResp := `{
			"id": 3,
			"firstName": "Geo",
			"lastName": "Dude",
			"birthDay": "1997-09-12",
			"gender": "m"
		}`
	assert := assert.New(t)
	assert.Equal(http.StatusOK, w.Code, "http Code doesn't match")
	assert.JSONEq(expectedResp, w.Body.String(), "Response body doesn't match")

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestDeleteEmployee(t *testing.T) {
	//Init mock db
	db, mock := newMock()

	h := handlers.New(db)
	//Init router
	router := gin.Default()
	router.DELETE("/employees/:id", h.DeleteEmployee) //delete specified employee

	//http request
	req, _ := http.NewRequest("DELETE", "/employees/3", nil)

	//employee to test
	emp := models.Employee{ID: 3, FirstName: "Joe", LastName: "Jones", BirthDay: "1997-09-12", Gender: "m"}

	//row for select query
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "birthday", "gender"}).AddRow(
		emp.ID, emp.FirstName, emp.LastName, emp.BirthDay, emp.Gender)

	mock.ExpectQuery(regexp.QuoteMeta(
		"DELETE FROM employees WHERE id = $1 RETURNING *")).WithArgs(strconv.Itoa(emp.ID)).WillReturnRows(rows)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expectedResp := `{
			"id": 3,
			"firstName": "Joe",
			"lastName": "Jones",
			"birthDay": "1997-09-12",
			"gender": "m"
		}`
	assert := assert.New(t)
	assert.Equal(http.StatusOK, w.Code, "http Code doesn't match")
	assert.JSONEq(expectedResp, w.Body.String(), "Response body doesn't match")

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestGetEvents(t *testing.T) {
	//Init mock db
	db, mock := newMock()

	h := handlers.New(db)
	//Init router
	router := gin.Default()
	router.GET("/events", h.GetEvents) //get all upcoming events

	//http request
	req, _ := http.NewRequest("GET", "/events", nil)

	//row for select query
	rows := sqlmock.NewRows([]string{"id", "name", "date"}).AddRow(
		1, "Costume Party", "2022-08-01").AddRow(2, "Escape Room", "2022-08-02")

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM events ORDER BY id")).WillReturnRows(rows)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expectedResp := `[
		{
			"id": 1,
			"name": "Costume Party",
			"date": "2022-08-01"
		},
		{
			"id": 2,
			"name": "Escape Room",
			"date": "2022-08-02"
		}
	]`
	assert := assert.New(t)
	assert.Equal(http.StatusOK, w.Code, "http Code doesn't match")
	assert.JSONEq(expectedResp, w.Body.String(), "Response body doesn't match")

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetEvent(t *testing.T) {
	//Init mock db
	db, mock := newMock()

	h := handlers.New(db)
	//Init router
	router := gin.Default()
	router.GET("/events/:id", h.GetEvent) //get specific event

	//http request
	req, _ := http.NewRequest("GET", "/events/1", nil)

	// Event to test
	event := models.Event{ID: 1, Name: "Costume Party", Date: "2022-08-01"}

	//row for select query
	rows := sqlmock.NewRows([]string{"id", "name", "date"}).AddRow(
		event.ID, event.Name, event.Date)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM events WHERE id = $1")).WithArgs(strconv.Itoa(event.ID)).WillReturnRows(rows)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expectedResp := `{
			"id": 1,
			"name": "Costume Party",
			"date": "2022-08-01"
		}`
	assert := assert.New(t)
	assert.Equal(http.StatusOK, w.Code, "http Code doesn't match")
	assert.JSONEq(expectedResp, w.Body.String(), "Response body doesn't match")

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

//Return a list of all employees that attend a specified event
func TestGetEmployeesForEvent(t *testing.T) {
	//Init mock db
	db, mock := newMock()

	h := handlers.New(db)
	//Init router
	router := gin.Default()
	router.GET("/events/:id/employees", h.GetEmployeesForEvent)

	//http request
	req, _ := http.NewRequest("GET", "/events/1/employees", nil)

	// Event to test
	event := models.Event{ID: 1, Name: "Costume Party", Date: "2022-08-01"}

	query := `SELECT id, first_name, last_name, birthday, gender FROM employees JOIN attendances 
		ON attendances.employee_id = employees.id WHERE attendances.event_id = $1 ORDER BY id`

	//row for select query
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "birthday", "gender"}).AddRow(
		1, "Son", "Nong", "1999-05-19", "m").AddRow(2, "Max", "Mustermann", "1998-04-18", "m").AddRow(3, "Joe", "Jones", "1997-09-12", "m")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(strconv.Itoa(event.ID)).WillReturnRows(rows)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expectedResp := `[
		{
			"id": 1,
			"firstName": "Son",
			"lastName": "Nong",
			"birthDay": "1999-05-19",
			"gender": "m"
		},
		{
			"id": 2,
			"firstName": "Max",
			"lastName": "Mustermann",
			"birthDay": "1998-04-18",
			"gender": "m"
		},
		{
			"id": 3,
			"firstName": "Joe",
			"lastName": "Jones",
			"birthDay": "1997-09-12",
			"gender": "m"
		}
	]`

	assert := assert.New(t)
	assert.Equal(http.StatusOK, w.Code, "http Code doesn't match")
	assert.JSONEq(expectedResp, w.Body.String(), "Response body doesn't match")

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

//Return a list of all employees that attend a specified event and need accommodation
func TestGetEmployeesForEventAcc(t *testing.T) {
	//Init mock db
	db, mock := newMock()

	h := handlers.New(db)
	//Init router
	router := gin.Default()
	router.GET("/events/:id/employees", h.GetEmployeesForEvent)

	//http request
	req, _ := http.NewRequest("GET", "/events/1/employees", nil)
	q := req.URL.Query()
	q.Add("accommodation", "true")
	req.URL.RawQuery = q.Encode()

	// Event to test
	event := models.Event{ID: 1, Name: "Costume Party", Date: "2022-08-01"}

	query := `SELECT id, first_name, last_name, birthday, gender FROM employees JOIN attendances 
		ON attendances.employee_id = employees.id WHERE attendances.event_id = $1 AND accommodation = true ORDER BY id`

	//row for select query
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "birthday", "gender"}).AddRow(
		1, "Son", "Nong", "1999-05-19", "m").AddRow(2, "Max", "Mustermann", "1998-04-18", "m")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(strconv.Itoa(event.ID)).WillReturnRows(rows)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expectedResp := `[
		{
			"id": 1,
			"firstName": "Son",
			"lastName": "Nong",
			"birthDay": "1999-05-19",
			"gender": "m"
		},
		{
			"id": 2,
			"firstName": "Max",
			"lastName": "Mustermann",
			"birthDay": "1998-04-18",
			"gender": "m"
		}
	]`

	assert := assert.New(t)
	assert.Equal(http.StatusOK, w.Code, "http Code doesn't match")
	assert.JSONEq(expectedResp, w.Body.String(), "Response body doesn't match")

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

//Return a list of all employees that attend a specified event and don't need accommodation
func TestGetEmployeesForEventNoAcc(t *testing.T) {
	//Init mock db
	db, mock := newMock()

	h := handlers.New(db)
	//Init router
	router := gin.Default()
	router.GET("/events/:id/employees", h.GetEmployeesForEvent)

	//http request
	req, _ := http.NewRequest("GET", "/events/1/employees", nil)
	q := req.URL.Query()
	q.Add("accommodation", "false")
	req.URL.RawQuery = q.Encode()

	// Event to test
	event := models.Event{ID: 1, Name: "Costume Party", Date: "2022-08-01"}

	query := `SELECT id, first_name, last_name, birthday, gender FROM employees JOIN attendances 
		ON attendances.employee_id = employees.id WHERE attendances.event_id = $1 AND accommodation = false ORDER BY id`

	//row for select query
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "birthday", "gender"}).AddRow(
		3, "Joe", "Jones", "1997-09-12", "m")

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(strconv.Itoa(event.ID)).WillReturnRows(rows)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expectedResp := `[
		{
			"id": 3,
			"firstName": "Joe",
			"lastName": "Jones",
			"birthDay": "1997-09-12",
			"gender": "m"
		}
	]`

	assert := assert.New(t)
	assert.Equal(http.StatusOK, w.Code, "http Code doesn't match")
	assert.JSONEq(expectedResp, w.Body.String(), "Response body doesn't match")

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
