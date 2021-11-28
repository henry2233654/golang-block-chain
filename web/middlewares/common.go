package middlewares

import (
	"fmt"
	"golang-block-chain/web/cores"
	"net/http"
	"reflect"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func AbortAndResponseError(c *gin.Context, statusCode int, message string, err error) {
	resp := cores.NewErrorResponse(statusCode, message, err)
	cores.GenerateGinResponse(c, resp)
	c.Abort()
}

func GetUserInfo(c *gin.Context) *cores.UserInfo {
	userInfo, exists := c.Get("user_info")
	if exists {
		return userInfo.(*cores.UserInfo)
	}

	userJwt := c.Request.Header.Get("X-USER-JWT")
	if userJwt == "" {
		AbortAndResponseError(c, http.StatusForbidden, "Forbidden", nil)
		return nil
	}

	tokenClaims, err := jwt.ParseWithClaims(
		userJwt, &cores.UserInfo{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(""), nil
		},
	)
	if err != nil {
		AbortAndResponseError(c, http.StatusInternalServerError, err.Error(), err)
		return nil
	}
	userInfo = tokenClaims.Claims.(*cores.UserInfo)
	c.Set("user_info", userInfo)
	return userInfo.(*cores.UserInfo)
}

func InjectionHandler(injectKey string, dependency interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(injectKey, dependency)
	}
}

func InjectFromPathParamName(injectKey string, param string, isUint bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var v interface{}
		v = c.Param(param)
		if isUint {
			var err error
			v, err = strconv.ParseInt(v.(string), 10, 64)
			if err != nil {
				AbortAndResponseError(c, http.StatusBadRequest, fmt.Sprintf("param [%v] is not acceptable, it only allowes type uint", v), nil)
				return
			}
		}
		c.Set(injectKey, v)
	}
}

func InjectFromPathParamSeq(injectKey string, seq int, isUint bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		v := c.Params[seq].Value
		if isUint {
			i64, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				AbortAndResponseError(c, http.StatusBadRequest, fmt.Sprintf("param [%v] is not acceptable, it only allowes type int64", v), nil)
				return
			}
			c.Set(injectKey, i64)
		} else {
			c.Set(injectKey, v)
		}
	}
}

func InjectGinContext(injectKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(injectKey, c)
	}
}

func InjectFromGinBind(injectKey string, v interface{}) gin.HandlerFunc {
	var t reflect.Type
	if _, ok := v.(reflect.Type); ok {
		t = v.(reflect.Type)
	} else {
		t = reflect.TypeOf(v)
	}
	needValidate := true
	if t.Kind() == reflect.Slice {
		needValidate = t.Elem().Kind() == reflect.Struct
	}
	return func(c *gin.Context) {
		modelPtr := reflect.New(t).Interface()
		if err := c.Bind(modelPtr); err != nil {
			AbortAndResponseError(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		model := reflect.ValueOf(modelPtr).Elem().Interface()
		if needValidate {
			if err := validator.New().Var(model, "dive"); err != nil {
				AbortAndResponseError(c, http.StatusBadRequest, err.Error(), nil)
				return
			}
		}
		c.Set(injectKey, model)
	}
}

func InjectUserInfo(injectKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo := GetUserInfo(c)
		if c.IsAborted() {
			return
		}
		c.Set(injectKey, userInfo)
	}
}
