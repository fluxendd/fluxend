package utils

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"strconv"
)

type PaginationParams struct {
	Page  int
	Limit int
	Sort  string
	Order string
}

func ExtractPaginationParams(c echo.Context) PaginationParams {
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

func GetUintQueryParam(c echo.Context, name string, required bool) (uint, error) {
	if required && c.QueryParam(name) == "" {
		return 0, fmt.Errorf("query parameter [%s] is required", name)
	}

	return ConvertStringToUint(c.QueryParam(name))
}

func GetUintPathParam(c echo.Context, name string, required bool) (uint, error) {
	if required && c.Param(name) == "" {
		return 0, fmt.Errorf("path parameter [%s] is required", name)
	}

	return ConvertStringToUint(c.Param(name))
}

func GetUUIDPathParam(c echo.Context, name string, required bool) (uuid.UUID, error) {
	if required && c.Param(name) == "" {
		return uuid.UUID{}, fmt.Errorf("path parameter [%s] is required", name)
	}

	return uuid.Parse(c.Param(name))
}

func ConvertStringToUint(param string) (uint, error) {
	value, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return 0, errors.New("provided value is not a valid integer")
	}

	return uint(value), nil
}
