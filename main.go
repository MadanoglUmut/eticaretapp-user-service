package main

import (
	"UserService/internal/handlers"
	"UserService/internal/repositories"
	"UserService/internal/services"
	"UserService/pkg/metrics"
	"UserService/pkg/psql"
	"UserService/pkg/redis"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-fiber/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {

	fmt.Println("Hello Rest")

	app := fiber.New()

	//err := godotenv.Load("../../.env")

	err := godotenv.Load(".env")
	if err != nil {

		log.Fatal("Env Dosyası Yüklenemedi", err)

	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	var db = psql.Connect(host, user, password, name, port)

	redisClient := redis.NewClient(redisHost, redisPort, redisPassword)

	ctx := context.Background()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis'e bağlanılamadı: %v", err)
	}
	log.Println("Redis'e başarıyla bağlanıldı")

	redisC := redis.NewRedis(redisClient, ctx)

	userRepo := repositories.NewUserRepository(db)

	userService := services.NewUserService(userRepo)

	redisService := services.NewRedisService(redisC)

	histogram := metrics.NewNamedHistogram("http_request_userservice_duration_seconds", []float64{0.001, 0.005, 0.01, 0.05, 0.1})

	registry := prometheus.NewRegistry()
	registry.MustRegister(histogram.Histogram)

	userHand := handlers.NewUserHandler(userService, redisService, histogram)

	userHand.UserSetRoutes(app)

	app.Get("/metrics", adaptor.HTTPHandler(metrics.GetHandler(registry)))

	sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"})

	sw.AddEndpoints(handlers.UserGetEndpoints())

	swagger.SwaggerHandler(app, sw.MustToJson(), swagger.WithPrefix("/swagger"))

	log.Fatal(app.Listen(":6060"))

}
