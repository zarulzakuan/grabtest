package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

func setCORS() gin.HandlerFunc {
	return func(c *gin.Context) {

		// First, we add the headers with need to enable CORS
		// Make sure to adjust these headers to your needs
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Content-Type", "application/json")

		// Second, we handle the OPTIONS problem
		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusOK)
		}
	}
}

func authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Enter auth")

		Token := UserToken{}
		if os.Getenv("ENV") != "PROD" {
			Token = UserToken{
				"UETk9f0woXbDT6Mb7n1sggrLnL53",
				"admin",
				true,
			}
			c.Set("Token", Token)
			return
		}
		app, ctx := createFirestoreClient()
		authClient, err := app.Auth(ctx)

		if err != nil {

			log.Println(err.Error())
		}
		if err != nil {
			log.Printf("error verifying ID token: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
			return
		}

		bearToken := c.Request.Header.Get("Authorization")
		strArr := strings.Split(bearToken, " ")
		if len(strArr) == 2 {
			token, err := authClient.VerifyIDToken(ctx, strArr[1])

			if err != nil {

				log.Printf("error verifying ID token: %v\n", err)
				matched, _ := regexp.MatchString(`expired`, err.Error())
				if matched {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Expired"})
					return
				}
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
				return
			}

			myclaims := token.Claims

			var role string
			var active bool

			if myclaims["role"] != nil {
				role = myclaims["role"].(string)
			}
			if myclaims["active"] != nil {
				active = myclaims["active"].(bool)
			}

			UID := token.UID
			log.Println(myclaims)
			log.Println(role)
			if role != "admin" || !active || UID == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
				return
			}

			// all good here. set context

			Token = UserToken{
				UID,
				role,
				active,
			}

			c.Set("Token", Token)

		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
			return
		}

	}
}

func createFirestoreClient() (*firebase.App, context.Context) {
	// Get a Firestore client.
	config := &firebase.Config{
		StorageBucket: "foresight-774f4.appspot.com",
		// ProjectID:     "foresightinfra", // comment this out to test on local
	}
	ctx := context.Background()
	sa := option.WithCredentialsFile(firebasePermissionFile)
	app, err := firebase.NewApp(ctx, config, sa)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// return
	}
	return app, ctx
}

func timeNowEpoch() int64 {
	now := time.Now()
	secs := now.Unix()
	return secs
}
