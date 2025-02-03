package views

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rishimalgwa/event-trigger-platform/api/schemas"
	"github.com/rishimalgwa/event-trigger-platform/config"
)

type ReturnMsg struct {
	Msg        string      `json:"msg"`
	Err        string      `json:"err,omitempty"`
	StatusCode int         `json:"statusCode,omitempty"`
	Body       interface{} `json:"body,omitempty"`
}

func InvalidJson(c *fiber.Ctx, err error) error {
	status := fiber.StatusUnprocessableEntity
	return c.Status(status).JSON(ReturnMsg{
		Msg:        "invalid json",
		Err:        fmt.Sprintf("%v", err),
		StatusCode: status,
	})
}

func NotFound(c *fiber.Ctx, err error) error {
	status := fiber.StatusNotFound
	return c.Status(status).JSON(ReturnMsg{
		Msg:        "not found",
		Err:        fmt.Sprintf("%v", err),
		StatusCode: status,
	})
}

func AlreadyExists(c *fiber.Ctx, err error) error {
	status := fiber.StatusConflict
	return c.Status(status).JSON(ReturnMsg{
		Msg:        "already exists",
		Err:        fmt.Sprintf("%v", err),
		StatusCode: status,
	})
}

func InternalServerError(c *fiber.Ctx, err error) error {
	if config.ENVIRONMENT == "development" || config.ENVIRONMENT == "staging" {
		log.Println("error: ", err.Error())
	}
	status := fiber.StatusInternalServerError
	return c.Status(status).JSON(ReturnMsg{
		Msg:        "something went wrong",
		Err:        "something went wrong",
		StatusCode: status,
	})
}

func Created(c *fiber.Ctx, body interface{}) error {
	status := fiber.StatusCreated
	return c.Status(status).JSON(ReturnMsg{
		Msg:        "created",
		StatusCode: status,
		Body:       body,
	})
}

func OK(c *fiber.Ctx, body interface{}) error {
	status := fiber.StatusOK
	return c.Status(status).JSON(ReturnMsg{
		Msg:        "ok",
		StatusCode: status,
		Body:       body,
	})
}

func Unauthorized(c *fiber.Ctx, err error) error {
	status := fiber.StatusUnauthorized
	return c.Status(status).JSON(ReturnMsg{
		Msg:        "incorrect or missing credentials",
		Err:        fmt.Sprintf("%v", err),
		StatusCode: status,
	})
}

func Forbidden(c *fiber.Ctx, err error) error {
	status := fiber.StatusForbidden
	return c.Status(status).JSON(ReturnMsg{
		Msg:        "you are forbidden from accessing this resource",
		Err:        fmt.Sprintf("%v", err),
		StatusCode: status,
	})
}

func ValidationError(c *fiber.Ctx, err []*schemas.ErrorResponse) error {
	status := fiber.StatusBadRequest
	return c.Status(status).JSON(ReturnMsg{
		Msg:        "Invalid JSON",
		Err:        "invalid json",
		StatusCode: status,
		Body:       err,
	})
}

func InvalidQuery(c *fiber.Ctx, err error) error {
	status := fiber.StatusExpectationFailed
	return c.Status(status).JSON(ReturnMsg{
		Msg:        "invalid json",
		Err:        fmt.Sprintf("%v", err),
		StatusCode: status,
	})
}
func TokenExpired(c *fiber.Ctx, err error) error {
	status := fiber.StatusUnauthorized
	return c.Status(status).JSON(ReturnMsg{
		Msg:        "token expired",
		Err:        fmt.Sprintf("%v", err),
		StatusCode: status,
	})
}

func BadRequest(c *fiber.Ctx, err error) error {
	status := fiber.StatusBadRequest
	return c.Status(status).JSON(ReturnMsg{
		Msg:        "bad request",
		Err:        fmt.Sprintf("%v", err),
		StatusCode: status,
	})
}
