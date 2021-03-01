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

// @Summary Create new job
// @Tags job
// @Accept  json
// @Param data body Job true "The input Job struct"
// @Produce  json
// @Success 200 {object} HTTPResponse
// @Router /job/api [post]
func createJob(c *gin.Context) {
	jobProfile := new(Job)
	resp := new(HTTPResponse)

	dbApp := c.MustGet("dbapp").(*firebase.App)
	dbClient, _ := dbApp.Firestore(c)
	defer dbClient.Close()

	if err := c.ShouldBindJSON(jobProfile); err != nil {
		resp.Status = "error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	jobProfile.ID = guuid.New().String()
	jobProfile.SearchKeys = strings.Split(jobProfile.Title, " ")

	//create patient
	_, err := dbClient.Collection("jobs").Doc(jobProfile.ID).Set(c, jobProfile)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		resp.Status = "error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	listJobs(c)

}

// @Summary List all job
// @Tags job
// @Accept  json
// @Produce  json
// @Success 200 {object} HTTPResponse
// @Router /job/api [get]
func listJobs(c *gin.Context) {
	jobProfiles := []*Job{}
	resp := new(HTTPResponse)

	dbApp := c.MustGet("dbapp").(*firebase.App)
	dbClient, _ := dbApp.Firestore(c)
	defer dbClient.Close()

	queryResult := dbClient.Collection("jobs").Documents(c)
	for {
		jobProfile := new(Job)
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

		err = mapstructure.Decode(doc.Data(), jobProfile)
		if err != nil {
			resp.Status = "fail"
			resp.Data = "No job found"
			c.JSON(http.StatusBadRequest, resp)
		}
		jobProfiles = append(jobProfiles, jobProfile)

	}

	resp.Status = "success"
	resp.Data = jobProfiles
	c.JSON(http.StatusOK, resp)

}

// @Summary Update job
// @Tags job
// @Accept  json
// @Param id path string true "Job ID"
// @Produce  json
// @Success 200 {object} HTTPResponse
// @Router /job/api/{id} [put]
func updateJob(c *gin.Context) {

	var jobProfile map[string]interface{}
	resp := new(HTTPResponse)

	id := c.Param("id")

	dbApp := c.MustGet("dbapp").(*firebase.App)
	dbClient, _ := dbApp.Firestore(c)
	defer dbClient.Close()

	if err := c.ShouldBindJSON(&jobProfile); err != nil {
		resp.Status = "error"
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	_, err := dbClient.Collection("jobs").Doc(id).Set(c, jobProfile, firestore.MergeAll)
	if err != nil {
		resp.Status = "error"
		resp.Data = err.Error()
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Status = "success"
	resp.Data = jobProfile
	c.JSON(http.StatusOK, resp)

}

// @Summary Delete job
// @Tags job
// @Accept  json
// @Param id path string true "Job ID"
// @Produce  json
// @Success 200 {object} HTTPResponse
// @Router /job/api/{id} [delete]
func deleteJob(c *gin.Context) {

	var jobProfile map[string]interface{}
	resp := new(HTTPResponse)

	id := c.Param("id")

	dbApp := c.MustGet("dbapp").(*firebase.App)
	dbClient, _ := dbApp.Firestore(c)
	defer dbClient.Close()

	if err := c.ShouldBindJSON(&jobProfile); err != nil {
		resp.Status = "error"
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	_, err := dbClient.Collection("jobs").Doc(id).Delete(c)
	if err != nil {
		resp.Status = "error"
		resp.Data = err.Error()
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Status = "success"
	resp.Data = jobProfile
	c.JSON(http.StatusOK, resp)

}

// @Summary Find job
// @Tags job
// @Accept  json
// @Param searchstring path string true "Search string"
// @Produce  json
// @Success 200 {object} HTTPResponse
// @Router /job/api/search/{searchstring} [get]
func searchJob(c *gin.Context) {
	jobProfiles := []*Job{}
	resp := new(HTTPResponse)

	searchString := c.Param("searchstring")

	dbApp := c.MustGet("dbapp").(*firebase.App)
	dbClient, _ := dbApp.Firestore(c)
	defer dbClient.Close()

	queryResult := dbClient.Collection("jobs").Where("searchkeys", "array-contains", searchString).Documents(c)
	for {
		jobProfile := new(Job)
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

		err = mapstructure.Decode(doc.Data(), jobProfile)
		if err != nil {
			resp.Status = "fail"
			resp.Data = "No job found"
			c.JSON(http.StatusBadRequest, resp)
		}
		jobProfiles = append(jobProfiles, jobProfile)

	}

	resp.Status = "success"
	resp.Data = jobProfiles
	c.JSON(http.StatusOK, resp)

}
