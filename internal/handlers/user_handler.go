package handlers

import (
	"UserService/internal/models"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type userService interface {
	GetUser(userId int) (models.Users, error)
	GetUserByEmail(userEmail string) (models.Users, error)
	CreateUser(createdUser models.CreateUsers) (models.Users, error)
	UpdateUser(userId int, updatedUser models.UpdateUsers) (models.Users, error)
	Delete(userId int) error
}

type redisService interface {
	Set(userId int, token string) error
	GetTokens(userId int) ([]string, error)
}

type UserHandler struct {
	userService  userService
	redisService redisService
}

func NewUserHandler(userService userService, redisService redisService) *UserHandler {

	return &UserHandler{
		userService:  userService,
		redisService: redisService,
	}

}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func (h *UserHandler) Login(c *fiber.Ctx) error {

	var loginRequest models.LoginRequest

	if err := c.BodyParser(&loginRequest); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{Error: "Body Parse Hatasi", Details: err.Error()})

	}

	user, err := h.userService.GetUserByEmail(loginRequest.Email)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.FailResponse{Error: "Kullanıcı Bulunamadi", Details: err.Error()})
	}

	if user.Password != loginRequest.Password {
		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{Error: "Şifre Hatali", Details: "Şifreniz Hatali Tekrar Deneyin"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.FailResponse{Error: "Token Oluşturulamadi", Details: err.Error()})
	}
	//REDİS SERVİS DEĞİL CLİENT OLCAK
	bearerToken := tokenString
	err = h.redisService.Set(user.ID, bearerToken)

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(models.FailResponse{
			Error:   "Token Redis'e yazılamadı",
			Details: err.Error()},
		)

	}
	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: tokenString})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: "Cikis Islemi Basarili"})

}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {

	randomNumber := rand.Intn(100)

	if randomNumber < 5 {
		return c.Status(fiber.StatusServiceUnavailable).JSON(models.FailResponse{
			Error:   "Sunucu  hizmet veremiyor",
			Details: "Tekrar deneyin",
		})
	}

	userClaims := c.Locals("user").(jwt.MapClaims)
	userId := int(userClaims["id"].(float64))

	user, err := h.userService.GetUser(userId)

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(models.FailResponse{Error: "Kullanici Bulunamadi", Details: err.Error()})

	}

	fmt.Println("Token Claims:", userClaims)

	fmt.Println("User ID:", userId)

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: user})

}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {

	var createdUser models.CreateUsers

	if err := c.BodyParser(&createdUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{Error: "Body Parse Hatasi", Details: err.Error()})
	}

	user, err := h.userService.CreateUser(createdUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.FailResponse{Error: "Kullanici Oluşturulamadi", Details: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccesResponse{SuccesData: user})

}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {

	userClaims := c.Locals("user").(jwt.MapClaims)
	userId := int(userClaims["id"].(float64))

	var updatedUser models.UpdateUsers

	if err := c.BodyParser(&updatedUser); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{Error: "Body Parse Hatasi", Details: err.Error()})

	}

	user, err := h.userService.UpdateUser(userId, updatedUser)

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(models.FailResponse{Error: "Kullanici Güncellenemedi", Details: err.Error()})

	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: user})

}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {

	userClaims := c.Locals("user").(jwt.MapClaims)
	userId := int(userClaims["id"].(float64))

	err := h.userService.Delete(userId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.FailResponse{Error: "Kullanici Silinemedi", Details: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: "Kullanıcı Silindi"})

}

func (h *UserHandler) JWTMiddleware(c *fiber.Ctx) error {
	fmt.Println("Hello")
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: models.ErrMissingAuthorization})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	fmt.Println("tokenString:", tokenString)
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: models.ErrMissingAuthorization})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("beklenmeyen imzalama metodu")
		}
		return jwtSecret, nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: models.ErrTokenIsExpired})
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: models.ErrInvalidToken})
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: models.ErrInvalidToken})
	}

	userId := int(claims["id"].(float64))

	tokens, err := h.redisService.GetTokens(userId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.FailResponse{
			Error:   "Redis'te token kontrolü sırasında hata oluştu",
			Details: err.Error(),
		})
	}

	fmt.Println(tokens[0])

	tokenFound := false
	for _, t := range tokens {
		if t == tokenString {
			tokenFound = true
			break
		}
	}

	if !tokenFound {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: models.ErrInvalidToken})
	}

	c.Locals("user", claims)
	fmt.Println("TOKEN OK")
	return c.Next()
}

func (h *UserHandler) UserSetRoutes(app *fiber.App) {

	app.Post("/login", h.Login)
	app.Post("/logout", h.Logout)
	app.Post("/users", h.CreateUser)

	userRoutesGroup := app.Group("/users")
	userRoutesGroup.Use(h.JWTMiddleware)

	userRoutesGroup.Get("/me", h.GetUser)
	userRoutesGroup.Put("/me", h.UpdateUser)
	userRoutesGroup.Delete("/me", h.DeleteUser)
}

func UserGetEndpoints() []*endpoint.EndPoint {

	return []*endpoint.EndPoint{
		endpoint.New(
			endpoint.POST,
			"/login",
			endpoint.WithTags("users"),
			endpoint.WithBody(models.LoginRequest{}),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(struct{ SuccesData string }{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{
				response.New(models.FailResponse{}, "400", "Bad Request"),
				response.New(models.FailResponse{}, "500", "Internal Server")}),
		),

		endpoint.New(
			endpoint.POST,
			"/logout",
			endpoint.WithTags("users"),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.SuccesResponse{}, "200", "OK")}),
		),

		endpoint.New(
			endpoint.POST,
			"/users",
			endpoint.WithTags("users"),
			endpoint.WithBody(models.CreateUsers{}),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.Users{}, "201", "Created")}),
			endpoint.WithErrors([]response.Response{
				response.New(models.FailResponse{}, "400", "Bad Request"),
				response.New(models.FailResponse{}, "500", "Internal Server")}),
		),
		endpoint.New(
			endpoint.GET,
			"/users/me",
			endpoint.WithTags("users"),
			endpoint.WithParams(parameter.StrParam("Authorization", parameter.Header, parameter.WithRequired())),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.Users{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.FailResponse{}, "500", "Internal Server")}),
		),

		endpoint.New(
			endpoint.PUT,
			"/users/me",
			endpoint.WithBody(models.UpdateUsers{}),
			endpoint.WithTags("users"),
			endpoint.WithParams(parameter.StrParam("Authorization", parameter.Header, parameter.WithRequired())),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.Users{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{
				response.New(models.FailResponse{}, "400", "Bad Request"),
				response.New(models.FailResponse{}, "500", "Internal Server")}),
		),

		endpoint.New(
			endpoint.DELETE,
			"/users/me",
			endpoint.WithTags("users"),
			endpoint.WithParams(parameter.StrParam("Authorization", parameter.Header, parameter.WithRequired())),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.SuccesResponse{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.FailResponse{}, "500", "Internal Server")}),
		),
	}

}
