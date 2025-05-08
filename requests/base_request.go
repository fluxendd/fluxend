package requests

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"strconv"
	"strings"
)

var (
	reservedTableNames = map[string]bool{
		"pg_catalog":         true,
		"information_schema": true,
	}

	reservedColumnNames = map[string]bool{
		"oid":      true,
		"xmin":     true,
		"cmin":     true,
		"xmax":     true,
		"cmax":     true,
		"tableoid": true,
	}

	allowedColumnTypes = map[string]bool{
		"int":       true,
		"serial":    true,
		"varchar":   true,
		"text":      true,
		"boolean":   true,
		"date":      true,
		"timestamp": true,
		"float":     true,
		"uuid":      true,
	}

	reservedIndexNames = map[string]bool{
		"primary": true,
		"unique":  true,
		"foreign": true,
		"exclude": true,
	}
)

type BaseRequest struct {
	ProjectUUID uuid.UUID `json:"projectUUID,omitempty"`
	Context     echo.Context
}

type PaginationParams struct {
	Page  int
	Limit int
	Sort  string
	Order string
}

func (r *BaseRequest) SetContext(c echo.Context) {
	r.Context = c
}

func (r *BaseRequest) WithProjectHeader(c echo.Context) error {
	projectUUID, err := uuid.Parse(c.Request().Header.Get("X-Project"))
	if err != nil {
		return errors.New("invalid project UUID")
	}

	r.ProjectUUID = projectUUID

	return nil
}

func (r *BaseRequest) ExtractValidationErrors(err error) []string {
	if err == nil {
		return nil
	}

	var errs []string
	if ve, ok := err.(validation.Errors); ok {
		for _, validationErr := range ve {
			errs = append(errs, validationErr.Error())
		}
	}

	return errs
}

func (r *BaseRequest) ExtractPaginationParams(c echo.Context) PaginationParams {
	defaultPage := 1
	defaultLimit := 10
	defaultSort := "id"
	defaultOrder := "asc"

	// Extract and parse query parameters
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = defaultPage
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = defaultLimit
	}

	sort := c.QueryParam("sort")
	if sort == "" {
		sort = defaultSort
	}

	order := c.QueryParam("order")
	if order != "asc" && order != "desc" {
		order = defaultOrder
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Order: order,
	}
}

func (r *BaseRequest) GetUUIDPathParam(c echo.Context, name string, required bool) (uuid.UUID, error) {
	if required && c.Param(name) == "" {
		return uuid.UUID{}, fmt.Errorf("path parameter [%s] is required", name)
	}

	return uuid.Parse(c.Param(name))
}

func (r *BaseRequest) GetUUIDQueryParam(c echo.Context, name string, required bool) (uuid.UUID, error) {
	if required && c.QueryParam(name) == "" {
		return uuid.UUID{}, fmt.Errorf("query parameter [%s] is required", name)
	}

	return uuid.Parse(strings.TrimSpace(c.QueryParam(name)))
}

func IsReservedTableName(name string) bool {
	if _, ok := reservedTableNames[name]; ok {
		return true
	}

	return false
}

func IsReservedColumnName(name string) bool {
	if _, ok := reservedColumnNames[name]; ok {
		return true
	}

	return false
}

func IsAllowedColumnType(columnType string) bool {
	if _, ok := allowedColumnTypes[columnType]; ok {
		return true
	}

	return false
}
func IsReservedIndexName(name string) bool {
	if _, ok := reservedIndexNames[name]; ok {
		return true
	}

	return false
}
