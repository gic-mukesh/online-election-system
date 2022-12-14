package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"online-election-system/auth"
	"online-election-system/dao"
	"online-election-system/helper"
	"online-election-system/model"
)

var uad = dao.UserDAO{}

func init() {
	uad.Server = "mongodb://localhost:27017/"
	uad.Database = "ElectionDB"
	uad.Collection = "User"

	uad.Connect()
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var dataBody model.User
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	if err := uad.Insert(dataBody); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request")
	} else {
		helper.RespondWithJson(w, http.StatusAccepted, "Record inserted successfully", true, nil)
	}
}

func VerifyUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var dataBody model.User
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	if err := uad.Update(dataBody); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request")
	} else {
		helper.RespondWithJson(w, http.StatusAccepted, "User verified successfully", true, nil)
	}
}

func SearchOneUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var dataBody model.User
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	id := dataBody.ID.String()
	user, err := uad.FindById(id)

	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request")
	} else {
		helper.RespondWithJson(w, http.StatusAccepted, "Record found successfully", true, user)
	}
}

func SearchMultipleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var dataBody model.User
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	id := dataBody.ID.String()
	user, err := uad.FindById(id)

	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request")
	} else {
		helper.RespondWithJson(w, http.StatusAccepted, "Record found successfully", true, user)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var dataBody model.User
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	id := dataBody.ID.String()

	if err := uad.Delete(id); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request")
	} else {
		helper.RespondWithJson(w, http.StatusAccepted, "Record Deleted successfully", true, nil)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var dataBody model.User
	if err := json.NewDecoder(r.Body).Decode(&dataBody); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	users, err := uad.FindByEmailAndPassword(dataBody.Email, dataBody.Password)

	if err != nil {
		http.Error(w, "No user found", http.StatusNotFound)
		return
	}

	user := users[0]

	fmt.Println(user)

	claim := auth.User_Claims{Username: user.Email, Authorized: true, Role: user.Role}

	jwt, err := auth.GenerateJWT(claim)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Token", jwt)

	helper.RespondWithJson(w, http.StatusOK, "Logged in successfully", true, nil)
}
