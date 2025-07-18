package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"ziparchive/storage"
)

const (
	serverPort    = ":8080"
	maxTasks      = 3
	maxZipArchive = 3
)

func CheckExtension(url string) bool {
	splt := strings.Split(url, ".")
	switch {
	case splt[len(splt)-1] == "pdf":
		return true
	case splt[len(splt)-1] == "jpeg":
		return true
	default:
		return false
	}
}

func main() {
	manager := storage.NewZipArchiveManager()

	//Create Zip Archive - создание нового архива
	http.HandleFunc("/createziparchive", func(w http.ResponseWriter, r *http.Request) {
		if manager.InProgress == maxZipArchive {
			http.Error(w, "Exceeded the number of ZipArchives", http.StatusBadRequest)
			return
		}
		za := manager.CreateZipArchive()
		manager.InProgress += 1
		fmt.Fprintf(w, "ZipArchiveID: %s, Status: %s", za.ID, za.Status)
	})

	//Zip Archives - посмотреть статусы имеющихся архивов
	http.HandleFunc("/ziparchives", func(w http.ResponseWriter, r *http.Request) {
		for _, v := range manager.ZipArchives {
			fmt.Fprintf(w, "ZipArchiveID: %s, Status: %s, TasksNumber: %v\n", v.ID, v.Status, len(v.Tasks))
		}
	})

	//Zip Archive / ZA.ID - посмотреть статус архива ID
	http.HandleFunc("/ziparchive/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/ziparchive/"):]
		if id == "" {
			http.Error(w, "Empty Zip Archive ID", http.StatusBadRequest)
			return
		}
		za, exist := manager.ZipArchives[id]
		if !exist {
			http.Error(w, "Invalid Zip Archive ID", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "ZipArchiveID: %s, Status: %s, TasksNumber: %v\n", za.ID, za.Status, len(za.Tasks))
		for _, v := range za.Tasks {
			fmt.Fprintf(w, "TaskID: %s, Status: %s\n", v.ID, v.Status)
		}
		if len(za.Tasks) == maxTasks {
			fmt.Fprintf(w, "To DOWNLOAD ZipArchive: %s\n", strings.Join([]string{"http://localhost:8080/download/", za.ID}, ""))
		}

	})

	//DOWNLOAD / ZA.ID - скачать архив ID
	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/download/"):]
		if id == "" {
			http.Error(w, "Empty Zip Archive ID", http.StatusBadRequest)
			return
		}
		za, exist := manager.ZipArchives[id]
		if !exist {
			http.Error(w, "Invalid Zip Archive ID", http.StatusBadRequest)
			return
		}
		if len(za.Tasks) == maxTasks {
			w.Header().Set("Content-Type", "application/zip")
			w.Header().Set("Content-Disposition", strings.Join([]string{"attachment; filename=", "archive", za.ID, ".zip"}, ""))

			if _, err := io.Copy(w, za.ZipBuffer); err != nil {
				log.Println("Error sending response:", err)
			}
		} else {
			http.Error(w, "ZipArchive is not ready for DOWNLOAD yet", http.StatusBadRequest)
		}
	})

	//Create Task / ZA.ID - создание нового задания в Zip Archive с ID
	http.HandleFunc("/createtask/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := r.URL.Path[len("/createtask/"):]
		if id == "" {
			http.Error(w, "Empty Zip Archive ID", http.StatusBadRequest)
			return
		}
		za, exist := manager.ZipArchives[id]
		if !exist {
			http.Error(w, "Invalid Zip Archive ID", http.StatusBadRequest)
			return
		}
		if len(za.Tasks) == maxTasks {
			http.Error(w, "Exceeded the number of tasks", http.StatusBadRequest)
			return
		}
		link := r.FormValue("url")
		ok := CheckExtension(link)
		if !ok {
			http.Error(w, "Wrong Extension", http.StatusBadRequest)
			return
		}
		tsk := za.CreateTask(r)
		err := za.AddFileToZA(link, strings.Join([]string{strconv.Itoa(len(za.Tasks)), strings.Split(link, ".")[len(strings.Split(link, "."))-1]}, "."))
		if err != nil {
			http.Error(w, "URL is unavailable", http.StatusBadRequest)
			tsk.Status = storage.TaskStatusCompletedUnsuccessfully
		} else {
			tsk.Status = storage.TaskStatusCompletedSuccessfully
		}
		if len(za.Tasks) == maxTasks {
			za.ZipWriter.Close()
			za.Status = storage.ZAStatusCompletedSuccessfully
			manager.InProgress -= 1
		}
		fmt.Fprintf(w, "ZipArchiveID: %s, TaskID: %s, Status: %s\n", za.ID, tsk.ID, tsk.Status)
	})

	//Task / task.ID - посмотреть статус таска ID
	http.HandleFunc("/task/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/task/"):]
		if id == "" {
			http.Error(w, "Invalid Task ID", http.StatusBadRequest)
			return
		}
		if tsk, err := manager.FindTask(id); err != nil {
			http.Error(w, "Invalid Task ID", http.StatusBadRequest)
			return
		} else {
			fmt.Fprintf(w, "TaskID: %s, Status: %s\n", tsk.ID, tsk.Status)
		}
	})
	/////////////////////////////////////////////////////////////////////////////////////////////////////
	server := &http.Server{
		Addr:    serverPort,
		Handler: http.DefaultServeMux,
	}

	log.Println("Server started on port", serverPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}
