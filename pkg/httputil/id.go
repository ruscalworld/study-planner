package httputil

import (
	"fmt"
	"strconv"

	"study-planner/pkg/stderrors"

	"github.com/gofiber/fiber/v2"
)

func ExtractId(ctx *fiber.Ctx, paramName string) (int64, error) {
	idParam := ctx.Params(paramName)
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, stderrors.BadRequest(fmt.Sprintf("invalid %s", paramName))
	}

	return id, nil
}
