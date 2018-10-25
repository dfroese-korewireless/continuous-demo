package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dfroese-korewireless/continuous-demo/messages"
	"github.com/dfroese-korewireless/continuous-demo/storage"
	"github.com/gorilla/mux"
)

// Accessor is the accessor for the api
type Accessor interface {
	CreateMessage(http.ResponseWriter, *http.Request)
	GetMessages(http.ResponseWriter, *http.Request)
	GetMessage(http.ResponseWriter, *http.Request)
}

type context struct {
	storage.Database
}

// New returns a new Context that can be used for the API calls
func New(db storage.Database) Accessor {
	return &context{
		Database: db,
	}
}

// CreateMessage stores a message in the database
func (ctx *context) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var msg messages.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		fmt.Printf("decoding request body: %s\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := ctx.Database.StoreMessage(msg)
	if err != nil {
		fmt.Printf("storing message in database: %s\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct{ ID uint64 }{ID: id})
}

// GetMessages retrieves all the messages from the database
func (ctx *context) GetMessages(w http.ResponseWriter, r *http.Request) {
	msgs, err := ctx.Database.GetAllMessages()
	if err != nil {
		fmt.Printf("getting messages from database: %s\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(msgs)
}

func (ctx *context) GetMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	msgID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Printf("converting id to int64: %s\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	msg, err := ctx.Database.GetMessage(msgID)
	if err != nil {
		fmt.Printf("retrieving message from database: %s\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(msg)
}
