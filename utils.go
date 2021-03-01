package main

import (
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
		log.Println("Set CORS")
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
		log.Println(c.Request.Header.Get("Authorization"))
		Token := UserToken{}
		if os.Getenv("ENV") != "PROD" {
			Token = UserToken{
				"UETk9f0woXbDT6Mb7n1sggrLnL53",
				"ahmed",
			}
			c.Set("Token", Token)
			return
		}
		dbApp := c.MustGet("dbapp").(*firebase.App)
		authClient, err := dbApp.Auth(c)

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
			token, err := authClient.VerifyIDToken(c, strArr[1])

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

			// all good here. set context

			user, err := authClient.GetUser(c, token.UID)
			if err != nil {
				log.Printf("error getting user %s: %v\n", token.UID, err)
			}

			Token = UserToken{
				token.UID,
				user.DisplayName,
			}

			c.Set("Token", Token)

		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
			return
		}

	}
}

func dbInit() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Initialize DB")
		sa := option.WithCredentialsJSON([]byte(`{
			"type": "service_account",
			"project_id": "grabtest-828e6",
			"private_key_id": "b52107f8f5c6a14a693ff9224c9e3e4ff9f7bf5a",
			"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQC1Uu8pwBCaMWJf\nJBhnQVkyO8nX2dGw3Uv9nC2GIlptTQ8mg45MWMzkz3JOVuL20ZUM+XtXo4VZzHN3\n0ngKD5L1YmupRSeQ/jqNMg9CIFJ9Op/Xl9qRqqkvV7WZoNwXoPFoKEl8h41gzApM\npUlBNozz/gRTiuWaCOwXWrY3RkJNeZpckQs1gk0vki3p8y+g9ASc7tm0qqDsTnnK\nZFU2JpvrFqIvwlZtvK3Q2dfyegM+zVc/0YbnccRWzzIbEHU5UOIx8phIv9swtR1m\nemITU5E1j2o2nf3yiQgqexbCdmhxxd+1l10Z3iN+1zxoTk8zpILfsM0OYu4QyJ3f\nvnSb555BAgMBAAECggEAAYsU9m6bdZFQCTreBLSQrsj4sDnUlN3KHmMS6Osn7xNs\nOpawdZCxhcA72wc9FuUa9w/nvD8Fc1ZpisgjnDcchOXAJtVFpfB6ZrNALMu1Ozpu\nTwzZOZfE3HgYCpiAGkKVBQVUXGP6KkW/HLHkP4YLiXp/zO1uUNmC/bV+8YkYJvkL\nUHFzxd2ToU44p4W8MIwfo68fTQjxN5ciT7SL5b/78/KjoZ1jMqCaQqk+gqHTUJRL\ne6fQK8XQTrcERUv1baeUoKyua+1TUdTiiOFwq1/vLqpptxoLiQTwtCfKSrcBE+vP\n8TtwNsjqzzt4yo4LcIjrjf4oty0V6em+LPAYRPi2RQKBgQDlDHWIRsuccoCVu4C6\nXaFhXDQGGRizcMWCWal8c4/bTLo3j/d8pn6hvsIimEFG7puhOxc3IuEg78IcCppg\nVxoQBsRb1ltSOTAAEKtgl7jJp7F+AOXyYxG97tcP/DZIdSeFERKMV+VqhX3cJmZb\nm0XiiWfwxzhvYhpJ7co0fSMQ9wKBgQDKqOILRvI3uQGkO01qZN2c0vewkFoYkAbD\nSTvmqvYqe5OXZHzqeiOp07hCw7eFVFH+AQMWlvPq2OeNEBSiyJnS4LCIWkbGHkC/\nRGFXNQRbVbGwZwcmhFMpDZx6YWwA5464hcc/51x//wCyZYD3dbEajGNWeqjBaeen\nzapD2mS0hwKBgQCWHXxKPFvlxQWRHLpZalQCQzO1a21M7XQE4k66SeLWj4rcL3a8\nM0J7L1J86dyeaHOHT/r/H9T8iSZmymwzB+ME7epzZiGj2ecjo8kuHUH2p/kj4+LQ\n4S0XlhlNWLca9e1YwL+vS0wIbET4rBIZp8I9nmCI5YiPN3STT01e6US6MQKBgQC3\n0kn+uqJ4ArrHcfb9e3I8jmuW3siqIPHRbvsDdq0EycSM/NwFfzYcE+u9u3MaX+pj\nQB3B/rhOm+Ij5KMjKFvGmIHnnCM4Dzbhhq/Yf0FtayRagolM2Zfo0+zMYNOrWl3t\nZ4LcpwoTG4VRS5qYW0uCbjaouWea0GoMnhNDqrPPvwKBgQDSl0xzUQvqNzuyl7rS\nTCensBxI+kGtJxQf6v1TUFCMPcb42upm2gF9Cgn930d9Bq0I185gmj7ysQZDLSpc\nltdCZeBA6Oli7EtsKmi55GUr5gFclueVevsoMYzDyoxw/6KhSwoi5bJF+0eEai5S\nG6vETnhN9mbO+o7yXTiIHqQfrg==\n-----END PRIVATE KEY-----\n",
			"client_email": "firebase-adminsdk-smde2@grabtest-828e6.iam.gserviceaccount.com",
			"client_id": "105092542092730027443",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-smde2%40grabtest-828e6.iam.gserviceaccount.com"
		  }
	  `))
		app, err := firebase.NewApp(c, nil, sa)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Set("dbapp", app)

	}
}

func timeNowEpoch() int64 {
	now := time.Now()
	secs := now.Unix()
	return secs
}
