package main

import (
	"UserService/internal/handlers"
	"UserService/internal/repositories"
	"UserService/internal/services"
	"UserService/pkg/psql"
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

	var db = psql.Connect(host, user, password, name, port)

	userRepo := repositories.NewUserRepository(db)

	userService := services.NewUserService(userRepo)

	userHand := handlers.NewUserHandler(userService)

	userHand.UserSetRoutes(app)

	sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"})

	sw.AddEndpoints(handlers.UserGetEndpoints())

	swagger.SwaggerHandler(app, sw.MustToJson(), swagger.WithPrefix("/swagger"))

	log.Fatal(app.Listen(":6060"))

}
