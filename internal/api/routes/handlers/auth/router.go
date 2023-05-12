package auth

import (
	"github.com/gofiber/fiber/v2"

	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/user"
	"github.com/romankravchuk/muerta/internal/services/auth"
)

func NewRouter(cfg *config.Config, client repositories.PostgresClient, logger *log.Logger, jware *jware.JWTMiddleware) *fiber.App {
	repo := user.New(client)
	svc := auth.New(cfg, repo)
	r := fiber.New()
	h := New(cfg, svc, logger)
	r.Post("/sign-up", h.SignUp)
	r.Post("/login", h.Login)
	r.Post("/logout", jware.DeserializeUser, h.Logout)
	r.Post("/refresh", jware.DeserializeUser, h.RefreshAccessToken)
	return r
}
