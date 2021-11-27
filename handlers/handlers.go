package handlers

import (
	"encoding/json"
	"gorm/db"
	"gorm/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

//Obtiene todos los datos en general
func GetUsers(rw http.ResponseWriter, r *http.Request) {
	users := models.Users{}
	//recupera los datos de la variable user
	db.Database.Find(&users)
	//y responde status ok
	sendData(rw, users, http.StatusOK)
}

func GetUser(rw http.ResponseWriter, r *http.Request) {
	//Capturamos un error
	if user, err := getUserById(r); err != nil {
		sendError(rw, http.StatusNotFound)
	} else {
		// Responswrite, user, y el statis Ok
		sendData(rw, user, http.StatusOK)
	}
}

func getUserById(r *http.Request) (models.User, *gorm.DB) {
	//Obtener ID
	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["id"])
	user := models.User{}
	//Capturamos un error si el usuario existe o no
	//Si no ha habido un error nos regresa nil, y si hay un error nos regresara el error
	if err := db.Database.First(&user, userId); err.Error != nil {
		return user, err
	} else {
		return user, nil
	}
}

//Crea un usuario
func CreateUser(rw http.ResponseWriter, r *http.Request) {
	user := models.User{}
	decoder := json.NewDecoder(r.Body)

	//Verifica si puede haber un error si lo hay devuelve el http
	if err := decoder.Decode(&user); err != nil {
		sendError(rw, http.StatusUnprocessableEntity)
	} else {
		//Sino lo tiene es el estatus que se ah creado
		db.Database.Save(&user)
		sendData(rw, user, http.StatusCreated)
	}
}

//Acuraliza un usuario
func UpdateUser(rw http.ResponseWriter, r *http.Request) {
	//Obtener el usuario por ID
	var userId int64

	if user_ant, err := getUserById(r); err != nil { //Capturamos un error
		sendError(rw, http.StatusNotFound)
	} else {
		userId = user_ant.Id
		user := models.User{}
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&user); err != nil {
			sendError(rw, http.StatusUnprocessableEntity)
		} else {
			user.Id = userId
			db.Database.Save(&user)
			sendData(rw, user, http.StatusOK)
		}
	}
}

func DeleteUser(rw http.ResponseWriter, r *http.Request) {
	if user, err := getUserById(r); err != nil { //Captura un error
		sendError(rw, http.StatusNotFound) //Si lo hay regresa el estatus
	} else {
		db.Database.Delete(&user)         //Sino lo hay localiza y elimina
		sendData(rw, user, http.StatusOK) //Y manda el mensaje de que se realizo de manera correcta
	}
}
