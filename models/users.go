package models

import (
	"gorm/db"
)

//Estrucura de la tabla en la base de datos
type User struct {
	Id       int64  `json:"id"`
	Empresa  string `json:"empresa"`
	Grupo    string `json:"grupo"`
	Miembros string `json:"miembros"`
}

type Users []User

//Conecta la base de datos y Podemos migrar una estructura en este caso "User"
func MigrarUser() {
	db.Database.AutoMigrate(User{})

}
