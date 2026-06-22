package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	pb "github.com/fadhlinw/go-ecommerce/ecomm-grpc/pb"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	client pb.EcommServiceClient
}

func NewHandler(client pb.EcommServiceClient) *Handler {
	return &Handler{client: client}
}

// ---------------------------------------------------------
// Product Handlers
// ---------------------------------------------------------

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req ProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pbReq := mapProductReqToPb(&req)
	res, err := h.client.CreateProduct(r.Context(), pbReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(IDRes{ID: int(res.Id)})
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	res, err := h.client.GetProduct(r.Context(), &pb.GetProductRequest{Id: int32(id)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapPbToProductRes(res.Product))
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	res, err := h.client.ListProducts(r.Context(), &pb.ListProductsRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := []ProductRes{}
	for _, p := range res.Products {
		out = append(out, *mapPbToProductRes(p))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req ProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pbReq := mapProductReqToUpdatePb(id, &req)
	_, err = h.client.UpdateProduct(r.Context(), pbReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DefaultRes{Message: "Product updated"})
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = h.client.DeleteProduct(r.Context(), &pb.DeleteProductRequest{Id: int32(id)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DefaultRes{Message: "Product deleted"})
}

// ---------------------------------------------------------
// Order Handlers
// ---------------------------------------------------------

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req OrderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pbReq := mapOrderReqToPb(&req)
	res, err := h.client.CreateOrder(r.Context(), pbReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(IDRes{ID: int(res.Id)})
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	res, err := h.client.GetOrder(r.Context(), &pb.GetOrderRequest{Id: int32(id)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapPbToOrderRes(res.Order))
}

func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	res, err := h.client.ListOrders(r.Context(), &pb.ListOrdersRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := []OrderRes{}
	for _, o := range res.Orders {
		out = append(out, *mapPbToOrderRes(o))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req OrderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pbReq := mapOrderReqToUpdatePb(id, &req)
	_, err = h.client.UpdateOrder(r.Context(), pbReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DefaultRes{Message: "Order updated"})
}

func (h *Handler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = h.client.DeleteOrder(r.Context(), &pb.DeleteOrderRequest{Id: int32(id)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DefaultRes{Message: "Order deleted"})
}
