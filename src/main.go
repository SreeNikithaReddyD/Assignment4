package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// --- Product Data Structures (Matching Assignment Requirements) ---

// Product structure must match the required fields: Name, Price, Quantity.
type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name" binding:"required"` // Ensures 400 Bad Request if missing
	Price    float64 `json:"price" binding:"required"`
	Quantity int     `json:"quantity" binding:"required"`
}

// In-memory storage and mutex for safe concurrent access
var productStore = make(map[string]Product)
var idCounter int = 1
var mutex = &sync.RWMutex{}

// --- Main Setup ---

func main() {
	router := gin.Default()

	// API endpoints for product management
	router.POST("/products", postProduct)
	router.GET("/products/:id", getProductByID)

	// --- NEW: Add a route to serve the OpenAPI specification ---
	router.GET("/api.yaml", func(c *gin.Context) {
		c.File("./api.yaml") // Serves the api.yaml file from the same directory
	})

	log.Println("Starting server on :8080")
	router.Run(":8080")
}

// --- Handler Functions ---

// postProduct handles POST /products. Returns 201 Created or 400 Bad Request.
func postProduct(c *gin.Context) {
	var newProduct Product

	// 1. Bind and Basic Validation
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input. Name, Price, and Quantity are required fields."})
		return
	}

	// 2. Custom Validation
	if newProduct.Price <= 0 || newProduct.Quantity < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be positive, and Quantity cannot be negative."})
		return
	}

	// 3. Generate ID and Store
	mutex.Lock()
	newProduct.ID = strconv.Itoa(idCounter)
	idCounter++
	productStore[newProduct.ID] = newProduct
	mutex.Unlock()

	// 201 Created response
	c.JSON(http.StatusCreated, newProduct)
}

// getProductByID handles GET /products/{id}. Returns 200 OK or 404 Not Found.
func getProductByID(c *gin.Context) {
	id := c.Param("id")

	mutex.RLock()
	product, ok := productStore[id]
	mutex.RUnlock()

	if ok {
		// 200 OK response
		c.JSON(http.StatusOK, product)
		return
	}

	// 404 Not Found response
	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}