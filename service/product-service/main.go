package main

import (
	"context"
	"fmt"
	"log"
	"net"
	database "rentless-services/internal/infrastructure/product_database/mongo"
	pb "rentless-services/service/product-service/product"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

// Define your Product struct
type Product struct {
	// gorm.Model
	ID          string  `bson:"_id"`
	Name        string  `bson:"name"`
	Price       float32 `bson:"price"`
	Description string  `bson:"description"`
	Owner       string  `bson:"owner"`
}

type server struct {
	coll *mongo.Collection
	pb.UnimplementedProductServiceServer
}

func (s *server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	product := &Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Owner:       req.Owner,
	}
	product.ID = primitive.NewObjectID().Hex()
	_, err := s.coll.InsertOne(context.Background(), product)
	if err != nil {
		return nil, err
	}
	// product.ID = res.InsertedID.(string)
	log.Println("New product has been created")
	return productToProto(product), nil
}

func (s *server) GetAllProducts(ctx context.Context, req *pb.GetAllProductsRequest) (*pb.GetAllProductsResponse, error) {
	cursor, err := s.coll.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("could not fetch products: %v", err)
	}
	defer cursor.Close(context.Background())

	var products []*pb.Product
	for cursor.Next(context.Background()) {
		var product Product
		if err := cursor.Decode(&product); err != nil {
			return nil, fmt.Errorf("could not decode product: %v", err)
		}
		products = append(products, productToProto(&product))
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return &pb.GetAllProductsResponse{
		Products: products,
	}, nil
}

func (s *server) ReadProduct(ctx context.Context, req *pb.ReadProductRequest) (*pb.Product, error) {
	var product Product
	res := s.coll.FindOne(context.Background(), bson.M{"_id": req.Id})
	if err := res.Decode(&product); err != nil {
		return nil, err
	}
	// product.ID = product.ID
	log.Println("Product has been read", product)
	return productToProto(&product), nil
}

func (s *server) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	var product Product
	s.coll.FindOne(context.Background(), bson.M{"_id": product.ID})

	product.Name = req.Name
	product.Price = req.Price
	product.Description = req.Description
	product.Owner = req.Owner
	log.Println("Product has been updated")

	_, err := s.coll.UpdateOne(context.Background(), bson.M{"_id": req.Id}, bson.M{"$set": bson.M{"name": product.Name, "price": product.Price, "description": product.Description, "owner": product.Owner}})
	if err != nil {
		return nil, fmt.Errorf("could not update product: %v", err)
	}
	product.ID = req.Id

	return productToProto(&product), nil
}

func (s *server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	// var product Product
	log.Println("Product has been deleted")
	// s.coll.FindOne(context.Background(), bson.M{"_id": req.Id})

	if req.Id == "" {
		return &pb.DeleteProductResponse{Success: false}, nil
	}

	s.coll.DeleteOne(context.Background(), bson.M{"_id": req.Id})
	return &pb.DeleteProductResponse{Success: true}, nil
}

func productToProto(product *Product) *pb.Product {
	return &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Owner:       product.Owner,
	}
}

func main() {
	// Specify the SQLite connection string
	// db, err := gorm.Open("sqlite3", "rentless.db")
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }
	// defer db.Close()
	// db.AutoMigrate(&Product{})

	db, err := database.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Client().Disconnect(context.Background())

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	service := &server{coll: database.Collection}
	pb.RegisterProductServiceServer(s, service)

	log.Println("Starting gRPC server on :50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
