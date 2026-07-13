package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryv1 "github.com/emreEngineering/kervan/gen/go/inventory/v1"
	grpchandler "github.com/emreEngineering/kervan/services/inventory-service/internal/adapters/grpc"
	"github.com/emreEngineering/kervan/services/inventory-service/internal/adapters/postgres"
	"github.com/emreEngineering/kervan/services/inventory-service/internal/application"
)

func main() {
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

	repo := postgres.NewPostgresStockRepo(db)
	inventoryApp := application.NewInventoryService(repo)
	inventoryHandler := grpchandler.NewInventoryHandler(inventoryApp)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Port dinlenemedi: %v", err)
	}

	s := grpc.NewServer()
	inventoryv1.RegisterInventoryServiceServer(s, inventoryHandler)
	reflection.Register(s)

	fmt.Println("gRPC sunucusu :50053 portunda başlatıldı")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Sunucu hatası: %v", err)
	}

}
