package auth

import (
	"github.com/gin-gonic/gin"
)

func M2MAuthorizationRequired(m2mAuthToken string) func(*gin.Context) {

	return gin.BasicAuth(gin.Accounts{
		"dev": m2mAuthToken,
	})

}
