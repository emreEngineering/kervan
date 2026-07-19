package grpc

import (
	"context"

	productv1 "github.com/emreEngineering/kervan/gen/go/product/v1"
	"google.golang.org/grpc"
)

type PrductClient struct {
	client productv1.ProductServiceClient
}

func NewProductClient(conn *grpc.ClientConn) *PrductClient {
	return &PrductClient{client: productv1.NewProductServiceClient(conn)}
}

func (c *PrductClient) GetProductPrice(ctx context.Context, productID int64) (float64, error) {
	resp, err := c.client.GetProduct(ctx, &productv1.GetProductRequest{Id: productID})
	if err != nil {
		return 0, err
	}
	return resp.Price, nil
}
