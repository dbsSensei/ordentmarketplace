// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CountCountry(ctx context.Context) (int64, error)
	CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error)
	CreateCountry(ctx context.Context, arg CreateCountryParams) (Country, error)
	CreateMerchant(ctx context.Context, arg CreateMerchantParams) (Merchant, error)
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (OrderItem, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteCategory(ctx context.Context, id int32) error
	DeleteProduct(ctx context.Context, id int32) error
	GetCategories(ctx context.Context) ([]Category, error)
	GetCategory(ctx context.Context, id int32) (Category, error)
	GetCountry(ctx context.Context, code string) (Country, error)
	GetMerchantByAdminID(ctx context.Context, adminID int32) (Merchant, error)
	GetOrder(ctx context.Context, arg GetOrderParams) (Order, error)
	GetOrderItems(ctx context.Context, orderID int32) ([]OrderItem, error)
	GetOrders(ctx context.Context, userID int32) ([]Order, error)
	GetProduct(ctx context.Context, id int32) (Product, error)
	GetProducts(ctx context.Context) ([]Product, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, id int32) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	UpdateMerchant(ctx context.Context, arg UpdateMerchantParams) (Merchant, error)
	UpdateOrderStatus(ctx context.Context, arg UpdateOrderStatusParams) (Order, error)
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)