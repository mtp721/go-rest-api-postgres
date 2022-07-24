# Assignment

Rest API written in GO using Postgres, Gin, lib/pq and database/sql. With tests using httptest and sqlmock.

## Endpoints:

• POST /employees --registers a new employee in the system--

• GET /employees --returns the list of all micobo employees--

• PUT /employees/{employee_id} --update the specified employee's information--

• DELETE /employees/{employee_id} --delete the specified employee from the system--

• GET /events --returns a list with all upcoming events--

• GET /events/{event_id} --returns the specific event--

• GET /events/{event_id}/employees --returns the list of the employees that are assisting to the event, should accept query parameters for filtering if they need or not accommodation--
