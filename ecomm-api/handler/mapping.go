package handler

import pb "github.com/fadhlinw/go-ecommerce/ecomm-grpc/pb"

func mapProductReqToPb(req *ProductReq) *pb.CreateProductRequest {
	return &pb.CreateProductRequest{
		Name:         req.Name,
		Image:        req.Image,
		Category:     req.Category,
		Description:  req.Description,
		Rating:       int32(req.Rating),
		NumReviews:   int32(req.NumReviews),
		Price:        req.Price,
		CountInStock: int32(req.CountInStock),
	}
}

func mapProductReqToUpdatePb(id int, req *ProductReq) *pb.UpdateProductRequest {
	return &pb.UpdateProductRequest{
		Id:           int32(id),
		Name:         req.Name,
		Image:        req.Image,
		Category:     req.Category,
		Description:  req.Description,
		Rating:       int32(req.Rating),
		NumReviews:   int32(req.NumReviews),
		Price:        req.Price,
		CountInStock: int32(req.CountInStock),
	}
}

func mapPbToProductRes(p *pb.Product) *ProductRes {
	if p == nil {
		return nil
	}
	return &ProductRes{
		ID:           int(p.Id),
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       int(p.Rating),
		NumReviews:   int(p.NumReviews),
		Price:        p.Price,
		CountInStock: int(p.CountInStock),
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func mapOrderReqToPb(req *OrderReq) *pb.CreateOrderRequest {
	var items []*pb.OrderItem
	for _, item := range req.Items {
		items = append(items, &pb.OrderItem{
			ProductId: int32(item.ProductID),
			Name:      item.Name,
			Quantity:  int32(item.Quantity),
			Image:     item.Image,
			Price:     int32(item.Price),
		})
	}
	return &pb.CreateOrderRequest{
		PaymentMethod: req.PaymentMethod,
		TaxPrice:      req.TaxPrice,
		ShippingPrice: req.ShippingPrice,
		TotalPrice:    req.TotalPrice,
		Items:         items,
	}
}

func mapOrderReqToUpdatePb(id int, req *OrderReq) *pb.UpdateOrderRequest {
	return &pb.UpdateOrderRequest{
		Id:            int32(id),
		PaymentMethod: req.PaymentMethod,
		TaxPrice:      req.TaxPrice,
		ShippingPrice: req.ShippingPrice,
		TotalPrice:    req.TotalPrice,
	}
}

func mapPbToOrderRes(o *pb.Order) *OrderRes {
	if o == nil {
		return nil
	}
	var items []OrderItemRes
	for _, item := range o.Items {
		items = append(items, OrderItemRes{
			ID:        int(item.Id),
			OrderID:   int(item.OrderId),
			ProductID: int(item.ProductId),
			Name:      item.Name,
			Quantity:  int(item.Quantity),
			Image:     item.Image,
			Price:     int(item.Price),
		})
	}
	return &OrderRes{
		ID:            int(o.Id),
		PaymentMethod: o.PaymentMethod,
		TaxPrice:      o.TaxPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
		Items:         items,
	}
}
