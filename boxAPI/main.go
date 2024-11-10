package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type DataModel struct {
	Users        []User        `json:"users"`
	Transactions []Transaction `json:"transactions"`
	Boxes        []Box         `json:"boxes"`
}

type Box struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Balance     float64 `json:"balance"`
}

type Transaction struct {
	ID     string  `json:"id"`
	BoxID  string  `json:"boxId"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var dataModel = DataModel{
	Users:        make([]User, 0),
	Transactions: make([]Transaction, 0),
	Boxes:        make([]Box, 0),
}
var boxes = make(map[string]Box)
var transactions = make(map[string]Transaction)
var users = make(map[string]User)

// main initializes the HTTP server, sets up route handlers for the main API endpoints,
// and starts listening on port 8000. If the server fails to start, it logs an error message.
func main() {
	r := mux.NewRouter()

	r.HandleFunc("/all-data", getAllData).Methods("GET")

	r.HandleFunc("/boxes", getBoxes).Methods("GET")
	r.HandleFunc("/boxes/{id}", getBox).Methods("GET")
	r.HandleFunc("/boxes", createBox).Methods("POST")
	r.HandleFunc("/boxes/{id}", deleteBox).Methods("DELETE")
	r.HandleFunc("/boxes/{id}/balance", updateBoxBalance).Methods("PATCH")

	r.HandleFunc("/transactions", getTransactions).Methods("GET")
	r.HandleFunc("/transactions/{id}", getTransaction).Methods("GET")
	r.HandleFunc("/transactions", createTransaction).Methods("POST")

	r.HandleFunc("/users/{id}", getUserId).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users", getUsers).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", r))
}

// getAllData responds to GET requests to /all-data by sending the entire contents of dataModel as JSON.
func getAllData(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(dataModel)
}

// getBoxes responds to GET requests to /boxes by sending a JSON map of all boxes
// in the system, where each key is a box ID and the value is a Box struct.
func getBoxes(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(boxes)
}

// getBox responds to GET requests to /boxes/{id} by sending the details of the specified
// box as JSON. If the box is not found, it returns a 404 status with an error message.
func getBox(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	boxID := params["id"]
	box, ok := boxes[boxID]
	if !ok {
		http.Error(w, "Box not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(box)
}

// createBox responds to POST requests to /boxes by adding a new box to the system.
// The request body should contain a JSON object representing the box to add.
// The response is the newly-created box as JSON.
func createBox(w http.ResponseWriter, r *http.Request) {
	var box Box
	err := json.NewDecoder(r.Body).Decode(&box)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	boxes[box.ID] = box
	dataModel.Boxes = append(dataModel.Boxes, box)

	json.NewEncoder(w).Encode(box)
}

// updateBoxBalance responds to PATCH requests to /boxes/{id}/balance by updating the Balance
// field of the specified box. The request body should contain a JSON object with the new
// balance value. The response is the updated box as JSON. If the box is not found or the
// request body is invalid, it returns an error status with an error message.
func updateBoxBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boxID := vars["id"]

	box, ok := boxes[boxID]
	if !ok {
		http.Error(w, "Box not found", http.StatusNotFound)
		return
	}

	var update struct {
		Balance float64 `json:"balance"`
	}
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	box.Balance = update.Balance
	boxes[boxID] = box

	json.NewEncoder(w).Encode(box)
}

// deleteBox responds to DELETE requests to /boxes/{id} by deleting the box with the specified ID.
// The response is a JSON object with a message indicating success.
func deleteBox(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	boxID := params["id"]
	delete(boxes, boxID)
	json.NewEncoder(w).Encode(map[string]string{"message": "Box deleted successfully"})
}

// getTransactions responds to GET requests to /transactions by sending a JSON map of all
// transactions in the system, where each key is a transaction ID and the value is a
// Transaction struct.
func getTransactions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(transactions)
}

// getTransaction responds to GET requests to /transactions/{id} by sending the details of the specified transaction as JSON. If the transaction is not found, it returns a 404 status with an error message.
func getTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	transactionID := params["id"]
	transaction, ok := transactions[transactionID]
	if !ok {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(transaction)
}

// createTransaction responds to POST requests to /transactions by creating a new Transaction
// in the system. The request body should contain a JSON object with the details of the
// transaction. The response is the newly created transaction as JSON. If the request body
// is invalid, it returns a 400 status with an error message.
func createTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	transactions[transaction.ID] = transaction
	dataModel.Transactions = append(dataModel.Transactions, transaction)

	json.NewEncoder(w).Encode(transaction)
}

// createUser responds to POST requests to /users by creating a new user in the system.
// The request body should contain a JSON object representing the user to add.
// The response is the newly-created user as JSON. If the request body is invalid,
// it returns a 400 status with an error message.
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users[user.ID] = user
	dataModel.Users = append(dataModel.Users, user)

	json.NewEncoder(w).Encode(user)
}

// getUsers responds to GET requests to /users by returning a JSON list of all
// users in the system.
func getUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

// getUserId responds to GET requests to /users/{id} by sending the details of the
// specified user as JSON. If the user is not found, it returns a 404 status with
// an error message.
func getUserId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]
	user, ok := users[userID]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}
