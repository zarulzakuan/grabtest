package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zarulzakuan/grabtest/docs"
)

// @title Grab Assessment
// @version 1.0
// @description Main service
// @termsOfService http://swagger.io/terms/
// @contact.name Zarul Zakuan
// @contact.email zarulzakuan@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /

const firebasePermissionFile = "./foresight-774f4-firebase-adminsdk-2p3wz-cf4c7c749a.json"

var router = gin.Default()
var appPort string
var appHost string
var ipay88ReqURL string
var backendURL string
var responseURL string
var merchantCode string
var merchantKey string
var appEnv string

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
		os.Exit(1)
	}

	var exists bool

	appPort, exists = os.LookupEnv("APP_PORT")

	if !exists || appPort == "" {
		log.Fatal("Environment Path not set")
		os.Exit(1)
	}

	appHost, exists = os.LookupEnv("APP_HOST")

	if !exists || appHost == "" {
		log.Fatal("Environment Path not set")
		os.Exit(1)
	}

	appEnv, exists = os.LookupEnv("ENV")
	if appEnv == "PROD" {
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", appHost, appPort)
	}
}

func main() {

	log.Println("Running...")

	//connect database
	router.Use(dbInit())

	// load html files
	router.LoadHTMLGlob("static/*")

	PublicRoutes(router)
	router.Use(static.Serve("/assets", static.LocalFile("./static", false)))

	router.Use(authenticate())
	ProtectedRoutes(router)

	router.Use(setCORS())

	// router.Use(authenticate())

	s := &http.Server{
		Addr:           fmt.Sprintf(":%v", appPort),
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()

}
