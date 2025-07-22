package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/welcomesite/controllers"
)

func NewRouter(taskController *controllers.CarController) *httprouter.Router {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.GET("/", taskController.ListTask)
	router.POST("/task/add", taskController.ProcessAddTask)
	router.POST("/task/update/:id", taskController.ProcessUpdateTask)
	router.POST("/task/delete/:id", taskController.DeleteTask)

	// Tambahkan route untuk WebSocket
	router.GET("/ws", taskController.HandleWebSocket)

	return router
}
