package main

import (
	"log"
	"net/http"
	"ziparchive/handlers"
	"ziparchive/storage"
)

const (
	serverPort    = ":8080"
	maxTasks      = 3
	maxZipArchive = 3
)

func main() {
	manager := storage.NewZipArchiveManager()

	//Create Zip Archive - создание нового архива
	http.HandleFunc("/createziparchive", handlers.HttpCreateZipArchive(manager))

	//Zip Archives - посмотреть статусы имеющихся архивов
	http.HandleFunc("/ziparchives", handlers.HttpZipArchives(manager))

	//Zip Archive / ZA.ID - посмотреть статус архива ID
	http.HandleFunc("/ziparchive/", handlers.HttpZipArchive(manager))

	//DOWNLOAD / ZA.ID - скачать архив ID
	http.HandleFunc("/download/", handlers.HttpDownload(manager))

	//Create Task / ZA.ID - создание нового задания в Zip Archive с ID
	http.HandleFunc("/createtask/", handlers.HttpCreateTask(manager))

	//Task / task.ID - посмотреть статус таска ID
	http.HandleFunc("/task/", handlers.HttpTask(manager))

	server := &http.Server{
		Addr:    serverPort,
		Handler: http.DefaultServeMux,
	}

	log.Println("Server started on port", serverPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}
