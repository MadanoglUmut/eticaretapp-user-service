package main

import (
	"UserService/internal/handlers"
	"UserService/internal/repositories"
	"UserService/internal/services"
	"UserService/pkg/psql"
	"UserService/pkg/redis"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-fiber/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	fmt.Println("Hello Rest")

	app := fiber.New()

	err := godotenv.Load("../../.env")
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

	redisClient := redis.NewClient(redisHost, redisPort, redisPassword)

	ctx := context.Background()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis'e bağlanılamadı: %v", err)
	}
	log.Println("Redis'e başarıyla bağlanıldı")

	redisC := redis.NewRedis(redisClient, ctx)

	var db = psql.Connect(host, user, password, name, port)

	userRepo := repositories.NewUserRepository(db)

	userService := services.NewUserService(userRepo)

	redisService := services.NewRedisService(redisC)

	userHand := handlers.NewUserHandler(userService, redisService)

	userHand.UserSetRoutes(app)

	sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"})

	sw.AddEndpoints(handlers.UserGetEndpoints())

	swagger.SwaggerHandler(app, sw.MustToJson(), swagger.WithPrefix("/swagger"))

	log.Fatal(app.Listen(":6060"))

}
