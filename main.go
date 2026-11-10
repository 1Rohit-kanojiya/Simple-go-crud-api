package main

import (
	"net/http"
	"encoding/json"
	"strconv"
)

type employee struct{
	ID string `json:"id"`
	Name string `json:"name"`
	Age string `json:"age"`
}
var employeedb = map[string]employee{}

func create_employee(w http.ResponseWriter, r *http.Request){
	var employee employee

	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil{
		http.Error(w,"failed to decode your request",http.StatusBadRequest)
		return
	}

	if employee.ID == "" {
		employee.ID = strconv.Itoa(len(employeedb)+1)
	}

	_,exist := employeedb[employee.ID]

	if exist {
		http.Error(w,"data alredy exist",http.StatusConflict)
		return
	}

	employeedb[employee.ID] = employee;

	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(employee)
}

func testing_server(w http.ResponseWriter,r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("this is my first server"))
}

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("GET /",testing_server)
	mux.HandleFunc("POST /employee/add",create_employee)
	http.ListenAndServe(":8080",mux)
}