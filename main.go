package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"

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

var produk = []Produk{
	{ID: 1, Nama: "Indomie Goreng", Harga: 3500, Stok: 100},
	{ID: 2, Nama: "Teh Botol", Harga: 3000, Stok: 50},
	{ID: 3, Nama: "Kecap Bango", Harga: 12000, Stok: 20},
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
	// Load environment variable
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	// Buat Data Model untuk menyimpan config variable
	type Config struct {
		Port    string `mapstructure:"PORT"`
		DBConn string `mapstructure:"DB_CONN"`
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	//Setup database
	db, err := database.InitDB(config.DBConn)
	fmt.Println("DB Connection String:", config.DBConn)
	if config.DBConn == "" {
		log.Fatal("DB_CONN environment variable is not set")
	}
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()


	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Setup routes
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)
	http.HandleFunc("/api/health", healthCheck)
	// Swagger UI
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Server running di http://localhost:" + config.Port)
	fmt.Println("Swagger UI: http://localhost:" + config.Port + "/swagger/index.html")

	err = http.ListenAndServe(":"+ config.Port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
