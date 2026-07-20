package grpc

import (
	"context"

	productv1 "github.com/emreEngineering/kervan/gen/go/product/v1"
	"google.golang.org/grpc"
)

type ProductClient struct {
	client productv1.ProductServiceClient
}

func NewProductClient(conn *grpc.ClientConn) *ProductClient {
	return &ProductClient{client: productv1.NewProductServiceClient(conn)}
}

func (c *ProductClient) GetProductPrice(ctx context.Context, productID int64) (float64, error) {
	resp, err := c.client.GetProduct(ctx, &productv1.GetProductRequest{Id: productID})
	if err != nil {
		return 0, err
	}
	return resp.Price, nil
}
