package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    StatusCode  `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OK(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code:    StatusOK,
		Message: message,
		Data:    data,
	})
}
func Error(ctx *gin.Context, code StatusCode, message string) {
	ctx.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}
func BindError(ctx *gin.Context, message string) {
	Error(ctx, StatusBindError, message)
}

func NotExistError(ctx *gin.Context, message string) {
	Error(ctx, StatusNotExist, message)
}

func AlreadyExisterror(ctx *gin.Context, message string) {
	Error(ctx, StatusAlreadyExist, message)
}

func RegexpCompileError(ctx *gin.Context, message string) {
	Error(ctx, StatusRegCompileError, message)
}

func AddError(ctx *gin.Context, msg string) {
	Error(ctx, StatusAddError, msg)
}

func DelError(ctx *gin.Context, message string) {
	Error(ctx, StatusDelError, message)
}

func ModError(ctx *gin.Context, message string) {
	Error(ctx, StatusModError, message)
}

func GetError(ctx *gin.Context, message string) {
	Error(ctx, StatusGetError, message)
}

func DatabaseAddError(ctx *gin.Context, message string) {
	Error(ctx, StatusAddDiskError, message)
}

func DatabaseDelError(ctx *gin.Context, message string) {
	Error(ctx, StatusDelDiskError, message)
}

func DatabaseModError(ctx *gin.Context, message string) {
	Error(ctx, StatusModDiskError, message)
}

func DatabaseGetError(ctx *gin.Context, message string) {
	Error(ctx, StatusGetDiskError, message)
}

func InvalidURLFormatError(ctx *gin.Context, message string) {
	Error(ctx, StatusInvalidURLFormat, message)
}
