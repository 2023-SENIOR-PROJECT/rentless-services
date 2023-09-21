package main

import (
	"context"
	"log"
	"net"
	pb "product-module/product"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"google.golang.org/grpc"
)

// Define your Product struct
type Product struct {
	gorm.Model
	Name        string
	Price       float32
	Description string
	Owner       string
}

type server struct {
	db *gorm.DB
	pb.UnimplementedProductServiceServer
}

func (s *server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	product := &Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Owner:       req.Owner,
	}
	s.db.Create(product)
	log.Println("New product has been created")
	return productToProto(product), nil
}

func (s *server) ReadProduct(ctx context.Context, req *pb.ReadProductRequest) (*pb.Product, error) {
	var product Product
	s.db.First(&product, req.Id)
	log.Println("Product has been read")
	return productToProto(&product), nil
}

func (s *server) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	var product Product
	s.db.First(&product, req.Id)

	product.Name = req.Name
	product.Price = req.Price
	product.Description = req.Description
	product.Owner = req.Owner
	log.Println("Product has been updated")

	s.db.Save(&product)

	return productToProto(&product), nil
}

func (s *server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	var product Product
	s.db.First(&product, req.Id)

	if product.ID == 0 {
		return &pb.DeleteProductResponse{Success: false}, nil
	}

	log.Println("Product has been deleted")

	s.db.Delete(&product)
	return &pb.DeleteProductResponse{Success: true}, nil
}

func productToProto(product *Product) *pb.Product {
	return &pb.Product{
		Id:          uint32(product.ID),
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Owner:       product.Owner,
	}
}

func main() {
	// Specify the SQLite connection string
	db, err := gorm.Open("sqlite3", "rentless.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	db.AutoMigrate(&Product{})

	// mongo.ConnectMongoDB()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	service := &server{db: db}
	pb.RegisterProductServiceServer(s, service)

	log.Println("Starting gRPC server on :50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
