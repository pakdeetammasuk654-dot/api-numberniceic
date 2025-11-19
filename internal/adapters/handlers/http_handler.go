package handlers

import (
	"api-numberniceic/internal/core/ports"
	"encoding/json"
	"net/http"
	"strings"
)

type HttpHandler struct {
	service ports.NumberService
}

func NewHttpHandler(service ports.NumberService) *HttpHandler {
	return &HttpHandler{
		service: service,
	}
}

func (h *HttpHandler) CreateNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type requestBody struct {
		Number  string `json:"number"`
		Meaning string `json:"meaning"`
	}
	var req requestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	created, err := h.service.Create(req.Number, req.Meaning)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *HttpHandler) GetNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// สมมติ URL: /numbers?n=081xxxxxxx
	numberStr := r.URL.Query().Get("n")

	result, err := h.service.Get(numberStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
