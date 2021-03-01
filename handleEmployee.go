package main

import (
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
)

// @Summary Create new employee
// @Tags employee
// @Accept  json
// @Param data body Employee true "The input Employee struct"
// @Produce  json
// @Success 200 {object} HTTPResponse
// @Router /employee/api [post]
func createEmployee(c *gin.Context) {
	employeeProfile := new(Employee)
	resp := new(HTTPResponse)

	dbApp := c.MustGet("dbapp").(*firebase.App)
	dbClient, _ := dbApp.Firestore(c)
	defer dbClient.Close()

	if err := c.ShouldBindJSON(employeeProfile); err != nil {
		resp.Status = "error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	employeeProfile.ID = guuid.New().String()
	employeeProfile.SearchKeys = strings.Split(employeeProfile.Name, " ")

	//create patient
	_, err := dbClient.Collection("employees").Doc(employeeProfile.ID).Set(c, employeeProfile)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		resp.Status = "error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	listEmployees(c)

}

// @Summary List all employee
// @Tags employee
// @Accept  json
// @Produce  json
// @Success 200 {object} HTTPResponse
// @Router /employee/api [get]
func listEmployees(c *gin.Context) {
	employeeProfiles := []*Employee{}
	resp := new(HTTPResponse)

	dbApp := c.MustGet("dbapp").(*firebase.App)
	dbClient, _ := dbApp.Firestore(c)
	defer dbClient.Close()

	queryResult := dbClient.Collection("employees").Documents(c)
	for {
		employeeProfile := new(Employee)
		doc, err := queryResult.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			resp.Status = "error"
			resp.Data = err.Error()
			c.JSON(http.StatusBadRequest, resp)

			return
		}

		err = mapstructure.Decode(doc.Data(), employeeProfile)
		if err != nil {
			resp.Status = "fail"
			resp.Data = "No employee found"
			c.JSON(http.StatusBadRequest, resp)
		}
		employeeProfiles = append(employeeProfiles, employeeProfile)

	}

	resp.Status = "success"
	resp.Data = employeeProfiles
	c.JSON(http.StatusOK, resp)

}

// @Summary Update employee
// @Tags employee
// @Accept  json
// @Param id path string true "Employee ID"
// @Produce  json
// @Success 200 {object} HTTPResponse
// @Router /employee/api/{id} [put]
func updateEmployee(c *gin.Context) {

	var employeeProfile map[string]interface{}
	resp := new(HTTPResponse)

	id := c.Param("id")

	dbApp := c.MustGet("dbapp").(*firebase.App)
	dbClient, _ := dbApp.Firestore(c)
	defer dbClient.Close()

	if err := c.ShouldBindJSON(&employeeProfile); err != nil {
		resp.Status = "error"
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	_, err := dbClient.Collection("employees").Doc(id).Set(c, employeeProfile, firestore.MergeAll)
	if err != nil {
		resp.Status = "error"
		resp.Data = err.Error()
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Status = "success"
	resp.Data = employeeProfile
	c.JSON(http.StatusOK, resp)

}

// @Summary Delete employee
// @Tags employee
// @Accept  json
// @Param id path string true "Employee ID"
// @Produce  json
// @Success 200 {object} HTTPResponse
// @Router /employee/api/{id} [delete]
func deleteEmployee(c *gin.Context) {

	var employeeProfile map[string]interface{}
	resp := new(HTTPResponse)

	id := c.Param("id")

	dbApp := c.MustGet("dbapp").(*firebase.App)
	dbClient, _ := dbApp.Firestore(c)
	defer dbClient.Close()

	if err := c.ShouldBindJSON(&employeeProfile); err != nil {
		resp.Status = "error"
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	_, err := dbClient.Collection("employees").Doc(id).Delete(c)
	if err != nil {
		resp.Status = "error"
		resp.Data = err.Error()
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Status = "success"
	resp.Data = employeeProfile
	c.JSON(http.StatusOK, resp)

}

// @Summary Find employee
// @Tags employee
// @Accept  json
// @Param searchstring path string true "Search string"
// @Produce  json
// @Success 200 {object} HTTPResponse
// @Router /employee/api/search/{searchstring} [get]
func searchEmployee(c *gin.Context) {
	employeeProfiles := []*Employee{}
	resp := new(HTTPResponse)

	searchString := c.Param("searchstring")

	dbApp := c.MustGet("dbapp").(*firebase.App)
	dbClient, _ := dbApp.Firestore(c)
	defer dbClient.Close()

	queryResult := dbClient.Collection("employees").Where("searchkeys", "array-contains", searchString).Documents(c)
	for {
		employeeProfile := new(Employee)
		doc, err := queryResult.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			resp.Status = "error"
			resp.Data = err.Error()
			c.JSON(http.StatusBadRequest, resp)

			return
		}

		err = mapstructure.Decode(doc.Data(), employeeProfile)
		if err != nil {
			resp.Status = "fail"
			resp.Data = "No employee found"
			c.JSON(http.StatusBadRequest, resp)
		}
		employeeProfiles = append(employeeProfiles, employeeProfile)

	}

	resp.Status = "success"
	resp.Data = employeeProfiles
	c.JSON(http.StatusOK, resp)

}
