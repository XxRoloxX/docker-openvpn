package rootrouter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"managed-openvpn/internal/auth"
	"os"
)

type RootRouter struct {
	GinRouter *gin.Engine
}

const M2M_AUTH_TOKEN_KEY = "M2M_AUTH_TOKEN"

func NewRootRouter() *RootRouter {

	m2mAuthToken, ok := os.LookupEnv(M2M_AUTH_TOKEN_KEY)
	if !ok {
		panic(fmt.Sprintf("%s is not set", M2M_AUTH_TOKEN_KEY))
	}

	root := gin.Default()
	root.Use(auth.M2MAuthorizationRequired(m2mAuthToken))

	root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &RootRouter{
		GinRouter: root,
	}
}
