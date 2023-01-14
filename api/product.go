package api

import (
	"database/sql"
	"fmt"
	db "github.com/dbssensei/ordentmarketplace/db/sqlc"
	"github.com/dbssensei/ordentmarketplace/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"strconv"
)

type createProductRequest struct {
	Name       string `json:"name" binding:"required"`
	Price      int32  `json:"price" binding:"required"`
	Status     string `json:"status" binding:"required"`
	CategoryID int32  `json:"category_id" binding:"required"`
}

type deleteProductResponse struct {
	Status string `json:"status"`
}

func (server *Server) createProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	existMerchant, _ := server.store.GetMerchantByAdminID(ctx, authPayload.UserId)
	if existMerchant.ID == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("you are not merchant admin")))
		return
	}
	arg := db.CreateProductParams{
		Name:       req.Name,
		MerchantID: existMerchant.ID,
		Price:      req.Price,
		Status:     req.Status,
		CategoryID: req.CategoryID,
	}

	product, err := server.store.CreateProduct(ctx, arg)
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

	rsp := db.Product{
		ID:         product.ID,
		Name:       product.Name,
		MerchantID: product.MerchantID,
		Price:      product.Price,
		Status:     product.Status,
		CategoryID: product.CategoryID,
		CreatedAt:  product.CreatedAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}

type updateProductRequest struct {
	Name       string `json:"name"`
	Price      int32  `json:"price"`
	Status     string `json:"status"`
	CategoryID int32  `json:"category_id"`
}

func (server *Server) updateProduct(ctx *gin.Context) {
	var req updateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	productId, _ := strconv.Atoi(ctx.Param("id"))
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	existMerchant, _ := server.store.GetMerchantByAdminID(ctx, authPayload.UserId)
	if existMerchant.ID == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("you are not merchant admin")))
		return
	}

	updatedProduct, err := server.store.UpdateProduct(ctx, db.UpdateProductParams{
		ID:         int32(productId),
		MerchantID: existMerchant.AdminID,
		Name: sql.NullString{
			String: req.Name,
			Valid:  len(req.Name) > 0,
		},
		Price: sql.NullInt32{
			Int32: req.Price,
			Valid: req.Price > 0,
		},
		Status: sql.NullString{
			String: req.Status,
			Valid:  len(req.Status) > 0,
		},
		CategoryID: sql.NullInt32{
			Int32: req.CategoryID,
			Valid: req.CategoryID > 0,
		},
	})
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

	rsp := updatedProduct
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getProducts(ctx *gin.Context) {
	products, err := server.store.GetProducts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := products
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getProduct(ctx *gin.Context) {
	productId, _ := strconv.Atoi(ctx.Param("id"))
	product, err := server.store.GetProduct(ctx, int32(productId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := product
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) deleteProduct(ctx *gin.Context) {
	productId, _ := strconv.Atoi(ctx.Param("id"))
	err := server.store.DeleteProduct(ctx, int32(productId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := deleteProductResponse{
		Status: "success",
	}
	ctx.JSON(http.StatusOK, rsp)
}
