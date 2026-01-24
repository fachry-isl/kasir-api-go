package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "kasir-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Kasir API
// @version         1.0
// @description     API untuk sistem kasir sederhana dengan Category
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@kasir.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// Produk struct
type Produk struct {
	ID    int    `json:"id" example:"1"`
	Nama  string `json:"nama" example:"Indomie Goreng"`
	Harga int    `json:"harga" example:"3500"`
	Stok  int    `json:"stok" example:"100"`
}

// Category struct
type Category struct {
	ID          int    `json:"id" example:"1"`
	Name        string `json:"name" example:"Makanan"`
	Description string `json:"description" example:"Kategori untuk produk makanan"`
}

var produk = []Produk{
	{ID: 1, Nama: "Indomie Goreng", Harga: 3500, Stok: 100},
	{ID: 2, Nama: "Teh Botol", Harga: 3000, Stok: 50},
	{ID: 3, Nama: "Kecap Bango", Harga: 12000, Stok: 20},
}

var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Kategori produk makanan"},
	{ID: 2, Name: "Minuman", Description: "Kategori produk minuman"},
	{ID: 3, Name: "Bumbu", Description: "Kategori bumbu dapur"},
}

// ==================== CATEGORY HANDLERS ====================

// GetAllCategories godoc
// @Summary      Get all categories
// @Description  Get list of all categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200  {array}   Category
// @Router       /api/categories [get]
func getAllCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// CreateCategory godoc
// @Summary      Create new category
// @Description  Add a new category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category  body      Category  true  "Category data"
// @Success      201       {object}  Category
// @Failure      400       {object}  map[string]string
// @Router       /api/categories [post]
func createCategory(w http.ResponseWriter, r *http.Request) {
	var categoryBaru Category
	err := json.NewDecoder(r.Body).Decode(&categoryBaru)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	categoryBaru.ID = len(categories) + 1
	categories = append(categories, categoryBaru)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(categoryBaru)
}

// GetCategoryByID godoc
// @Summary      Get category by ID
// @Description  Get single category by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  Category
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/categories/{id} [get]
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}

// UpdateCategory godoc
// @Summary      Update category
// @Description  Update category by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id        path      int       true  "Category ID"
// @Param        category  body      Category  true  "Category data"
// @Success      200       {object}  Category
// @Failure      400       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Router       /api/categories/{id} [put]
func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updateData Category
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for i := range categories {
		if categories[i].ID == id {
			updateData.ID = id
			categories[i] = updateData
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateData)
			return
		}
	}

	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}

// DeleteCategory godoc
// @Summary      Delete category
// @Description  Delete category by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/categories/{id} [delete]
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Category berhasil dihapus",
			})
			return
		}
	}

	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}

// ==================== PRODUK HANDLERS (existing) ====================

// GetAllProduk godoc
// @Summary      Get all products
// @Description  Get list of all products
// @Tags         produk
// @Accept       json
// @Produce      json
// @Success      200  {array}   Produk
// @Router       /api/produk [get]
func getAllProduk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produk)
}

// CreateProduk godoc
// @Summary      Create new product
// @Description  Add a new product to inventory
// @Tags         produk
// @Accept       json
// @Produce      json
// @Param        produk  body      Produk  true  "Product data"
// @Success      201     {object}  Produk
// @Failure      400     {object}  map[string]string
// @Router       /api/produk [post]
func createProduk(w http.ResponseWriter, r *http.Request) {
	var produkBaru Produk
	err := json.NewDecoder(r.Body).Decode(&produkBaru)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	produkBaru.ID = len(produk) + 1
	produk = append(produk, produkBaru)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(produkBaru)
}

// GetProdukByID godoc
// @Summary      Get product by ID
// @Description  Get single product by ID
// @Tags         produk
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  Produk
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/produk/{id} [get]
func getProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

// UpdateProduk godoc
// @Summary      Update product
// @Description  Update product by ID
// @Tags         produk
// @Accept       json
// @Produce      json
// @Param        id      path      int     true  "Product ID"
// @Param        produk  body      Produk  true  "Product data"
// @Success      200     {object}  Produk
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Router       /api/produk/{id} [put]
func updateProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updateData Produk
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for i := range produk {
		if produk[i].ID == id {
			updateData.ID = id
			produk[i] = updateData
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateData)
			return
		}
	}

	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

// DeleteProduk godoc
// @Summary      Delete product
// @Description  Delete product by ID
// @Tags         produk
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/produk/{id} [delete]
func deleteProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, p := range produk {
		if p.ID == id {
			produk = append(produk[:i], produk[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Produk berhasil dihapus",
			})
			return
		}
	}

	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Check if API is running
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /api/health [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "API Running",
	})
}

func main() {
	// ========== CATEGORY ROUTES ==========
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}
	})

	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getAllCategories(w, r)
		} else if r.Method == "POST" {
			createCategory(w, r)
		}
	})

	// ========== PRODUK ROUTES ==========
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProdukByID(w, r)
		} else if r.Method == "PUT" {
			updateProduk(w, r)
		} else if r.Method == "DELETE" {
			deleteProduk(w, r)
		}
	})

	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getAllProduk(w, r)
		} else if r.Method == "POST" {
			createProduk(w, r)
		}
	})

	http.HandleFunc("/api/health", healthCheck)

	// Swagger UI
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Port untuk Zeabur (environment variable)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running di http://localhost:" + port)
	fmt.Println("Swagger UI: http://localhost:" + port + "/swagger/index.html")

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
