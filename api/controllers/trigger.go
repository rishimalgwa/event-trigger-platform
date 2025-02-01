package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rishimalgwa/event-trigger-platform/api/views"
)

type TriggerController struct{}

func (t *TriggerController) GetAllTriggers(ctx *fiber.Ctx) error {

	return views.OK(ctx, nil)
}

func (t *TriggerController) GetTriggerById(ctx *fiber.Ctx) error {
	return views.OK(ctx, nil)
}

func (t *TriggerController) CreateTrigger(ctx *fiber.Ctx) error {
	return views.OK(ctx, nil)
}

func (t *TriggerController) DeleteTrigger(ctx *fiber.Ctx) error {
	return views.OK(ctx, nil)
}

func (t *TriggerController) UpdateTrigger(ctx *fiber.Ctx) error {
	return views.OK(ctx, nil)
}
