package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rishimalgwa/event-trigger-platform/api/db"
	"github.com/rishimalgwa/event-trigger-platform/api/schemas"
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

	trigger, err := db.TriggerSvc.Find(triggerId)
	if err != nil {
		return views.InternalServerError(ctx, err)
	}

	// sent to queue

	// // Handle different trigger types
	// switch models.TriggerType(trigger.Type) {
	// case models.ScheduledTrigger:
	// 	// Send to Kafka for future execution
	err = db.TriggerSvc.ScheduleTrigger(trigger)
	if err != nil {
		return views.InternalServerError(ctx, err)
	}

	// case models.APITrigger:
	// 	// Execute immediately
	// 	err := db.TriggerSvc.ExecuteTrigger(trigger)
	// 	if err != nil {
	// 		return views.InternalServerError(ctx, err)
	// 	}
	// }

	return views.OK(ctx, schemas.CreateTriggerResponse{
		ID:                  trigger.ID,
		Type:                string(trigger.Type),
		APIURL:              trigger.APIURL,
		APIPayload:          trigger.APIPayload,
		Schedule:            trigger.ScheduleTime,
		NumberOfOccurrences: trigger.NumberOfOccurrences,
		CreatedAt:           trigger.CreatedAt,
	})
}

func (t *TriggerController) DeleteTrigger(ctx *fiber.Ctx) error {
	return views.OK(ctx, nil)
}

func (t *TriggerController) UpdateTrigger(ctx *fiber.Ctx) error {
	return views.OK(ctx, nil)
}
