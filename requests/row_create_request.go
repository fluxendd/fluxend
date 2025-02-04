package requests

import (
	"encoding/json"
	"fluxton/utils"
	"github.com/labstack/echo/v4"
	"io"
)

type RowCreateRequest struct {
	OrganizationID uint                   `json:"-"`
	Fields         map[string]interface{} `json:"-"`
}

func (r *RowCreateRequest) BindAndValidate(c echo.Context) []string {
	// we read request manually because we want to store it in the request nested struct
	// and default echo.Context.Bind() doesn't allow that
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return []string{"Failed to read request body"}
	}

	if err := json.Unmarshal(body, &r.Fields); err != nil {
		return []string{"Invalid JSON format"}
	}

	organizationID, err := utils.ConvertStringToUint(c.Request().Header.Get("X-OrganizationID"))
	if err != nil {
		return []string{"Organization ID is required and must be a number"}
	}

	r.OrganizationID = organizationID

	return nil
}
