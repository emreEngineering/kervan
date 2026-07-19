package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcclient "github.com/emreEngineering/kervan/services/order-service/internal/adapters/grpc"
	httphandler "github.com/emreEngineering/kervan/services/order-service/internal/adapters/http"
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
	orderHandler := httphandler.NewOrderHandler(orderApp)

	// HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/orders", orderHandler.CreateOrder)
	mux.HandleFunc("/orders/", orderHandler.GetOrder)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Port dinlenemedi: %v", err)
	}

	fmt.Println("REST sunucusu :8080 portunda başlatıldı")
	if err := http.Serve(lis, mux); err != nil {
		log.Fatalf("Sunucu hatası: %v", err)
	}
}
