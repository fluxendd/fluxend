package controllers

import (
	"fluxton/requests"
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

// List retrieves all buckets
//
// @Summary List all buckets
// @Description Retrieve a list of buckets in a specified project.
// @Tags Buckets
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Number of items per page"
// @Param sort query string false "Field to sort by"
// @Param order query string false "Sort order (asc or desc)"
//
// @Success 200 {object} responses.Response{content=[]resources.BucketResponse} "List of buckets"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /storage [get]
func (bc *BucketController) List(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	paginationParams := utils.ExtractPaginationParams(c)
	buckets, err := bc.bucketService.List(paginationParams, request.ProjectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.BucketResourceCollection(buckets))
}

// Show retrieves details of a specific bucket.
//
// @Summary Show details of a single bucket
// @Description Get details of a specific bucket
// @Tags Buckets
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param bucketUUID path string true "Bucket UUID"
//
// @Success 200 {object} responses.Response{content=resources.BucketResponse} "Bucket details"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /storage/{bucketUUID} [get]
func (bc *BucketController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	bucketUUID, err := utils.GetUUIDPathParam(c, "bucketUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	bucket, err := bc.bucketService.GetByUUID(bucketUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.BucketResource(&bucket))
}

// Store creates a new bucket
//
// @Summary Create a new bucket
// @Description Add a new bucket to a project
// @Tags Buckets
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
// @Param bucket body bucket_requests.CreateRequest true "Bucket details"
//
// @Success 201 {object} responses.Response{content=resources.BucketResponse} "Bucket created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /storage [post]
func (bc *BucketController) Store(c echo.Context) error {
	var request bucket_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	bucket, err := bc.bucketService.Create(&request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.BucketResource(&bucket))
}

// Update a bucket
//
// @Summary Update a bucket
// @Description Modify an existing bucket's details
// @Tags Buckets
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param bucketUUID path string true "Bucket UUID"
// @Param bucket body bucket_requests.CreateRequest true "Bucket details"
//
// @Success 200 {object} responses.Response{content=resources.BucketResponse} "Bucket updated"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /storage/{bucketUUID} [put]
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

	updatedBucket, err := bc.bucketService.Update(bucketUUID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.BucketResource(updatedBucket))
}

// Delete a bucket
//
// @Summary Delete a bucket
// @Description Remove a bucket from a project
// @Tags Buckets
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param bucketUUID path string true "Bucket UUID"
//
// @Success 204 "Bucket deleted"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /storage/{bucketUUID} [delete]
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
