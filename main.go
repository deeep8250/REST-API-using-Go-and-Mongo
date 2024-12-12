package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var client *mongo.Client
var collection *mongo.Collection

func DB_connect() {

	var err error

	uri := "mongodb+srv://deep82500:deep82500@deep.jqe1i.mongodb.net/?retryWrites=true&w=majority&appName=deep"

	dblocation := options.Client().ApplyURI(uri)

	client, err = mongo.Connect(context.TODO(), dblocation)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("NGO").Collection("employee")
	fmt.Println("Successfully connected to the MongDB")

}

type data_Structure struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Position     string    `json:"position"`
	Joining_Time time.Time `json:"joining time"`
}

func addData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var Data data_Structure
	json.NewDecoder(r.Body).Decode(&Data)

	result, err := collection.InsertOne(context.TODO(), bson.M{
		"id":           Data.ID,
		"name":         Data.Name,
		"position":     Data.Position,
		"joining time": time.Now(),
	})

	if err != nil {
		http.Error(w, "data added unsuccesfull", http.StatusInternalServerError)
	}

	if result.InsertedID == 0 {
		http.Error(w, "No data added", http.StatusOK)
	} else {
		response := map[string]string{"message": "data added successfully"}
		json.NewEncoder(w).Encode(response)
	}

}

func FilterData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var showData data_Structure
	vars := mux.Vars(r)
	id := vars["id"]

	id1, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	err1 := collection.FindOne(context.TODO(), bson.M{"id": id1}).Decode(&showData)
	if err1 == mongo.ErrNoDocuments {
		http.Error(w, "No data found", http.StatusNotFound)
	} else {
		log.Print(err)
	}

	json.NewEncoder(w).Encode(showData)
}

func UpdateData(w http.ResponseWriter, r *http.Request) {

	var InputData map[string]interface{}

	vars := mux.Vars(r)
	id := vars["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "url input not converted into int", http.StatusInternalServerError)
	}

	filter := bson.M{"id": idInt}

	if err := json.NewDecoder(r.Body).Decode(&InputData); err != nil {
		http.Error(w, "user input from post men  problem", http.StatusBadRequest)
	}

	updateInputData := bson.M{"$set": InputData}

	err = collection.FindOne(context.TODO(), filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "no document matches", http.StatusOK)
		} else {
			http.Error(w, "find one problem ", http.StatusInternalServerError)
		}
		return
	}

	result, err := collection.UpdateOne(context.TODO(), filter, updateInputData)
	if err != nil {
		http.Error(w, "Update data problem", http.StatusInternalServerError)
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "no data updated", http.StatusOK)
	} else {
		response := map[string]string{"message": "data  updated successfully"}
		json.NewEncoder(w).Encode(response)
	}

}

func deleteData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/josn")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "cant convert the id into int", http.StatusInternalServerError)
	}

	filter := bson.M{"id": idInt}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "no data found", http.StatusOK)
		} else {
			http.Error(w, "error occured in delete one ", http.StatusInternalServerError)
		}
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "no data deleted", http.StatusOK)
	} else {
		response := map[string]string{"message": "data deleted successfully"}
		json.NewEncoder(w).Encode(response)
	}
}

func main() {

	DB_connect()
	r := mux.NewRouter()

	r.HandleFunc("/add", addData).Methods("POST")
	r.HandleFunc("/filter/{id}", FilterData).Methods("GET")
	r.HandleFunc("/update/{id}", UpdateData).Methods("PUT")
	r.HandleFunc("/delete/{id}", deleteData).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))

}
