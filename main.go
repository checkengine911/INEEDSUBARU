package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ImageResponse struct {
	URL string `json:"url"`
}

var images []string

func init() {
	// Получаем абсолютный путь к директории проекта
	execPath, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory:", err)
	}
	log.Printf("Current working directory: %s", execPath)

	// Проверяем существование директории images
	imagesDir := filepath.Join(execPath, "images")
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		log.Fatal("Images directory does not exist:", imagesDir)
	}

	files, err := os.ReadDir(imagesDir)
	if err != nil {
		log.Fatal("Error reading images directory:", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			images = append(images, "/images/"+file.Name())
			log.Printf("Added image: /images/%s", file.Name())
		}
	}

	if len(images) == 0 {
		log.Fatal("No images found in images directory")
	}

	log.Printf("Total images loaded: %d", len(images))
}

// Middleware для CORS
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func getRandomImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rand.Seed(time.Now().UnixNano())
	randomImage := images[rand.Intn(len(images))]
	log.Printf("Selected random image: %s", randomImage)

	response := ImageResponse{URL: randomImage}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Настраиваем обработку статических файлов
	imagesHandler := http.FileServer(http.Dir("images"))
	http.Handle("/images/", http.StripPrefix("/images/", imagesHandler))

	// Регистрируем обработчики с CORS middleware
	http.HandleFunc("/api/random-image", corsMiddleware(getRandomImage))

	// Настраиваем обработку статических файлов для остальных путей
	staticHandler := http.FileServer(http.Dir("static"))
	http.Handle("/", staticHandler)

	log.Println("Server starting on :8080")
	log.Println("Available images:")
	for _, img := range images {
		log.Printf("- %s", img)
	}

	// Запускаем сервер
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
