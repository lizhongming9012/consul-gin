package app

import (
	"NULL/consul-gin/pkg/e"
	"NULL/consul-gin/pkg/logging"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// BindAndValid binds and validates data
func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		logging.Error("BIND:%v", err)
		return http.StatusBadRequest, e.INVALID_PARSE_FORM
	}
	logging.Debug(fmt.Sprintf("%+v", form))
	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		log.Printf("VERIFY: %v", err)
		return http.StatusInternalServerError, e.ERROR
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, e.INVALID_PARAMS_VERIFY
	}

	return http.StatusOK, e.SUCCESS
}
