package di

import (
	"fmt"
	"golang-block-chain/web/cores"
	"golang-block-chain/web/middlewares"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

func Di(fn interface{}, injectKeys []string, dependencies []gin.HandlerFunc) []gin.HandlerFunc {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("fn is not a function")
	}
	fnValue := reflect.ValueOf(fn)
	numIn := fnType.NumIn()
	handler := func(c *gin.Context) {
		in := make([]reflect.Value, numIn)
		for i, injectKey := range injectKeys {
			in[i] = reflect.ValueOf(c.MustGet(injectKey))
		}
		out := fnValue.Call(in)
		response := convertToResponse(out)
		cores.GenerateGinResponse(c, response)
	}
	result := make([]gin.HandlerFunc, 0)
	result = append(result, dependencies...)
	return append(result, handler)
}

func AutoDi(fn interface{}) []gin.HandlerFunc {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("fn is not a function")
	}

	numIn := fnType.NumIn()
	injectKeys := make([]string, numIn)
	dependencies := make([]gin.HandlerFunc, numIn)
	var pathParamCount int = 0
	for i := 0; i < numIn; i++ {
		injectKeys[i] = fmt.Sprint(i)
		t := fnType.In(i)
		ordinal := fmt.Sprint(i)
		if t == reflect.TypeOf((uint)(0)) {
			dependencies[i] = middlewares.InjectFromPathParamSeq(ordinal, pathParamCount, true)
			pathParamCount++
		} else if t == reflect.TypeOf((string)("")) {
			dependencies[i] = middlewares.InjectFromPathParamSeq(ordinal, pathParamCount, false)
			pathParamCount++
		} else if t == reflect.TypeOf((*gin.Context)(nil)) {
			dependencies[i] = middlewares.InjectGinContext(ordinal)
		} else if t == reflect.TypeOf((*cores.UserInfo)(nil)) {
			dependencies[i] = middlewares.InjectUserInfo(ordinal)
		} else {
			dependencies[i] = middlewares.InjectFromGinBind(ordinal, t)
		}
	}
	return Di(fn, injectKeys, dependencies)
}

func convertToResponse(fnOut []reflect.Value) *cores.Response {
	var response *cores.Response
	var body interface{}
	var err error
L:
	for _, outValue := range fnOut {
		switch v := outValue.Interface().(type) {
		case *cores.Response:
			response = v
		case cores.Response:
			response = &v
		case error:
			err = v
			break L
		default:
			if body == nil {
				body = v
			}
		}
	}
	if err != nil {
		return cores.NewErrorResponse(http.StatusInternalServerError, err.Error(), nil)
	}
	if response != nil {
		return response
	}
	if body == nil {
		return cores.NewResponse(http.StatusNoContent, nil)
	}
	return cores.NewResponse(http.StatusOK, body)
}
