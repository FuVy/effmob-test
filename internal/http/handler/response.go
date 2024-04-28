package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type messageResponse struct {
	Error   *CodeResponse `json:"error"`
	Data    interface{}   `json:"data"`
	Message string        `json:"message"`
} // @name ApiResponse

type CodeResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
} // @name ApiCodeResponse

func badRequest(ctx *fiber.Ctx, code int, message string) error {
	ctx.Response().SetStatusCode(http.StatusOK)
	return ctx.JSON(messageResponse{
		Error: &CodeResponse{
			Code:    code,
			Message: message,
		},
	})
}

func respondMessage(ctx *fiber.Ctx, message string) error {
	return ctx.JSON(messageResponse{
		Message: message,
		Data:    CodeResponse{},
	})
}

func respondJSON(ctx *fiber.Ctx, statusCode int, msg, errMsg string, data interface{}) error {
	resp := messageResponse{
		Data:    data,
		Message: msg,
	}
	if statusCode >= 400 {
		resp.Error = &CodeResponse{
			Code:    statusCode,
			Message: errMsg,
		}
	}
	return ctx.JSON(resp)
}

func writeApiError(ctx *fiber.Ctx, codeResponse *CodeResponse) error {
	ctx.Response().SetStatusCode(http.StatusBadRequest)
	return ctx.JSON(*codeResponse)
}

func writeInvalidRequest(ctx *fiber.Ctx) error {
	ctx.Response().SetStatusCode(http.StatusUnprocessableEntity)
	return ctx.SendString(("can't parse json body"))
}

func internalError(ctx *fiber.Ctx, err error) error {
	ctx.Context().SetStatusCode(http.StatusInternalServerError)
	return err
}
