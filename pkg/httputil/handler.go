package httputil

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

type (
	Endpoint[I, O any]    func(ctx *fiber.Ctx, request I) (*O, error)
	SimpleEndpoint[O any] func(ctx *fiber.Ctx) (*O, error)
)

func MakeHandler[I, O any](endpoint Endpoint[I, O]) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var request I
		err := json.Unmarshal(ctx.Body(), &request)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		response, err := endpoint(ctx, request)
		if err != nil {
			return unwrapError(err)
		}

		return handleResponse(ctx, response)
	}
}

func MakeSimpleHandler[O any](endpoint SimpleEndpoint[O]) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		response, err := endpoint(ctx)
		if err != nil {
			return unwrapError(err)
		}

		return handleResponse(ctx, response)
	}
}

func handleResponse[O any](ctx *fiber.Ctx, response *O) error {
	if response != nil {

		return ctx.JSON(response)
	} else {
		ctx.Status(fiber.StatusNoContent)
	}

	return ctx.Send(nil)
}
