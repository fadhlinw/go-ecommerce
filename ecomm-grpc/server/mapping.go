package server

import (
	"time"

	pb "github.com/fadhlinw/go-ecommerce/ecomm-grpc/pb"
	"github.com/fadhlinw/go-ecommerce/ecomm-grpc/storer"
)

// ---------------------------------------------------------
// Product Mapping
// ---------------------------------------------------------

func mapPbCreateReqToStorerProduct(req *pb.CreateProductRequest) *storer.Product {
	return &storer.Product{
		Name:         req.Name,
		Image:        req.Image,
		Category:     req.Category,
		Description:  req.Description,
		Rating:       int(req.Rating),
		NumReviews:   int(req.NumReviews),
		Price:        req.Price,
		CountInStock: int(req.CountInStock),
	}
}

func mapPbUpdateReqToStorerProduct(req *pb.UpdateProductRequest) *storer.Product {
	return &storer.Product{
		ID:           int(req.Id),
		Name:         req.Name,
		Image:        req.Image,
		Category:     req.Category,
		Description:  req.Description,
		Rating:       int(req.Rating),
		NumReviews:   int(req.NumReviews),
		Price:        req.Price,
		CountInStock: int(req.CountInStock),
	}
}

func mapStorerProductToPb(p *storer.Product) *pb.Product {
	var created, updated string
	if p.CreatedAt != nil {
		created = p.CreatedAt.Format(time.RFC3339)
	}
	if p.UpdatedAt != nil {
		updated = p.UpdatedAt.Format(time.RFC3339)
	}
	return &pb.Product{
		Id:           int32(p.ID),
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       int32(p.Rating),
		NumReviews:   int32(p.NumReviews),
		Price:        p.Price,
		CountInStock: int32(p.CountInStock),
		CreatedAt:    created,
		UpdatedAt:    updated,
	}
}

// ---------------------------------------------------------
// Order Mapping
// ---------------------------------------------------------

func mapPbCreateReqToStorerOrder(req *pb.CreateOrderRequest) (*storer.Order, []storer.OrderItem) {
	o := &storer.Order{
		PaymentMethod: req.PaymentMethod,
		TaxPrice:      req.TaxPrice,
		ShippingPrice: req.ShippingPrice,
		TotalPrice:    req.TotalPrice,
	}
	var items []storer.OrderItem
	for _, item := range req.Items {
		items = append(items, storer.OrderItem{
			ProductID: int(item.ProductId),
			Name:      item.Name,
			Quantity:  int(item.Quantity),
			Image:     item.Image,
			Price:     int(item.Price),
		})
	}
	return o, items
}

func mapPbUpdateReqToStorerOrder(req *pb.UpdateOrderRequest) *storer.Order {
	return &storer.Order{
		ID:            int(req.Id),
		PaymentMethod: req.PaymentMethod,
		TaxPrice:      req.TaxPrice,
		ShippingPrice: req.ShippingPrice,
		TotalPrice:    req.TotalPrice,
	}
}

func mapStorerOrderToPb(o *storer.Order, items []storer.OrderItem) *pb.Order {
	var pbItems []*pb.OrderItem
	for _, item := range items {
		pbItems = append(pbItems, &pb.OrderItem{
			Id:        int32(item.ID),
			OrderId:   int32(item.OrderID),
			ProductId: int32(item.ProductID),
			Name:      item.Name,
			Quantity:  int32(item.Quantity),
			Image:     item.Image,
			Price:     int32(item.Price),
		})
	}
	var created, updated string
	if o.CreatedAt != nil {
		created = o.CreatedAt.Format(time.RFC3339)
	}
	if o.UpdatedAt != nil {
		updated = o.UpdatedAt.Format(time.RFC3339)
	}
	return &pb.Order{
		Id:            int32(o.ID),
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		CreatedAt:     created,
		UpdatedAt:     updated,
		Items:         pbItems,
	}
}
