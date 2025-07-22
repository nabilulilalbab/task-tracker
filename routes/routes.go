package routes

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/welcomesite/controllers"
)

func NewRouter(taskController *controllers.CarController, staticFS http.FileSystem) *httprouter.Router {
	router := httprouter.New()

	// Handler kustom untuk menyajikan file statis.
	// Ini akan menyajikan dari sistem file fisik untuk /static/uploads/*
	// dan dari embed.FS untuk semua path /static/* lainnya.
	fileHandler := http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "uploads/") {
			http.ServeFile(w, r, filepath.Join("static", r.URL.Path))
		} else {
			http.FileServer(staticFS).ServeHTTP(w, r)
		}
	}))
	router.Handler(http.MethodGet, "/static/*filepath", fileHandler)

	router.GET("/", taskController.ListTask)
	router.POST("/task/add", taskController.ProcessAddTask)
	router.POST("/task/update/:id", taskController.ProcessUpdateTask)
	router.POST("/task/delete/:id", taskController.DeleteTask)

	// Tambahkan route untuk WebSocket
	router.GET("/ws", taskController.HandleWebSocket)

	return router
}
