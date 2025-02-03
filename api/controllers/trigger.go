package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rishimalgwa/event-trigger-platform/api/db"
	"github.com/rishimalgwa/event-trigger-platform/api/schemas"
	"github.com/rishimalgwa/event-trigger-platform/api/views"
)

type TriggerController struct{}

func (t *TriggerController) GetAllTriggers(ctx *fiber.Ctx) error {

	triggers, err := db.TriggerSvc.FindAll()
	if err != nil {
		return views.InternalServerError(ctx, err)
	}

	return views.OK(ctx, triggers)
}

func (t *TriggerController) GetTriggerById(ctx *fiber.Ctx) error {
	return views.OK(ctx, nil)
}

func (t *TriggerController) CreateTrigger(ctx *fiber.Ctx) error {
	// validate
	payload := new(schemas.CreateTriggerRequest)
	if err := ctx.BodyParser(payload); err != nil {
		return views.InvalidJson(ctx, err)
	}

	error := payload.Validate()
	if error != nil {
		return views.ValidationError(ctx, error)
	}

	// store in db
	triggerId, err := db.TriggerSvc.Save(payload)
	if err != nil {
		return views.InternalServerError(ctx, err)
	}

	// sent to queue

	// // Handle different trigger types
	// switch models.TriggerType(trigger.Type) {
	// case models.ScheduledTrigger:
	// 	// Send to Kafka for future execution
	// err = db.TriggerSvc.ScheduleTrigger(trigger)
	// if err != nil {
	// 	return views.InternalServerError(ctx, err)
	// }

	// case models.APITrigger:
	// 	// Execute immediately
	// 	err := db.TriggerSvc.ExecuteTrigger(trigger)
	// 	if err != nil {
	// 		return views.InternalServerError(ctx, err)
	// 	}
	// }

	return views.OK(ctx, schemas.CreateTriggerResponse{
		ID:                  *triggerId,
		Type:                string(payload.Type),
		APIURL:              payload.APIURL,
		APIPayload:          payload.APIPayload,
		Schedule:            payload.ScheduleTime,
		NumberOfOccurrences: payload.NumberOfOccurrences,
	})
}

func (t *TriggerController) DeleteTrigger(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	triggerId, err := uuid.Parse(id)
	if err != nil {
		return views.BadRequest(ctx, err)
	}

	err = db.TriggerSvc.DeleteTrigger(triggerId)
	if err != nil {
		return views.InternalServerError(ctx, err)
	}
	return views.OK(ctx, nil)
}

func (t *TriggerController) UpdateTrigger(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	triggerId, err := uuid.Parse(id)
	if err != nil {
		return views.BadRequest(ctx, err)
	}

	payload := new(schemas.UpdateTriggerRequest)
	if err := ctx.BodyParser(payload); err != nil {
		return views.InvalidJson(ctx, err)
	}

	error := payload.Validate()
	if error != nil {
		return views.ValidationError(ctx, error)
	}

	err = db.TriggerSvc.UpdateTrigger(triggerId, payload)
	if err != nil {
		return views.InternalServerError(ctx, err)
	}
	return views.OK(ctx, nil)
}
