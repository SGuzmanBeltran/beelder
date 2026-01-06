package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateBody is a generic middleware for validating request bodies
func ValidateBody[T any](c *fiber.Ctx) error {
	body := new(T)

	// Parse request body
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate body
	if errors := ValidateStruct(body); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	// Store validated body in context for the next handler
	c.Locals("validated", body)

	return c.Next()
}

// ValidateQuery is a generic middleware for validating query parameters
func ValidateQuery[T any](c *fiber.Ctx) error {
	params := new(T)

	// Parse query parameters
	if err := c.QueryParser(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	// Validate params
	if errors := ValidateStruct(params); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	// Store validated params in context for the next handler
	c.Locals("validated", params)

	return c.Next()
}

// ValidateStruct is a generic validation function that can validate any struct
func ValidateStruct[T any](data *T) []ValidationError {
	var errors []ValidationError

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Field:   err.Field(),
				Message: getErrorMsg(err),
			})
		}
	}

	return errors
}

func toJSONName(fieldName string) string {
	if fieldName == "" {
		return ""
	}
	return strings.ToLower(fieldName[:1]) + fieldName[1:]
}

// getErrorMsg converts validator errors into user-friendly messages
func getErrorMsg(fe validator.FieldError) string {
	fieldName := toJSONName(fe.Field())
	feParam := toJSONName(fe.Param())

	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "oneof":
		return "Must be one of: " + feParam
	case "min":
		return "Should be at least " + feParam
	case "max":
		return "Should be at most " + feParam
	case "uuid":
		return "Must be a valid UUID"
	case "datetime":
		return "Must be a valid datetime"
	case "number":
		return "Must be a valid number"
	case "url":
		return "Must be a valid URL"
	case "gtfield":
		return "Must be greater than " + feParam
	case "ltfield":
		return "Must be less than " + feParam
	case "gtefield":
		return "Must be greater than or equal to " + feParam
	case "ltefield":
		return "Must be less than or equal to " + feParam
	case "required_with":
		return "Required when " + feParam + " is present"
	case "required_without":
		return "Required when " + feParam + " is not present"
	case "json":
		return "Must be a valid JSON"
	case "gt":
		return "Must be greater than " + feParam
	case "lt":
		return "Must be less than " + feParam
	case "gte":
		return "Must be greater than or equal to " + feParam
	case "lte":
		return "Must be less than or equal to " + feParam
	default:
		return "Invalid value for field " + fieldName
	}
}
