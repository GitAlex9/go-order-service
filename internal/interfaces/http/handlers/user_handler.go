package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GitAlex9/go-order-service/internal/application/contracts"
	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/interfaces/http/response"
)

type UserHandler struct {
	service contracts.UserService
}

func NewUserHandler(service contracts.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	user, err := h.service.Create(r.Context(), req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, user)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	user, err := h.service.Get(r.Context(), id)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	offset, limit := parsePagination(r)

	users, err := h.service.List(r.Context(), offset, limit)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, users)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := h.service.ChangePassword(r.Context(), id, req); err != nil {
		response.HandleError(w, err)
		return
	}

	response.NoContent(w)
}

func (h *UserHandler) ChangeEmail(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var req dto.ChangeUserEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	user, err := h.service.ChangeEmail(r.Context(), id, req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	user, err := h.service.Activate(r.Context(), id)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	user, err := h.service.Deactivate(r.Context(), id)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, user)
}
