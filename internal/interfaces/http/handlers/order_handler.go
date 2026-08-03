package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GitAlex9/go-order-service/internal/application/contracts"
	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/interfaces/http/response"
)

type OrderHandler struct {
	service contracts.OrderService
}

func NewOrderHandler(service contracts.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	order, err := h.service.Create(r.Context(), req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	order, err := h.service.Get(r.Context(), id)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, order)
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	offset, limit := parsePagination(r)

	orders, err := h.service.List(r.Context(), offset, limit)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) Pay(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	order, err := h.service.Pay(r.Context(), id)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, order)
}

func (h *OrderHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	order, err := h.service.Cancel(r.Context(), id)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, order)
}

func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
