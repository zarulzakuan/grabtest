package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/zarulzakuan/grabtest/docs"
)

// NewRoutes All routes defintion
func NewRoutes(route *gin.Engine) {

	request := route.Group("/request")
	{
		_ = request
		// request.GET("", sendPaymentRequest)
	}

}
