// from: https://github.com/chenghonour/formulago

package middleware

import (
	"context"
	"fmt"
	Data "go-social-network/data"
	"strconv"
	"time"

	"go-social-network/biz/domain"
	"go-social-network/biz/logic"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/jwt"
)

type jwtLogin struct {
	Username string `form:"username,required" json:"username,required"` //lint:ignore SA5008 ignoreCheck
	Password string `form:"password,required" json:"password,required"` //lint:ignore SA5008 ignoreCheck
}

// jwt identityKey
var (
	identityKey   = "jwt-id"
	jwtMiddleware = new(jwt.HertzJWTMiddleware)
)

// GetJWTMiddleware returns a new JWT middleware.
func GetJWTMiddleware(d *Data.Data) *jwt.HertzJWTMiddleware {
	jwtMiddleware, err := newJWT(d)
	if err != nil {
		hlog.Fatal(err, "JWT Init Error")
	}
	return jwtMiddleware
}

func newJWT(db *Data.Data) (jwtMiddleware *jwt.HertzJWTMiddleware, err error) {
	// the jwt middleware
	jwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:       "go-social-network",
		Key:         []byte("test"),
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			// take map which have roleID, userID as Payload
			if v, ok := data.(map[string]interface{}); ok {
				return jwt.MapClaims{
					identityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			payloadMap, ok := claims[identityKey].(map[string]interface{})
			if !ok {
				hlog.Error("get payloadMap error", "claims data:", claims[identityKey])
				return nil
			}
			fmt.Println("payload userID:", payloadMap["userID"])
			c.Set("userID", payloadMap["userID"])
			return payloadMap
		},
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			res := new(domain.UserLoginResp)
			// normal jwtLogin
			var loginVal jwtLogin
			if err := c.BindAndValidate(&loginVal); err != nil {
				return "", err
			}

			// Login
			username := loginVal.Username
			password := loginVal.Password
			res, err = logic.NewUser(db).Login(ctx, username, password)
			if err != nil {
				hlog.Error(err, "jwtLogin error")
				return nil, err
			}
			// return the payload
			// take str roleID, userID into PayloadMap
			payloadMap := make(map[string]interface{})
			payloadMap["userID"] = strconv.Itoa(int(res.UserID))
			return payloadMap, nil
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, map[string]interface{}{
				"code":    code,
				"message": message,
			})
		},
	})

	return
}
