package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GitAlex9/go-order-service/internal/application/contracts"
	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/interfaces/http/response"
)

type CustomerHandler struct {
	service contracts.CustomerService
}

func NewCustomerHandler(service contracts.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	customer, err := h.service.Create(r.Context(), req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, customer)
}

func (h *CustomerHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	customer, err := h.service.Get(r.Context(), id)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	offset, limit := parsePagination(r)

	customers, err := h.service.List(r.Context(), offset, limit)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, customers)
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var req dto.UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	customer, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		response.HandleError(w, err)
		return
	}

	response.NoContent(w)
}
