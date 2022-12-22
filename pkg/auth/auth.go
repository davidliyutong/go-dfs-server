package auth

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go-dfs-server/pkg/nameserver/server"
	"time"
)

type ClientLoginCredential struct {
	AccessKey string `form:"accessKey" json:"accessKey" binding:"required"`
	SecretKey string `form:"secretKey" json:"secretKey" binding:"required"`
}

type JWTResponse struct {
	Code   int       `json:"code"`
	Expire time.Time `json:"expire"`
	Token  string    `json:"token"`
}

type User struct {
	AccessKey string
}

func RegisterAuthModule(engine *gin.Engine, loginPath string, tokenRefreshPath string, timeout time.Duration, authnFn func(string, string) bool, authzFn func(string) bool) (*jwt.GinJWTMiddleware, error) {
	ginJWT, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            server.GlobalServerDesc.Opt.Auth.Domain,
		SigningAlgorithm: "HS256",
		Key:              []byte(server.GlobalServerDesc.Opt.Auth.SecretKey),
		Timeout:          timeout,
		MaxRefresh:       timeout,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginValues ClientLoginCredential
			if err := c.ShouldBind(&loginValues); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			accessKey := loginValues.AccessKey
			secretKey := loginValues.SecretKey

			if authnFn(accessKey, secretKey) {
				return &User{
					AccessKey: accessKey,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					IdentityKeyStr: v.AccessKey,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				AccessKey: claims[IdentityKeyStr].(string),
			}
		},
		IdentityKey: IdentityKeyStr,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && authzFn(v.AccessKey) {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.IndentedJSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		SendCookie:    true,
		TimeFunc:      time.Now,
	})

	authStrategy := NewJWTStrategy(*ginJWT)

	engine.POST(loginPath, authStrategy.LoginHandler)
	engine.POST(tokenRefreshPath, authStrategy.RefreshHandler)

	return ginJWT, nil
}

func CreateJWTAuthGroup(ginEngine *gin.Engine, ginJWT *jwt.GinJWTMiddleware, relativePath string) (*gin.RouterGroup, error) {
	authStrategy := NewJWTStrategy(*ginJWT)
	authGroup := ginEngine.Group(relativePath)
	authGroup.Use(authStrategy.AuthFunc())

	return authGroup, nil
}

func MakeJWTAuthGroup(authGroup *gin.RouterGroup, ginJWT *jwt.GinJWTMiddleware) error {
	authStrategy := NewJWTStrategy(*ginJWT)
	authGroup.Use(authStrategy.AuthFunc())

	return nil
}
