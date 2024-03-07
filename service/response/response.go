package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SUCCESS_CODE          = 200
	FAIL_CODE             = 300
	SUPER_ADMIN_ID uint64 = 10086
)

type ResponseModel struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseModelBase struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponsePageData struct {
	Total int64       `json:"total"`
	Items interface{} `json:"items"`
}

type ResponsePage struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    ResponsePageData `json:"data"`
}

// 响应JSON数据
func ResJSON(c *gin.Context, status int, v interface{}) {
	c.JSON(status, v)
	c.Abort()
}

// 响应成功
func ResSuccess(c *gin.Context, v interface{}) {
	result := ResponseModel{
		Code:    SUCCESS_CODE,
		Message: "ok",
		Data:    v,
	}
	ResJSON(c, http.StatusOK, &result)
}

func ResSuccessMsg(c *gin.Context) {
	result := ResponseModelBase{
		Code:    SUCCESS_CODE,
		Message: "ok",
	}
	ResJSON(c, http.StatusOK, &result)
}

func ResFail(c *gin.Context, msg string) {
	result := ResponseModelBase{
		Code:    FAIL_CODE,
		Message: msg,
	}
	ResJSON(c, http.StatusOK, &result)
}

func ResFailCode(c *gin.Context, msg string, code int) {
	result := ResponseModelBase{
		Code:    code,
		Message: msg,
	}
	ResJSON(c, http.StatusOK, &result)
}

func ResErrSrv(c *gin.Context, err error) {
	result := ResponseModelBase{
		Code:    FAIL_CODE,
		Message: "服务器故障...",
	}
	ResJSON(c, http.StatusOK, &result)
}

func ResErrCli(c *gin.Context, err error) {
	result := ResponseModelBase{
		Code:    FAIL_CODE,
		Message: "err",
	}
	ResJSON(c, http.StatusOK, &result)
}

func ResSuccessPage(c *gin.Context, total int64, list interface{}) {
	result := ResponsePage{
		Code:    200,
		Message: "ok",
		Data:    ResponsePageData{Total: total, Items: list},
	}
	ResJSON(c, 200, &result)
}

func ResErrMsg(c *gin.Context, err error) {
	result := ResponseModelBase{
		Code:    FAIL_CODE,
		Message: err.Error(),
	}
	ResJSON(c, http.StatusOK, &result)
}
