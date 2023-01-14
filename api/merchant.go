package api

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dbssensei/ordentmarketplace/token"
	"github.com/lib/pq"
	"net/http"
	"time"

	db "github.com/dbssensei/ordentmarketplace/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createMerchantRequest struct {
	Name        string `json:"name" binding:"required"`
	CountryCode string `json:"country_code" binding:"required"`
}

type merchantResponse struct {
	AdminID     int32     `json:"admin_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	CountryCode string    `json:"country_code" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
}

func newMerchantResponse(merchant db.Merchant) merchantResponse {
	return merchantResponse{
		AdminID:     merchant.AdminID,
		CountryCode: merchant.CountryCode,
		Name:        merchant.Name,
		CreatedAt:   merchant.CreatedAt,
	}
}

func (server *Server) createMerchant(ctx *gin.Context) {
	var req createMerchantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateMerchantParams{
		AdminID:     authPayload.UserId,
		Name:        req.Name,
		CountryCode: req.CountryCode,
	}

	existMerchant, _ := server.store.GetMerchantByAdminID(ctx, arg.AdminID)
	if existMerchant.ID != 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("one user can only have one merchant")))
		return
	}

	merchant, err := server.store.CreateMerchant(ctx, arg)
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

	rsp := newMerchantResponse(merchant)
	ctx.JSON(http.StatusOK, rsp)
}

type updateMerchantRequest struct {
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
}

func (server *Server) updateMerchant(ctx *gin.Context) {
	var req updateMerchantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	updatedMerchant, err := server.store.UpdateMerchant(context.Background(), db.UpdateMerchantParams{
		AdminID: authPayload.UserId,
		Name: sql.NullString{
			String: req.Name,
			Valid:  len(req.Name) > 0,
		},
		CountryCode: sql.NullString{
			String: req.CountryCode,
			Valid:  len(req.CountryCode) > 0,
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

	rsp := newMerchantResponse(updatedMerchant)
	ctx.JSON(http.StatusOK, rsp)
}
