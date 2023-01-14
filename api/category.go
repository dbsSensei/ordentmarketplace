package api

import (
	"database/sql"
	db "github.com/dbssensei/ordentmarketplace/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"strconv"
)

type createCategoryRequest struct {
	Name     string        `json:"name" binding:"required"`
	ParentID sql.NullInt32 `json:"parent_id"`
}

type categoryResponse struct {
	ID       int32         `json:"id"`
	Name     string        `json:"name"`
	ParentID sql.NullInt32 `json:"parent_id"`
}

type deleteCategoryResponse struct {
	Status string `json:"status"`
}

func (server *Server) createCategory(ctx *gin.Context) {
	var req createCategoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateCategoryParams{
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	category, err := server.store.CreateCategory(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := categoryResponse{
		ID:       category.ID,
		Name:     category.Name,
		ParentID: category.ParentID,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getCategories(ctx *gin.Context) {
	categories, err := server.store.GetCategories(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := categories
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getCategory(ctx *gin.Context) {
	categoryId, _ := strconv.Atoi(ctx.Param("id"))
	category, err := server.store.GetCategory(ctx, int32(categoryId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := category
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) deleteCategory(ctx *gin.Context) {
	categoryId, _ := strconv.Atoi(ctx.Param("id"))
	err := server.store.DeleteCategory(ctx, int32(categoryId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := deleteCategoryResponse{
		Status: "success",
	}
	ctx.JSON(http.StatusOK, rsp)
}
