package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sangkips/revenue-system/internal/config"
	"github.com/sangkips/revenue-system/internal/db"
	"github.com/sangkips/revenue-system/internal/domain/counties"
	"github.com/sangkips/revenue-system/internal/domain/revenue"
	"github.com/sangkips/revenue-system/internal/domain/taxpayers"
	"github.com/sangkips/revenue-system/internal/domain/user"
	"github.com/sangkips/revenue-system/internal/middleware/auth"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := config.Load()
	sqlDB := db.ConnectAndMigrate(cfg.DBURL)
	defer sqlDB.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	usersHandler := user.NewHandler(sqlDB)
	r.Route("/users", func(r chi.Router) {
		r.Use(auth.JWTAuth(cfg.JWTSecret))
		usersHandler.RegisterUserRoutes(r)
	})

	countiesHandler := counties.NewHandler(sqlDB)
	r.Route("/counties", func(r chi.Router) {
		r.Use(auth.JWTAuth(cfg.JWTSecret))
		countiesHandler.RegisterCountyRoutes(r)
	})

	taxpayerHandler := taxpayers.NewHandler(sqlDB)
	r.Route("/taxpayers", func(r chi.Router) {
		r.Use(auth.JWTAuth(cfg.JWTSecret))
		taxpayerHandler.RegisterTaxpayerRoutes(r)
	})

	revenueHandler := revenue.NewHandler(sqlDB)
	r.Route("/revenues", func(r chi.Router) {
		r.Use(auth.JWTAuth(cfg.JWTSecret))
		revenueHandler.RegisterRevenueRoutes(r)
	})

	authService := auth.NewAuthService(user.NewRepository(sqlDB), cfg.JWTSecret)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			var req auth.LoginRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			resp, err := authService.Login(r.Context(), req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			json.NewEncoder(w).Encode(resp)
		})

		r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			var req auth.RegisterRequest

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			user, err := authService.Register(r.Context(), req)
			if err != nil {
				log.Error().Err(err).Msg("Failed to register user")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			type UserResponse struct {
				ID        string `json:"id"`
				Email     string `json:"email"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Role      string `json:"role"`
			}
			resp := UserResponse{
				ID:        user.ID.String(),
				Email:     user.Email,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Role:      user.Role,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(resp)
		})
	})

	log.Info().Msgf("Server starting on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}

}
