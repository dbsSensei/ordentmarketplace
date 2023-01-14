package api

import (
	db "github.com/dbssensei/ordentmarketplace/db/sqlc"
	"github.com/dbssensei/ordentmarketplace/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"strconv"
	"time"
)

func (server *Server) createOrder(ctx *gin.Context) {
	var req db.OrderTxParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.OrderTxParams{
		UserID: authPayload.UserId,
		Status: "order_received",
		Items:  req.Items,
	}

	order, err := server.store.OrderTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			//case "unique_violation":
			//	ctx.JSON(http.StatusForbidden, errorResponse(err))
			//	return
			//}
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		rsp := db.OrderTxResult{
			ID:        order.ID,
			Status:    order.Status,
			Items:     order.Items,
			CreatedAt: order.CreatedAt,
		}
		ctx.JSON(http.StatusOK, rsp)
	}
}

func (server *Server) getOrders(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	order, err := server.store.GetOrders(ctx, authPayload.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := order
	ctx.JSON(http.StatusOK, rsp)
}

type detailOrderResponse struct {
	ID        int32              `json:"id"`
	Status    string             `json:"status"`
	Items     []detailOrderItems `json:"items"`
	CreatedAt time.Time          `json:"created_at"`
}

type detailOrderItems struct {
	ProductID  int32  `json:"product_id"`
	Name       string `json:"name"`
	MerchantID int32  `json:"merchant_id"`
	Price      int32  `json:"price"`
	CategoryID int32  `json:"category_id"`
	Quantity   int32  `json:"quantity"`
}

func (server *Server) getOrder(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	orderId, _ := strconv.Atoi(ctx.Param("id"))
	order, err := server.store.GetOrder(ctx, db.GetOrderParams{
		ID:     int32(orderId),
		UserID: authPayload.UserId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	orderItems, err := server.store.GetOrderItems(ctx, order.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var items []detailOrderItems
	for _, orderItem := range orderItems {
		product, err := server.store.GetProduct(ctx, orderItem.ProductID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		items = append(items, detailOrderItems{
			ProductID:  product.ID,
			Name:       product.Name,
			MerchantID: product.MerchantID,
			Price:      product.Price,
			CategoryID: product.CategoryID,
			Quantity:   orderItem.Quantity,
		})
	}

	rsp := detailOrderResponse{
		ID:        order.ID,
		Status:    order.Status,
		Items:     items,
		CreatedAt: order.CreatedAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}

type updateOrderStatusRequest struct {
	Status string `json:"status"`
}

func (server *Server) updateOrderStatus(ctx *gin.Context) {
	var req updateOrderStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	productId, _ := strconv.Atoi(ctx.Param("id"))
	updatedOrderStatus, err := server.store.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		ID:     int32(productId),
		Status: req.Status,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			//case "unique_violation":
			//	ctx.JSON(http.StatusForbidden, errorResponse(err))
			//	return
			//}
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		rsp := updatedOrderStatus
		ctx.JSON(http.StatusOK, rsp)
	}
}
