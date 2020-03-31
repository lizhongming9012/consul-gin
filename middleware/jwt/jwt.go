package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"NULL/consul-gin/pkg/e"
	"NULL/consul-gin/pkg/util"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.GetHeader("X-Access-Token")
		auth := c.GetHeader("Authorization")
		//token := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					expiresAt := claims.(jwt.MapClaims)["exp"].(float64)
					expire := time.Unix(int64(expiresAt), 0).Format("2006-01-02")
					today := time.Now().Format("2006-01-02")
					if expire != today {
						code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
					}
				default:
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
