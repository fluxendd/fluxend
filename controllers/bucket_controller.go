package controllers

import (
	"fluxton/requests/bucket_requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type BucketController struct {
	bucketService services.BucketService
}

func NewBucketController(injector *do.Injector) (*BucketController, error) {
	bucketService := do.MustInvoke[services.BucketService](injector)

	return &BucketController{bucketService: bucketService}, nil
}

func (bc *BucketController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, "Invalid project UUID")
	}

	paginationParams := utils.ExtractPaginationParams(c)
	buckets, err := bc.bucketService.List(paginationParams, projectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.BucketResourceCollection(buckets))
}

func (bc *BucketController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	bucketUUID, err := utils.GetUUIDPathParam(c, "bucketUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	bucket, err := bc.bucketService.GetByUUID(bucketUUID, projectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.BucketResource(&bucket))
}

func (bc *BucketController) Store(c echo.Context) error {
	var request bucket_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	bucket, err := bc.bucketService.Create(projectUUID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.BucketResource(&bucket))
}

func (bc *BucketController) Update(c echo.Context) error {
	var request bucket_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	bucketUUID, err := utils.GetUUIDPathParam(c, "bucketUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	updatedForm, err := bc.bucketService.Update(bucketUUID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.BucketResource(updatedForm))
}

func (bc *BucketController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	bucketUUID, err := utils.GetUUIDPathParam(c, "bucketUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := bc.bucketService.Delete(bucketUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
