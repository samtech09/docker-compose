package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sam-learn/docker-compose/reposit"
	"strconv"
	"time"
)

// DBhost=localhost DBname=postgres DBuser=postgres DBpwd=dockpgr DBport=5432 go run .

func main() {
	// Configure Logging
	LOG_FILE_LOCATION := os.Getenv("LOG_FILE_LOCATION")
	if LOG_FILE_LOCATION != "" {
		f, err := os.OpenFile(LOG_FILE_LOCATION, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)
	}

	// database config
	host := os.Getenv("DBhost")
	dbname := os.Getenv("DBname")
	dbuser := os.Getenv("DBuser")
	dbpwd := os.Getenv("DBpwd")
	dbport, err := strconv.Atoi(os.Getenv("DBport"))
	if err != nil {
		log.Fatalf("Invalid database port %v", err)
	}

	retry := 0

retrydb:
	retry++
	err = reposit.InitDB(host, dbname, dbuser, dbpwd, dbport)
	if err != nil {
		log.Printf("Database Connection error, will retry: %v\n", err)
		if retry < 10 {
			// wait then retry to allow database container to be ready
			time.Sleep(1 * time.Second)
			goto retrydb
		}
		log.Println("Give Up! failed to connect to database")
		return
	}
	log.Println("Database connection initialized.")

	helloHandler := http.HandlerFunc(hello)
	dataHandler := http.HandlerFunc(getData)
	http.HandleFunc("/", index)
	//http.HandleFunc("/db", dbinfo)
	http.HandleFunc("/healthcheck", health)
	http.HandleFunc("/readiness", ready)
	http.Handle("/hello", loggingHandler(helloHandler))
	http.Handle("/data", loggingHandler(dataHandler))

	fmt.Println("Server listening on port 3001...")
	http.ListenAndServe(":3001", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bare minimum API server in go with JSON and Database calls")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body><h1>Hello from golang!</h1></body></html>")
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func ready(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getData(w http.ResponseWriter, r *http.Request) {
	lst, err := reposit.GetData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// set header first before writing to response
	w.Header().Set("Content-Type", "application/json")

	// write data to response
	json.NewEncoder(w).Encode(lst)
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

// func dbinfo(w http.ResponseWriter, r *http.Request) {
// 	// host := os.Getenv("DBhost")
// 	// dbname := os.Getenv("DBname")
// 	// dbuser := os.Getenv("DBuser")
// 	// dbpwd := os.Getenv("DBpwd")
// 	// dbport, err := strconv.Atoi(os.Getenv("DBport"))
// 	// if err != nil {
// 	// 	log.Fatalf("Invalid database port %v", err)
// 	// }

// 	// err = reposit.InitDB(host, dbname, dbuser, dbpwd, dbport)
// 	// if err != nil {
// 	// 	fmt.Fprintf(w, "Host: %s, Port: %s<br/>%v", os.Getenv("DBhost"), os.Getenv("DBport"), err)
// 	// 	return
// 	// }
// 	fmt.Fprintf(w, "Host: %s, Port: %s", os.Getenv("DBhost"), os.Getenv("DBport"))
// }
