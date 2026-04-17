package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/megadoge1337/maxon/internal/domain"
	"github.com/megadoge1337/maxon/internal/dto"
	"github.com/megadoge1337/maxon/internal/service"
	"github.com/megadoge1337/maxon/pkg/response"
)

type AuthHandlerDeps struct {
	AuthService    service.AuthService
	AuthMiddleware func(http.Handler) http.Handler
}

type AuthHandler struct {
	service        service.AuthService
	authMiddleware func(http.Handler) http.Handler
}

func NewAuthHandler(deps AuthHandlerDeps) *AuthHandler {
	return &AuthHandler{
		service:        deps.AuthService,
		authMiddleware: deps.AuthMiddleware,
	}
}

func (h *AuthHandler) RoutesV1() http.Handler {
	r := chi.NewRouter()

	// public
	r.Post("/login", h.Login)

	// auth protected
	r.Group(func(r chi.Router) {
		r.Use(h.authMiddleware)

		r.Post("/refresh", h.Refresh)
	})

	return r
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var createSessionDto dto.CreateSessionDto
	if err := json.NewDecoder(r.Body).Decode(&createSessionDto); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	createSession := domain.CreateSessionCommand{
		Login:    createSessionDto.Login,
		Password: createSessionDto.Password,
	}

	access, refresh, err := h.service.Login(createSession)
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

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var refreshSessionDto dto.RefreshSessionDto
	if err := json.NewDecoder(r.Body).Decode(&refreshSessionDto); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	refreshSession := domain.RefreshSessionCommand{
		Refresh: refreshSessionDto.Refresh,
	}

	access, refresh, err := h.service.Refresh(refreshSession)
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
