package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"unicode"
	"github.com/gorilla/mux"
    "strings"
    "github.com/google/uuid"
)

type Item struct {
    ShortDescription string  `json:"shortDescription"`
    Price            string  `json:"price"`
}

type Receipt struct {
    Retailer     string `json:"retailer"`
    PurchaseDate string `json:"purchaseDate"`
    PurchaseTime string `json:"purchaseTime"`
    Items        []Item `json:"items"`
    Total        string `json:"total"`
}

type idMapStruct struct {
    idMap map[string]int
}

var myStruct = idMapStruct{
    idMap: make(map[string]int),
}

// processReceiptHandler handles POST requests to /receipts/process
func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var receipt Receipt
    err := json.NewDecoder(r.Body).Decode(&receipt)
    if err != nil {
        http.Error(w, "Error decoding JSON", http.StatusBadRequest)
        return
    }

    score := 0

    // Checking each character in retailer name
    for _, char := range receipt.Retailer {
        if unicode.IsLetter(char) || unicode.IsDigit(char) {
            score++
        }
    }

    cents := receipt.Total[len(receipt.Total)-2:]
    // Check if the total is a round dollar amount
    if cents == "00" {
        score += 50
    }

    floatTotal, err := strconv.ParseFloat(receipt.Total, 64)
    if err != nil {
        fmt.Println("Error converting string to float:", err)
        return
    }

    remainder := math.Mod(floatTotal, 0.25)

    // Check if the remainder is zero
    if remainder == 0 {
        score += 25
    }

    // Add a point for every 2 items on the receipt
    score += (len(receipt.Items) / 2) * 5

    // Looping over all the items, checking if the trimmed item description is divisble by 3, then doing the score calculation
    for _, item := range receipt.Items {
        if len(strings.Trim(item.ShortDescription, " ")) % 3 == 0 {
            floatPrice, err := strconv.ParseFloat(item.Price, 64)
            if err != nil {
                fmt.Println("Error converting string to float:", err)
                return
            }

            score += int(math.Ceil(floatPrice * 0.2))
        }
    }

    // Calculating if the day is odd
    day := receipt.PurchaseDate[len(receipt.PurchaseDate)-2:]
    intDay, err := strconv.Atoi(day)
    if err != nil {
        fmt.Println("Error converting string to int:", err)
        return
    }

    if intDay % 2 == 1 {
        score += 6
    }

    // Calculating if the time is between 2 and 4 pm
    time := strings.Split(receipt.PurchaseTime, ":")

    hour, err := strconv.Atoi(time[0])
    if err != nil {
        fmt.Println("Error converting string to int:", err)
        return
    }
    minutes, err := strconv.Atoi(time[1])
    if err != nil {
        fmt.Println("Error converting string to int:", err)
        return
    }

    if hour >= 14 && hour <= 15 {
        if (hour == 14 && minutes == 0){
            // Since time is exactly 2:00pm, don't add to score
        } else {
            score += 10
        }
    }

    id := uuid.New().String()

    myStruct.idMap[id] = score

    response := map[string]string{"id": id}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

// getPointsHandler handles GET requests to /receipts/{id}/points
func getPointsHandler(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    segments := strings.Split(path, "/")
    if len(segments) < 3 {
        http.Error(w, "Invalid URL path", http.StatusBadRequest)
        return
    }
    id := segments[2]
    if points, exists := myStruct.idMap[id]; exists {
        response := map[string]int{"points": points}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    } else {
        http.Error(w, "ID not found", http.StatusNotFound)
    }
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/receipts/process", processReceiptHandler).Methods("POST")
    r.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")
    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
