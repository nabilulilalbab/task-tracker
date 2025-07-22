package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"

	"github.com/nabilulilalbab/welcomesite/models"
	"github.com/nabilulilalbab/welcomesite/services"
	"github.com/nabilulilalbab/welcomesite/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type CarController struct {
	service  services.TaskService
	template *template.Template
}

func NewTaskController(service services.TaskService, tmpl *template.Template) *CarController {
	return &CarController{service: service, template: tmpl}
}

func (c *CarController) ListTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tasks, err := c.service.GetAllTasks()
	if err != nil {
		http.Error(w, "gagal ambil task nih", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title":     "Home Task",
		"Tasks":     tasks,
		"OS":        utils.GetOS(),
		"Terminals": utils.GetAvailableTerminals(),
	}

	if err := c.template.ExecuteTemplate(w, "indextask.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
}

func (c *CarController) ProcessAddTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("Tidak dapat mem-parsing multipart form: %v", err)
		http.Error(w, "Request tidak valid", http.StatusBadRequest)
		return
	}
	file, fileHeader, err := r.FormFile("cover")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("Gagal mengambil file: %v", err)
		http.Error(w, "Gagal memproses file cover", http.StatusInternalServerError)
		return
	}
	if file != nil {
		defer file.Close()
	}
	task := &models.Task{
		Judul:   r.FormValue("judul"),
		Tipe:    r.FormValue("tipe"),
		Tags:    r.FormValue("tags"),
		Catatan: r.FormValue("catatan"),
		Status:  "todo",
	}
	pathProjectVal := r.FormValue("path_project")
	if pathProjectVal != "" {
		task.PathProject = &pathProjectVal
	}
	linkWebsiteVal := r.FormValue("link_website")
	if linkWebsiteVal != "" {
		task.LinkWebsite = &linkWebsiteVal
	}
	_, err = c.service.CreateTask(task, fileHeader)
	if err != nil {
		log.Printf("Error saat memanggil service CreateTask: %v", err)
		http.Error(w, "Gagal menyimpan data task", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *CarController) HandleWebSocket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Gagal upgrade ke WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Gagal membaca pesan WebSocket:", err)
			break
		}

		var msg map[string]string
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Gagal unmarshal JSON:", err)
			continue
		}

		terminalCmd := msg["terminal"]
		path := msg["path"]

		if err := utils.OpenTerminal(terminalCmd, path); err != nil {
			log.Printf("Gagal membuka terminal: %v", err)
		} else {
			log.Printf("Berhasil membuka %s di %s", path, terminalCmd)
		}
	}
}

func (c *CarController) ProcessUpdateTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Request tidak valid", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("cover")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Gagal memproses file cover", http.StatusInternalServerError)
		return
	}
	if file != nil {
		defer file.Close()
	}

	taskInput := &models.Task{
		Judul:   r.FormValue("judul"),
		Tipe:    r.FormValue("tipe"),
		Tags:    r.FormValue("tags"),
		Catatan: r.FormValue("catatan"),
		Status:  r.FormValue("status"),
	}
	pathProjectVal := r.FormValue("path_project")
	if pathProjectVal != "" {
		taskInput.PathProject = &pathProjectVal
	}
	linkWebsiteVal := r.FormValue("link_website")
	if linkWebsiteVal != "" {
		taskInput.LinkWebsite = &linkWebsiteVal
	}
	_, err = c.service.UpdateTask(uint(id), taskInput, fileHeader)
	if err != nil {
		log.Printf("Error saat memanggil service UpdateTask: %v", err)
		http.Error(w, "Gagal mengupdate data task", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
