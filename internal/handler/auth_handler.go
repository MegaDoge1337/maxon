package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/dto"
	"megadoge1337/maxon/internal/usecase/auth"
	"megadoge1337/maxon/pkg/response"

	"github.com/go-chi/chi/v5"
)

type AuthHandlerDeps struct {
	LoginUC        auth.Login
	RegisterUC     auth.Register
	RefreshUC      auth.Refresh
	AuthMiddleware func(http.Handler) http.Handler
}

type AuthHandler struct {
	loginUC        auth.Login
	registerUC     auth.Register
	refreshUC      auth.Refresh
	authMiddleware func(http.Handler) http.Handler
}

func NewAuthHandler(deps AuthHandlerDeps) *AuthHandler {
	return &AuthHandler{
		loginUC:        deps.LoginUC,
		registerUC:     deps.RegisterUC,
		refreshUC:      deps.RefreshUC,
		authMiddleware: deps.AuthMiddleware,
	}
}

func (h *AuthHandler) RoutesV1() http.Handler {
	r := chi.NewRouter()

	// public
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)

	// auth protected
	r.Group(func(r chi.Router) {
		r.Use(h.authMiddleware)
		r.Post("/refresh", h.Refresh)
	})

	return r
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginDto dto.LoginDto
	if err := json.NewDecoder(r.Body).Decode(&loginDto); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	loginCommand := domain.LoginCommand{
		Username: loginDto.Username,
		Password: loginDto.Password,
	}

	access, refresh, err := h.loginUC.Execute(r.Context(), loginCommand)
	if err != nil {
		slog.Error("failed to login", slog.Any("error", err))
		response.WriteError(w, http.StatusBadRequest, "failed to login")
		return
	}

	tokensDto := dto.TokensDto{
		Access:  access,
		Refresh: refresh,
	}

	response.WriteJSON(w, http.StatusCreated, tokensDto)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerDto dto.RegisterDto
	if err := json.NewDecoder(r.Body).Decode(&registerDto); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	registerCommand := domain.RegisterCommand{
		Username: registerDto.Username,
		Email:    registerDto.Email,
		Password: registerDto.Password,
	}

	newUser, err := h.registerUC.Execute(r.Context(), registerCommand)
	if err != nil {
		slog.Error("failed to register", slog.Any("error", err))
		response.WriteError(w, http.StatusBadRequest, "failed to login")
		return
	}

	userDto := dto.UserDto{
		ID:       newUser.ID,
		Username: newUser.Username,
		Email:    newUser.Email,
		Created:  newUser.Created,
	}

	response.WriteJSON(w, http.StatusCreated, userDto)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var refreshSessionDto dto.RefreshDto
	if err := json.NewDecoder(r.Body).Decode(&refreshSessionDto); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	refreshCommand := domain.RefreshCommand{
		Refresh: refreshSessionDto.Refresh,
	}

	access, refresh, err := h.refreshUC.Execute(r.Context(), refreshCommand)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "failed to login")
		return
	}

	tokensDto := dto.TokensDto{
		Access:  access,
		Refresh: refresh,
	}

	response.WriteJSON(w, http.StatusCreated, tokensDto)
}
