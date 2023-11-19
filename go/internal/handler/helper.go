package handler

import (
	"context"
	"encoding/json"
	"fancy-todo/internal/constant"
	"fancy-todo/internal/libs"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func CreateContext(c echo.Context) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, libs.RequestId, c.Response().Header().Get(echo.HeaderXRequestID))
	ctx = context.WithValue(ctx, libs.IpAddress, c.RealIP())
	ctx = context.WithValue(ctx, libs.UserID, c.Get("user_id"))
	ctx = context.WithValue(ctx, libs.UserID, c.Get("user_email"))
	return ctx
}

func PreprocessedRequest(c echo.Context, validate *validator.Validate, body any) (context.Context, error) {
	ctx := CreateContext(c)
	
	if body != nil {
		err := json.NewDecoder(c.Request().Body).Decode(&body)
		if err != nil {
			return ctx, libs.DefaultInternalServerError(err)
		}
		
		err = validate.Struct(body)
		if err != nil {
			return ctx, libs.DefaultInternalServerError(err)
		}
	}

	return ctx, nil
}

func GetUserIdFromContext(c echo.Context) (int64, error) {
	userId := c.Get("user_id")

	switch v := userId.(type) {
	case int64:
		return v, nil
	default:
		return 0, libs.DefaultInternalServerError(fmt.Errorf("invalid userID %v", userId))
	}
}

func GetIdFromPathParam(c echo.Context, key string) (int64, error) {
	idParam := c.Param(key)
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, libs.CustomError{
			HTTPCode: http.StatusBadRequest,
			Message: fmt.Sprintf("invalid %s ID",key),
		}
	}

	return id, nil
}

func ConvertCommonQueryParam(c echo.Context) (libs.QueryParam, error) {
	pageSizeStr := c.QueryParam("page_size")
	pageNumberStr := c.QueryParam("page_number")
	sortKey := c.QueryParam("sort_key")
	sortOrder := c.QueryParam("sort_order")

	var queryParam libs.QueryParam
	var pageSize, pageNumber int
	var err error

	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			return queryParam, libs.CustomError{
				HTTPCode: http.StatusBadRequest,
				Message: "invalid page size query",
			}
		}

		if pageSize < constant.PAGE_SIZE_MIN || pageSize > constant.PAGE_SIZE_MAX {
			pageSize = constant.PAGE_SIZE_DEFAULT
		}
	} else {
		pageSize = constant.PAGE_SIZE_DEFAULT
	} 

	if pageNumberStr != "" {
		pageNumber, err = strconv.Atoi(pageNumberStr)
		if err != nil {
			return queryParam, libs.CustomError{
				HTTPCode: http.StatusBadRequest,
				Message: "invalid page number query",
			}
		}

		if pageNumber < constant.PAGE_NUMBER_MIN || pageNumber > constant.PAGE_NUMBER_MAX {
			pageNumber = constant.PAGE_NUMBER_DEFAULT
		}
	} else {
		pageNumber = constant.PAGE_NUMBER_DEFAULT
	}

	queryParam.PageNumber = pageNumber
	queryParam.PageSize = pageSize
	queryParam.SortKey = sortKey
	queryParam.SortOrder = sortOrder

	return queryParam, nil
}