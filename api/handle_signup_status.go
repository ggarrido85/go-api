package api

import (
	"net/http"

	//"github.com/checkmarble/marble-backend/usecases"
	"github.com/gin-gonic/gin"
)

func handleSignupStatus( ) func(c *gin.Context) {
	return func(c *gin.Context) {


		c.JSON(http.StatusOK, gin.H{
			"migrations_run":      "   dddd",
			"has_an_organization": "   dddd",
			"has_a_user":          "   dddd",
		})
	}
}
