package httputil

import (
	"errors"
	"log"

	"study-planner/pkg/stderrors"

	"github.com/gofiber/fiber/v2"
)

type errorResponse struct {
	Error string `json:"error"`
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	err = unwrapError(err)
	var fiberErr *fiber.Error

	if errors.As(err, &fiberErr) {
		return ctx.Status(fiberErr.Code).JSON(errorResponse{Error: fiberErr.Message})
	}

	log.Println("error while processing request:", err)
	return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "unknown error"})
}

func unwrapError(err error) error {
	var e stderrors.WebError
	if errors.As(err, &e) {
		return fiber.NewError(e.GetStatusCode(), e.Error())
	}

	return err
}
