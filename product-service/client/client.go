package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	pb "product-module/product"
)

func main() {
	// Set up a connection to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := pb.NewProductServiceClient(conn)

	// Call the CreateProduct RPC to create a product
	createProductRequest := &pb.CreateProductRequest{
		Name:        "Sample Product",
		Price:       19.99,
		Description: "This is a sample product",
		Owner:       "John Doe",
	}
	createProductResponse, err := client.CreateProduct(context.Background(), createProductRequest)
	if err != nil {
		log.Fatalf("CreateProduct failed: %v", err)
	}
	log.Printf("CreateProduct response: %v", createProductResponse)

	// Call the ReadProduct RPC to retrieve the product's ID
	readProductRequest := &pb.ReadProductRequest{Id: createProductResponse.Id}
	readProductResponse, err := client.ReadProduct(context.Background(), readProductRequest)
	if err != nil {
		log.Fatalf("ReadProduct failed: %v", err)
	}
	log.Printf("ReadProduct response: %v", readProductResponse)

	// Call the UpdateProduct RPC to update the product
	updateProductRequest := &pb.UpdateProductRequest{
		Id:          readProductResponse.Id,
		Name:        "Updated Product",
		Price:       29.99,
		Description: "This is an updated product",
		Owner:       "Jane Doe",
	}
	updateProductResponse, err := client.UpdateProduct(context.Background(), updateProductRequest)
	if err != nil {
		log.Fatalf("UpdateProduct failed: %v", err)
	}
	log.Printf("UpdateProduct response: %v", updateProductResponse)

	// Call the DeleteProduct RPC to delete the product
	deleteProductRequest := &pb.DeleteProductRequest{Id: updateProductResponse.Id}
	deleteProductResponse, err := client.DeleteProduct(context.Background(), deleteProductRequest)
	if err != nil {
		log.Fatalf("DeleteProduct failed: %v", err)
	}

	if deleteProductResponse.Success {
		log.Printf("DeleteProduct successful")
	} else {
		log.Printf("DeleteProduct failed: Product not found")
	}

	// Call the GetAllProducts RPC to retrieve all products
	getAllProductsRequest := &pb.GetAllProductsRequest{}
	getAllProductsResponse, err := client.GetAllProducts(context.Background(), getAllProductsRequest)
	if err != nil {
		log.Fatalf("GetAllProducts failed: %v", err)
	}

	// Print the list of products
	log.Println("List of all products:")
	for _, product := range getAllProductsResponse.Products {
		log.Printf("Product ID: %s, Name: %s, Price: %.2f, Description: %s, Owner: %s",
			product.Id, product.Name, product.Price, product.Description, product.Owner)
	}
}
