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

	authv1 "github.com/emreEngineering/kervan/gen/go/auth/v1"
	grpchandler "github.com/emreEngineering/kervan/services/auth-service/internal/adapters/grpc"
	"github.com/emreEngineering/kervan/services/auth-service/internal/adapters/postgres"
	"github.com/emreEngineering/kervan/services/auth-service/internal/application"
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

	repo := postgres.NewPostgresUserRepo(db)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "test-secret"
	}
	authApp := application.NewAuthService(repo, jwtSecret)
	authHandler := grpchandler.NewAuthHandler(authApp)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Port dinlenemedi: %v", err)
	}

	s := grpc.NewServer()
	authv1.RegisterAuthServiceServer(s, authHandler)
	reflection.Register(s)

	fmt.Println("gRPC sunucusu :50051 portunda başlatıldı")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Sunucu hatası: %v", err)
	}
}
