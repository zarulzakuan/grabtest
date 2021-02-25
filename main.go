package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zarulzakuan/grabtest/docs"
)

// @title Ipay88 Payment Gateway Client API
// @version 1.0
// @description Main service
// @termsOfService http://swagger.io/terms/
// @contact.name Zarul Zakuan
// @contact.email zarulzakuan@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8084
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
	ipay88ReqURL, exists = os.LookupEnv("IPAY88_OPSG_REQUEST_URL")

	if !exists || ipay88ReqURL == "" {
		log.Fatal("Environment Path not set")
		os.Exit(1)
	}

	backendURL, exists = os.LookupEnv("BACKEND_URL")

	if !exists || backendURL == "" {
		log.Fatal("Environment Path not set")
		os.Exit(1)
	}

	responseURL, exists = os.LookupEnv("RESPONSE_URL")

	if !exists || responseURL == "" {
		log.Fatal("Environment Path not set")
		os.Exit(1)
	}

	appPort, exists = os.LookupEnv("APP_PORT")

	if !exists || appPort == "" {
		log.Fatal("Environment Path not set")
		os.Exit(1)
	}

	merchantKey, exists = os.LookupEnv("MERCHANT_KEY")

	if !exists || merchantKey == "" {
		log.Fatal("Environment Path not set")
		os.Exit(1)
	}

	merchantCode, exists = os.LookupEnv("MERCHANT_KEY")

	if !exists || merchantCode == "" {
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

	router.Use(setCORS())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// router.Use(authenticate())
	NewRoutes(router)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%v", appPort),
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()

}
