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

func get_employees(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(employeedb)
}

func get_single_employee(w http.ResponseWriter,r *http.Request){
	id := r.PathValue("id")
	employee,exist := employeedb[id]
	if !exist {
		http.Error(w,"not found",http.StatusNotFound)
		return
	}
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(employee)
}

func update_employee(w http.ResponseWriter,r *http.Request){
	id := r.PathValue("id")
	_,exist := employeedb[id]
	if !exist{
		http.Error(w,"employee not exist",http.StatusBadRequest)
		return
	}

	var updated_data employee
	if err:=json.NewDecoder(r.Body).Decode(&updated_data);err != nil{
		http.Error(w,"Failed to decode the Request",http.StatusBadRequest)
		return
	}
	updated_data.ID = id
	employeedb[id] = updated_data
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated_data)
}

func delete_employee(w http.ResponseWriter, r *http.Request){
	id := r.PathValue("id")
	_,exist := employeedb[id] 
	if !exist {
		http.Error(w,"Data not found",http.StatusBadRequest)
		return
	}
	delete(employeedb,id)
	w.WriteHeader(http.StatusOK)
	response_message := map[string]string{"message":"deleted succfully"}
	json.NewEncoder(w).Encode(response_message)

}
func testing_server(w http.ResponseWriter,r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("this is my first server"))
}

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("GET /",testing_server)
	mux.HandleFunc("POST /employee/add",create_employee)
	mux.HandleFunc("GET /employees/{$}",get_employees)
	mux.HandleFunc("GET /employee/{id}",get_single_employee)
	mux.HandleFunc("PUT /employee/update/{id}",update_employee)
	mux.HandleFunc("DELETE /employee/remove/{id}",delete_employee)
	http.ListenAndServe(":8080",mux)
}