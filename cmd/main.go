package main

import (
	"log"
	"net/http"

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
	// Task
	taskRepo := repositories.NewTaskRepository(config.DB)
	taskService := services.NewTaskService(taskRepo)
	taskCtrl := controllers.NewTaskController(taskService, cachedTemplates)
	router := routes.NewRouter(taskCtrl)

	port := ":8080"
	log.Printf("Server berjalan di http://localhost%s\n", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err)
	}
}
