package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Code    StatusCode  `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OK(ctx echo.Context, message string, data interface{}) error {
	return ctx.JSON(http.StatusOK, Response{
		Code:    StatusOK,
		Message: message,
		Data:    data,
	})
}
func Error(ctx echo.Context, code StatusCode, message string) error {
	return ctx.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}
func BindError(ctx echo.Context, message string) error {
	return Error(ctx, StatusBindError, message)
}

func NotExistError(ctx echo.Context, message string) error {
	return Error(ctx, StatusNotExist, message)
}

func AlreadyExisterror(ctx echo.Context, message string) error {
	return Error(ctx, StatusAlreadyExist, message)
}

func RegexpCompileError(ctx echo.Context, message string) error {
	return Error(ctx, StatusRegCompileError, message)
}

func AddError(ctx echo.Context, msg string) error {
	return Error(ctx, StatusAddError, msg)
}

func DelError(ctx echo.Context, message string) error {
	return Error(ctx, StatusDelError, message)
}

func ModError(ctx echo.Context, message string) error {
	return Error(ctx, StatusModError, message)
}

func GetError(ctx echo.Context, message string) error {
	return Error(ctx, StatusGetError, message)
}

func DatabaseAddError(ctx echo.Context, message string) error {
	return Error(ctx, StatusAddDiskError, message)
}

func DatabaseDelError(ctx echo.Context, message string) error {
	return Error(ctx, StatusDelDiskError, message)
}

func DatabaseModError(ctx echo.Context, message string) error {
	return Error(ctx, StatusModDiskError, message)
}

func DatabaseGetError(ctx echo.Context, message string) error {
	return Error(ctx, StatusGetDiskError, message)
}

func InvalidURLFormatError(ctx echo.Context, message string) error {
	return Error(ctx, StatusInvalidURLFormat, message)
}
