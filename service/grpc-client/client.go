package main

import (
	"context"
	"log"
	"net/http"

	pb "rentless-services/service/product-service/product"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

const serverAddress = "localhost:50051" // Replace with the actual server address

func main() {
	r := gin.Default()

	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewProductServiceClient(conn)

	// Endpoint to create a new product
	r.POST("/api/products", func(c *gin.Context) {
		var product pb.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createProductResponse, err := createProduct(client, &product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createProductResponse)
	})

	// Endpoint to get all products
	r.GET("/api/products", func(c *gin.Context) {
		getAllProductsResponse, err := getAllProducts(client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, getAllProductsResponse)
	})

	// Endpoint to get a single product
	r.GET("/api/products/:id", func(c *gin.Context) {
		productID := c.Param("id")
		readProductResponse, err := readProduct(client, productID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusOK, readProductResponse)
	})

	// Endpoint to update a product
	r.PUT("/api/products/:id", func(c *gin.Context) {
		productID := c.Param("id")
		var product pb.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		product.Id = productID

		updateProductResponse, err := updateProduct(client, &product)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusOK, updateProductResponse)
	})

	// Endpoint to delete a product
	r.DELETE("/api/products/:id", func(c *gin.Context) {
		productID := c.Param("id")
		deleteProductResponse, err := deleteProduct(client, productID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusOK, deleteProductResponse)
	})

	// Start the Gin server
	err = r.Run(":8081") // Replace with the desired port
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}

func createProduct(client pb.ProductServiceClient, product *pb.Product) (*pb.Product, error) {
	createProductRequest := &pb.CreateProductRequest{
		Name:         product.Name,
		Slug:         product.Slug,
		Image:        product.Image,
		Category:     product.Category,
		Brand:        product.Brand,
		Price:        product.Price,
		CountInStock: product.CountInStock,
		Description:  product.Description,
		Owner:        product.Owner,
	}

	createProductResponse, err := client.CreateProduct(context.Background(), createProductRequest)
	return createProductResponse, err
}

func getAllProducts(client pb.ProductServiceClient) (*pb.GetAllProductsResponse, error) {
	getAllProductsRequest := &pb.GetAllProductsRequest{}
	getAllProductsResponse, err := client.GetAllProducts(context.Background(), getAllProductsRequest)
	return getAllProductsResponse, err
}

func readProduct(client pb.ProductServiceClient, productID string) (*pb.Product, error) {
	readProductRequest := &pb.ReadProductRequest{
		Id: productID,
	}
	readProductResponse, err := client.ReadProduct(context.Background(), readProductRequest)
	return readProductResponse, err
}

func updateProduct(client pb.ProductServiceClient, product *pb.Product) (*pb.Product, error) {
	updateProductRequest := &pb.UpdateProductRequest{
		Id:   product.Id,
		Name: product.Name, // Update other fields as needed
	}
	updateProductResponse, err := client.UpdateProduct(context.Background(), updateProductRequest)
	return updateProductResponse, err
}

func deleteProduct(client pb.ProductServiceClient, productID string) (*pb.DeleteProductResponse, error) {
	deleteProductRequest := &pb.DeleteProductRequest{
		Id: productID,
	}
	deleteProductResponse, err := client.DeleteProduct(context.Background(), deleteProductRequest)
	return deleteProductResponse, err
}
