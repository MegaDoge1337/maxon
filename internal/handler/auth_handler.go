package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/megadoge1337/maxon/internal/domain"
	"github.com/megadoge1337/maxon/internal/dto"
	"github.com/megadoge1337/maxon/internal/service"
	"github.com/megadoge1337/maxon/pkg/helper"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) RoutesV1() http.Handler {
	r := chi.NewRouter()

	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)

	return r
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var createSessionDto dto.CreateSessionDto
	if err := json.NewDecoder(r.Body).Decode(&createSessionDto); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	createSession := domain.CreateSessionCommand{
		Login:    createSessionDto.Login,
		Password: createSessionDto.Password,
	}

	access, refresh, err := h.service.Login(createSession)
	if err != nil {
		slog.Error("failed to login", slog.Any("error", err))
		helper.WriteError(w, http.StatusBadRequest, "failed to login")
		return
	}

	tokensDto := dto.TokensDto{
		Access:  access,
		Refresh: refresh,
	}

	helper.WriteJSON(w, http.StatusCreated, tokensDto)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var refreshSessionDto dto.RefreshSessionDto
	if err := json.NewDecoder(r.Body).Decode(&refreshSessionDto); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	refreshSession := domain.RefreshSessionCommand{
		Refresh: refreshSessionDto.Refresh,
	}

	access, refresh, err := h.service.Refresh(refreshSession)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "failed to login")
		return
	}

	tokensDto := dto.TokensDto{
		Access:  access,
		Refresh: refresh,
	}

	helper.WriteJSON(w, http.StatusCreated, tokensDto)
}
