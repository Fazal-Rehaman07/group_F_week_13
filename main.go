package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

var db *sql.DB

// Connect to MySQL database
func initDatabase() {
	var err error
	conn := "admin:admin@4321@tcp(127.0.0.1:3306)/UserData" // Replace with your MySQL credentials
	db, err = sql.Open("mysql", conn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}
	fmt.Println("Connected to MySQL Database Successfully!")
}

// Handler for /current-time endpoint
func currentTimeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Username  string `json:"username"`
		IPAddress string `json:"ip_address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		log.Printf("Error decoding JSON: %v", err)
		return
	}

	if requestData.Username == "" || requestData.IPAddress == "" {
		http.Error(w, "Username and IP Address are required", http.StatusBadRequest)
		log.Println("Missing username or IP address in request")
		return
	}

	// Convert the current time to Toronto time zone
	location, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Error loading time zone", http.StatusInternalServerError)
		log.Printf("Error loading time zone: %v", err)
		return
	}
	torontoTime := time.Now().In(location).Format("2006-01-02 15:04:05")

	// Log data into database
	_, err = db.Exec("INSERT INTO time_log (username, ip_address, logged_time) VALUES (?, ?, ?)", requestData.Username, requestData.IPAddress, torontoTime)
	if err != nil {
		http.Error(w, "Error logging data to database", http.StatusInternalServerError)
		log.Printf("Error logging data to database: %v", err)
		return
	}

	// Respond with success message
	response := map[string]string{
		"message":      "Data logged successfully",
		"username":     requestData.Username,
		"ip_address":   requestData.IPAddress,
		"current_time": torontoTime,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

func main() {
	// Initialize the database
	initDatabase()
	defer db.Close()

	// Set up HTTP routes
	http.HandleFunc("/current-time", currentTimeHandler)

	fs := http.FileServer(http.Dir("./templates"))
	http.Handle("/", fs)

	// Start the server
	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
