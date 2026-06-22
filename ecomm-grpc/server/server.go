package server

import (
	"context"

	pb "github.com/fadhlinw/go-ecommerce/ecomm-grpc/pb"
	"github.com/fadhlinw/go-ecommerce/ecomm-grpc/storer"
)

type EcommServer struct {
	pb.UnimplementedEcommServiceServer
	store *storer.Storer
}

func NewEcommServer(s *storer.Storer) *EcommServer {
	return &EcommServer{store: s}
}

// ---------------------------------------------------------
// Product Methods
// ---------------------------------------------------------

func (s *EcommServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	p := mapPbCreateReqToStorerProduct(req)
	id, err := s.store.CreateProduct(ctx, p)
	if err != nil {
		return nil, err
	}
	return &pb.CreateProductResponse{Id: int32(id)}, nil
}

func (s *EcommServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := s.store.GetProductByID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{Product: mapStorerProductToPb(p)}, nil
}

func (s *EcommServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := s.store.ListProducts(ctx)
	if err != nil {
		return nil, err
	}
	var pbProducts []*pb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, mapStorerProductToPb(&p))
	}
	return &pb.ListProductsResponse{Products: pbProducts}, nil
}

func (s *EcommServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	p := mapPbUpdateReqToStorerProduct(req)
	err := s.store.UpdateProduct(ctx, p)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateProductResponse{}, nil
}

func (s *EcommServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	err := s.store.DeleteProduct(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteProductResponse{}, nil
}

// ---------------------------------------------------------
// Order Methods
// ---------------------------------------------------------

func (s *EcommServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	o, items := mapPbCreateReqToStorerOrder(req)
	id, err := s.store.CreateOrder(ctx, o, items)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{Id: int32(id)}, nil
}

func (s *EcommServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	o, err := s.store.GetOrderByID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	items, err := s.store.GetOrderItemsByOrderID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.GetOrderResponse{Order: mapStorerOrderToPb(o, items)}, nil
}

func (s *EcommServer) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.store.ListOrders(ctx)
	if err != nil {
		return nil, err
	}
	var pbOrders []*pb.Order
	for _, o := range orders {
		pbOrders = append(pbOrders, mapStorerOrderToPb(&o, nil))
	}
	return &pb.ListOrdersResponse{Orders: pbOrders}, nil
}

func (s *EcommServer) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	o := mapPbUpdateReqToStorerOrder(req)
	err := s.store.UpdateOrder(ctx, o)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateOrderResponse{}, nil
}

func (s *EcommServer) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	err := s.store.DeleteOrder(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteOrderResponse{}, nil
}
