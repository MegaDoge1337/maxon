package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/dto"
	"megadoge1337/maxon/internal/service"
	"megadoge1337/maxon/pkg/response"

	"github.com/go-chi/chi/v5"
)

type UserHandlerDeps struct {
	UserService    service.UserService
	AuthMiddleware func(http.Handler) http.Handler
}

type UserHandler struct {
	service        service.UserService
	authMiddleware func(http.Handler) http.Handler
}

func NewUserHandler(deps UserHandlerDeps) *UserHandler {
	return &UserHandler{
		service:        deps.UserService,
		authMiddleware: deps.AuthMiddleware,
	}
}

func (h *UserHandler) RoutesV1() http.Handler {
	r := chi.NewRouter()

	// public
	r.Get("/username/{username}", h.GetByUsername)

	// auth protected
	r.Group(func(r chi.Router) {
		r.Use(h.authMiddleware)

		r.Get("/", h.GetAll)
		r.Post("/", h.Create)
		r.Get("/{id}", h.GetById)
		r.Put("/{id}", h.UpdateById)
		r.Delete("/{id}", h.DeleteById)
	})

	return r
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll(r.Context())
	if err != nil {
		slog.Error("GetAll failed", slog.Any("error", err))
		response.WriteError(w, http.StatusInternalServerError, "failed to get users")
		return
	}

	var usersDtos []dto.UserDto

	for _, user := range users {
		usersDtos = append(usersDtos, dto.UserDto{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Created:  user.Created,
		})
	}

	response.WriteJSON(w, http.StatusOK, usersDtos)
}

func (h *UserHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := response.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.service.GetById(r.Context(), id)
	if err != nil {
		slog.Error("GetById failed", slog.Int("id", id), slog.Any("error", err))
		response.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	userDto := dto.UserDto{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Created:  user.Created,
	}

	response.WriteJSON(w, http.StatusOK, userDto)
}

func (h *UserHandler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	user, err := h.service.GetByUsername(r.Context(), username)
	if err != nil {
		slog.Error("GetByUsername failed", slog.String("username", username), slog.Any("error", err))
		response.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	userDto := dto.UserDto{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Created:  user.Created,
	}

	response.WriteJSON(w, http.StatusOK, userDto)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createUserDto dto.CreateUserDto
	if err := json.NewDecoder(r.Body).Decode(&createUserDto); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	createUser := domain.CreateUserCommand{
		Username: createUserDto.Username,
		Email:    createUserDto.Email,
		Password: createUserDto.Password,
	}

	newUser, err := h.service.Create(r.Context(), createUser)
	if err != nil {
		slog.Error("Create failed", slog.Any("error", err))
		response.WriteError(w, http.StatusInternalServerError, "failed to create user")
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

func (h *UserHandler) UpdateById(w http.ResponseWriter, r *http.Request) {
	id, err := response.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var updateUserDto dto.UpdateUserDto
	if err := json.NewDecoder(r.Body).Decode(&updateUserDto); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updateUserCommand := domain.UpdateUserCommand{
		Username:    updateUserDto.Username,
		Email:       updateUserDto.Email,
		OldPassword: updateUserDto.OldPassword,
		NewPassword: updateUserDto.NewPassword,
	}

	user, err := h.service.UpdateById(updateUserCommand, id)
	if err != nil {
		slog.Error("Update failed", slog.Int("id", id), slog.Any("error", err))
		response.WriteError(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	userDto := dto.UserDto{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Created:  user.Created,
	}

	response.WriteJSON(w, http.StatusOK, userDto)
}

func (h *UserHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	id, err := response.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.service.DeleteById(id)
	if err != nil {
		slog.Error("DeleteById failed", slog.Int("id", id), slog.Any("error", err))
		response.WriteError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	userDto := dto.UserDto{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Created:  user.Created,
	}

	response.WriteJSON(w, http.StatusOK, userDto)
}
