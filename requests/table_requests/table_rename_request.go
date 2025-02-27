package table_requests

import (
	"fluxton/configs"
	"fluxton/requests"
	"fluxton/utils"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"regexp"
)

type TableRenameRequest struct {
	requests.BaseRequest
	Name string `json:"name"`
}

func (r *TableRenameRequest) BindAndValidate(c echo.Context) []string {
	if err := c.Bind(r); err != nil {
		return []string{"Invalid request payload"}
	}

	err := validation.ValidateStruct(r,
		validation.Field(
			&r.Name, validation.Required.Error("Name is required"),
			validation.Match(
				regexp.MustCompile(utils.AlphanumericWithUnderscorePattern()),
			).Error("Table name must be alphanumeric with underscores"),
			validation.Length(
				configs.MinTableNameLength, configs.MaxTableNameLength,
			).Error(
				fmt.Sprintf(
					"Name must be between %d and %d characters",
					configs.MinTableNameLength,
					configs.MaxTableNameLength,
				),
			),
			validation.By(validateName),
		),
	)

	return r.ExtractValidationErrors(err)
}
