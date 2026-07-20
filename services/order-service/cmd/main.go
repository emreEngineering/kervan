package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	orderv1 "github.com/emreEngineering/kervan/gen/go/order/v1"
	grpcclient "github.com/emreEngineering/kervan/services/order-service/internal/adapters/grpc"
	"github.com/emreEngineering/kervan/services/order-service/internal/adapters/postgres"
	"github.com/emreEngineering/kervan/services/order-service/internal/application"
)

func main() {
	// PostgreSQL bağlantısı
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://kervan:kervan_password@localhost:5432/kervan?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Veritabanı bağlantısı açılamadı: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Veritabanı ping başarısız: %v", err)
	}
	fmt.Println("Veritabanı bağlantısı başarılı")

	// gRPC client bağlantıları
	authConn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	cartConn, _ := grpc.Dial("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	inventoryConn, _ := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	productConn, _ := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))

	authClient := grpcclient.NewAuthClient(authConn)
	cartClient := grpcclient.NewCartClient(cartConn)
	inventoryClient := grpcclient.NewInventoryClient(inventoryConn)
	productClient := grpcclient.NewProductClient(productConn)

	orderRepo := postgres.NewOrderRepo(db)
	orderApp := application.NewOrderService(authClient, cartClient, inventoryClient, productClient, orderRepo)
	grpcHandler := grpcclient.NewOrderHandler(orderApp)

	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Port dinlenemedi: %v", err)
	}

	s := grpc.NewServer()
	orderv1.RegisterOrderServiceServer(s, grpcHandler)
	reflection.Register(s)

	fmt.Println("gRPC sunucusu :50055 portunda başlatıldı")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Sunucu hatası: %v", err)
	}
}
