package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/zarulzakuan/grabtest/docs"
)

func PublicRoutes(route *gin.Engine) {

	access := route.Group("/access")
	{
		_ = access
	}

	route.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	route.GET("/employees", func(c *gin.Context) {
		c.HTML(http.StatusOK, "employee.html", nil)
	})
	route.GET("/jobs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "job.html", nil)
	})
	route.GET("/reports", func(c *gin.Context) {
		c.HTML(http.StatusOK, "report.html", nil)
	})
	route.GET("/log", func(c *gin.Context) {
		c.HTML(http.StatusOK, "log.html", nil)
	})
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func ProtectedRoutes(route *gin.Engine) {

	job := route.Group("/job/api")
	{
		job.POST("", createJob)
		job.GET("", listJobs)
		// employee.GET("/search/:searchstring", searchJob)
		job.PUT("/:id", updateJob)
		job.DELETE("/:id", deleteJob)
	}

	employee := route.Group("/employee/api")
	{
		employee.POST("", createEmployee)
		employee.GET("", listEmployees)
		// employee.GET("/search/:searchstring", searchEmployee)
		employee.PUT("/:id", updateEmployee)
		employee.DELETE("/:id", deleteEmployee)
	}

	log := route.Group("/log/api")
	{
		log.GET("", listLogs)
		log.GET("/in", recordLogIn)
		log.GET("/out", recordLogOut)
		log.GET("/visit/:page", recordPageVisit)

	}

}
