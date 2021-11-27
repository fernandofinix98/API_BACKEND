package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//Conexion y verificacion con la base de datos
var dsn = "root:022528Sepmaydic@tcp(localhost:3306)/pia_b?charset=utf8mb4&parseTime=True&loc=Local"

//funcion anonima para la prueba de la conexion
var Database = func() (db *gorm.DB) {
	//Abre la concexion con la funcion Open y dentro de la funcion open abre con gorm el dsn
	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		fmt.Println("Error de conexion", err)
		panic(err)
	} else {
		fmt.Println("Conexion exitosa")
		return db
	}
}()
