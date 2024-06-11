package handlers

import (
	"Shaunak/Employee_manager/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
)

// RegisterEmployeeHandlers registers employee-related HTTP handlers
func RegisterEmployeeHandlers(r *gin.Engine) {
	r.GET("/employee/:id", GetEmployeeById)
	r.GET("/employees", GetPaginatedEmployees)
	r.POST("/insertEmployee", InsertEmployee)
	r.PUT("/updateEmployee/:id", UpdateEmployee)
	r.DELETE("/deleteEmployee/:id", DeleteEmployee)
}

func GetEmployeeById(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	id, _ := strconv.Atoi(c.Param("id"))

	var employee models.Employee
	err := db.QueryRow("SELECT id, name, position, salary FROM employees WHERE id = $1", id).
		Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Salary)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func GetPaginatedEmployees(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	var requestBody struct {
		First int `json:"first"`
		After int `json:"after"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rows, err := db.Query("SELECT id, name, position, salary FROM employees ORDER BY id LIMIT $1 OFFSET $2", requestBody.First, requestBody.After)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employees"})
		return
	}
	defer rows.Close()

	employees := make([]models.Employee, 0)
	for rows.Next() {
		var employee models.Employee
		if err := rows.Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Salary); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employees"})
			return
		}
		employees = append(employees, employee)
	}

	c.JSON(http.StatusOK, employees)
}

func DeleteEmployee(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	id, _ := strconv.Atoi(c.Param("id"))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := db.Exec("DELETE FROM employees WHERE id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
			return
		}
	}()
	wg.Wait()
	c.JSON(http.StatusOK, gin.H{"Message": "Employee deleted successfully"})
	c.Status(http.StatusOK)
}

func InsertEmployee(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id int64
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := db.QueryRow("INSERT INTO employees (name, position, salary) VALUES ($1, $2, $3) RETURNING id",
			employee.Name, employee.Position, employee.Salary).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
			return
		}
	}()
	wg.Wait()

	c.JSON(http.StatusOK, gin.H{"id": id, "message": "Employee created successfully"})
	c.Status(http.StatusCreated)
}

func UpdateEmployee(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	id, _ := strconv.Atoi(c.Param("id"))

	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Construct the update query dynamically
	query := "UPDATE employees SET "
	var args []interface{}
	argCount := 1

	// Update the name field if provided
	if employee.Name != "" {
		query += "name = $" + strconv.Itoa(argCount) + ", "
		args = append(args, employee.Name)
		argCount++
	}

	// Update the position field if provided
	if employee.Position != "" {
		query += "position = $" + strconv.Itoa(argCount) + ", "
		args = append(args, employee.Position)
		argCount++
	}

	// Update the salary field if provided
	if employee.Salary != 0 {
		query += "salary = $" + strconv.Itoa(argCount) + ", "
		args = append(args, employee.Salary)
		argCount++
	}

	// Remove the trailing comma and space
	query = query[:len(query)-2]

	// Add the WHERE clause
	query += " WHERE id = $" + strconv.Itoa(argCount)
	args = append(args, id)

	// Execute the update query
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := db.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
			return
		}
	}()
	wg.Wait()
	c.JSON(http.StatusOK, gin.H{"Message": "Employee Updated Successfully"})
	c.Status(http.StatusOK)
}
