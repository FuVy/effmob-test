package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	log "github.com/sirupsen/logrus"
)

func commonMiddleware(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	err := ctx.Next()
	if err != nil {
		ferr, ok := err.(*fiber.Error)
		if ok && ferr.Code < 500 {
			return badRequest(ctx, ferr.Code, ferr.Message)
		}
		ctx.Context().SetStatusCode(http.StatusInternalServerError)
		log.WithError(err).WithField("request-id", ctx.Locals(requestid.ConfigDefault.ContextKey)).Error("request error:", err)
	}
	return nil
}
