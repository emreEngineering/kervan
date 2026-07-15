package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/emreEngineering/kervan/services/cart-service/internal/domain"
)

type CartRepo struct {
	client *redis.Client
}

func NewCartRepo(client *redis.Client) *CartRepo {
	return &CartRepo{client: client}
}

func (r *CartRepo) cartKey(userID int64) string {
	return fmt.Sprintf("cart:%d", userID)
}

func (r *CartRepo) GetCart(ctx context.Context, userID int64) (*domain.Cart, error) {
	data, err := r.client.HGetAll(ctx, r.cartKey(userID)).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	cart := domain.NewCart(userID)
	for productIDStr, quantityStr := range data {
		var item domain.CartItem
		item.ProductID = parseInt(productIDStr)
		item.Quantity = int32(parseInt(quantityStr))
		cart.Items = append(cart.Items, item)
	}
	return cart, nil
}

func (r *CartRepo) SaveCart(ctx context.Context, cart *domain.Cart) error {
	key := r.cartKey(cart.UserID)
	for _, item := range cart.Items {
		productIDStr := fmt.Sprintf("%d", item.ProductID)
		quantityStr := fmt.Sprintf("%d", item.Quantity)
		if err := r.client.HSet(ctx, key, productIDStr, quantityStr).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r *CartRepo) ClearCart(ctx context.Context, userID int64) error {
	return r.client.Del(ctx, r.cartKey(userID)).Err()
}

func parseInt(s string) int64 {
	var n int64
	fmt.Sscanf(s, "%d", &n)
	return n
}
