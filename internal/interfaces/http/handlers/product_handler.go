package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GitAlex9/go-order-service/internal/application/contracts"
	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/interfaces/http/response"
)

type ProductHandler struct {
	service contracts.ProductService
}

func NewProductHandler(service contracts.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	product, err := h.service.Create(r.Context(), req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	product, err := h.service.Get(r.Context(), id)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	offset, limit := parsePagination(r)

	products, err := h.service.List(r.Context(), offset, limit)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, products)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var req dto.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	product, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *ProductHandler) IncreaseStock(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var req dto.AdjustStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	product, err := h.service.IncreaseStock(r.Context(), id, req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) DecreaseStock(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var req dto.AdjustStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	product, err := h.service.DecreaseStock(r.Context(), id, req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	product, err := h.service.Activate(r.Context(), id)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	product, err := h.service.Deactivate(r.Context(), id)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, product)
}
