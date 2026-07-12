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

	productv1 "github.com/emreEngineering/kervan/gen/go/product/v1"
	grpchandler "github.com/emreEngineering/kervan/services/product-service/internal/adapters/grpc"
	"github.com/emreEngineering/kervan/services/product-service/internal/adapters/postgres"
	"github.com/emreEngineering/kervan/services/product-service/internal/application"
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

	repo := postgres.NewPostgresProductRepo(db)
	productApp := application.NewProductService(repo)
	productHandler := grpchandler.NewProductHandler(productApp)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Port dinlenemedi: %v", err)
	}

	s := grpc.NewServer()
	productv1.RegisterProductServiceServer(s, productHandler)
	reflection.Register(s)

	fmt.Println("gRPC sunucusu :50052 portunda başlatıldı")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Sunucu hatası: %v", err)
	}
}
