package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/megadoge1337/maxon/internal/domain"
	"github.com/megadoge1337/maxon/internal/service"
	"github.com/megadoge1337/maxon/pkg/helper"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", h.GetAll)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetById)
	r.Get("/username/{username}", h.GetByUsername)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.DeleteById)

	return r
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll()
	if err != nil {
		slog.Error("GetAll failed", slog.Any("error", err))
		helper.WriteError(w, http.StatusInternalServerError, "failed to get users")
		return
	}

	helper.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.service.GetById(id)
	if err != nil {
		slog.Error("GetById failed", slog.Int("id", id), slog.Any("error", err))
		helper.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	helper.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	user, err := h.service.GetByUsername(username)
	if err != nil {
		slog.Error("GetByUsername failed", slog.String("username", username), slog.Any("error", err))
		helper.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	helper.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Create(input)
	if err != nil {
		slog.Error("Create failed", slog.Any("error", err))
		helper.WriteError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	helper.WriteJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input domain.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Update(input, id)
	if err != nil {
		slog.Error("Update failed", slog.Int("id", id), slog.Any("error", err))
		helper.WriteError(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	helper.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.service.DeleteById(id)
	if err != nil {
		slog.Error("DeleteById failed", slog.Int("id", id), slog.Any("error", err))
		helper.WriteError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	helper.WriteJSON(w, http.StatusOK, user)
}
