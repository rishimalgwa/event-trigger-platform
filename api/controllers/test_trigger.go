package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rishimalgwa/event-trigger-platform/api/views"
)

type TestTriggerController struct{}

func (t *TestTriggerController) TestTrigger(ctx *fiber.Ctx) error {
	return views.OK(ctx, nil)
}
