package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	cartv1 "github.com/emreEngineering/kervan/gen/go/cart/v1"
	grpchandler "github.com/emreEngineering/kervan/services/cart-service/internal/adapters/grpc"
	redisrepo "github.com/emreEngineering/kervan/services/cart-service/internal/adapters/redis"
	"github.com/emreEngineering/kervan/services/cart-service/internal/application"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer client.Close()

	repo := redisrepo.NewCartRepo(client)
	cartApp := application.NewCartService(repo)
	cartHandler := grpchandler.NewCartHandler(cartApp)

	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("port dinlenemedi: %v", err)
	}

	s := grpc.NewServer()
	cartv1.RegisterCartServiceServer(s, cartHandler)
	reflection.Register(s)

	fmt.Println("gRPC sunucus :50054 portunda başlatıldı")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("sunucu hatası: %v", err)
	}
}
