package main

import (
	"log"
	"net/http"

	"github.com/nabilulilalbab/welcomesite"
	"github.com/nabilulilalbab/welcomesite/config"
	"github.com/nabilulilalbab/welcomesite/controllers"
	"github.com/nabilulilalbab/welcomesite/repositories"
	"github.com/nabilulilalbab/welcomesite/routes"
	"github.com/nabilulilalbab/welcomesite/services"
	"github.com/nabilulilalbab/welcomesite/view"
)

func main() {
	config.InitDatabase()
	cachedTemplates := view.ParseTemplates()
	// Definisikan path untuk unggahan dan buat direktori jika belum ada
	uploadsPath := "static/uploads/tasks"
	// Task
	taskRepo := repositories.NewTaskRepository(config.DB)
	taskService := services.NewTaskService(taskRepo, uploadsPath)
	taskCtrl := controllers.NewTaskController(taskService, cachedTemplates)
	// Inisialisasi router dengan static file system
	router := routes.NewRouter(taskCtrl, welcomesite.StaticFS)

	port := ":8080"
	log.Printf("Server berjalan di http://localhost%s\n", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err)
	}
}
