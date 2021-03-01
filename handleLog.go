package main

import (
	"log"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
)

func recordLogIn(c *gin.Context) {
	recordLog := new(RecordLog)
	resp := new(HTTPResponse)

	dbApp := c.MustGet("dbapp").(*firebase.App)
	dbClient, _ := dbApp.Firestore(c)
	defer dbClient.Close()

	userToken := c.MustGet("Token").(UserToken)

	recordLog.DateTime = time.Now().Format("2006.01.02 15:04:05")
	recordLog.ID = guuid.New().String()
	recordLog.Type = "login"
	recordLog.UserID = userToken.UID
	recordLog.UserName = userToken.UserName
	//create patient
	_, err := dbClient.Collection("recordlog").Doc(recordLog.ID).Set(c, recordLog)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		resp.Status = "error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
}

func recordLogOut(c *gin.Context) {
	recordLog := new(RecordLog)
	resp := new(HTTPResponse)

	dbApp := c.MustGet("dbapp").(*firebase.App)
	dbClient, _ := dbApp.Firestore(c)
	defer dbClient.Close()

	userToken := c.MustGet("Token").(UserToken)

	recordLog.DateTime = time.Now().Format("2006.01.02 15:04:05")
	recordLog.ID = guuid.New().String()
	recordLog.Type = "logout"
	recordLog.UserID = userToken.UID
	recordLog.UserName = userToken.UserName
	//create patient
	_, err := dbClient.Collection("recordlog").Doc(recordLog.ID).Set(c, recordLog)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		resp.Status = "error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
}

func recordPageVisit(c *gin.Context) {
	dbApp := c.MustGet("dbapp").(*firebase.App)
	dbClient, _ := dbApp.Firestore(c)
	defer dbClient.Close()
	recordLog := new(RecordLog)
	userToken := c.MustGet("Token").(UserToken)

	pageName := c.Param("page")

	recordLog.DateTime = time.Now().Format("2006.01.02 15:04:05")
	recordLog.ID = guuid.New().String()
	recordLog.Type = "page visit"
	recordLog.UserID = userToken.UID
	recordLog.UserName = userToken.UserName
	recordLog.PageName = pageName
	//create patient
	_, err := dbClient.Collection("recordlog").Doc(recordLog.ID).Set(c, recordLog)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func listLogs(c *gin.Context) {
	recordLogs := []*RecordLog{}
	resp := new(HTTPResponse)

	dbApp := c.MustGet("dbapp").(*firebase.App)
	dbClient, _ := dbApp.Firestore(c)
	defer dbClient.Close()

	queryResult := dbClient.Collection("recordlog").Documents(c)
	for {
		recordLog := new(RecordLog)
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

		err = mapstructure.Decode(doc.Data(), recordLog)
		if err != nil {
			resp.Status = "fail"
			resp.Data = "No job found"
			c.JSON(http.StatusBadRequest, resp)
		}
		recordLogs = append(recordLogs, recordLog)

	}

	resp.Status = "success"
	resp.Data = recordLogs
	c.JSON(http.StatusOK, resp)
}
