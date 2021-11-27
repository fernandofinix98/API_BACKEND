### Base de Datos


En este apartado hacemos la conexión de la base de datos, usando gorm para la conexión con mysql, al igual incluimos la url de la base de datos realizada en mysql. En la variable Database() colocamos la conexión junto con una forma de verificar que efectivamente se realizó la conexión

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

##Datos de la tabla

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
    
Aquí podemos ver nuestras variables declaradas en el type para que se impriman en formato json, junto con la funcion Migrauser que es la que nos conecta la base de datos 

###CRUD
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

En la func GetUser obtendremos todos los datos en general que vayamos ingresando

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

En esta función el usuario obtendrá el dato que requiera, usando solamente su id

    
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
Aquí se están creando cada uno de los datos que irán en el registro

    
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
    

Esta función actualiza cualquier dato requerido, solamente hay que decirle el id del dato que se quiere actualizar 

    
    func DeleteUser(rw http.ResponseWriter, r *http.Request) {
    	if user, err := getUserById(r); err != nil { //Captura un error
    		sendError(rw, http.StatusNotFound) //Si lo hay regresa el estatus
    	} else {
    		db.Database.Delete(&user)         //Sino lo hay localiza y elimina
    		sendData(rw, user, http.StatusOK) //Y manda el mensaje de que se realizo de manera correcta
    	}
    }
Aquí se eliminarán cualquier dato que se necesite borrar, solo con buscar su id 

###Respuestas Postman
    import (
    	"encoding/json"
    	"fmt"
    	"net/http"
    )
    
    func sendData(rw http.ResponseWriter, data interface{}, status int) {
    	//modifica el estatus en el header
    	rw.Header().Set("Contet-Type", "aplication/json")
    	rw.WriteHeader(status)
    
    	//Responde al Cliente
    	//La data que recibo la convertimos en json
    	output, _ := json.Marshal(&data)
    	// y el json lo convierte en string
    	fmt.Fprintln(rw, string(output))
    }
    
    //
    func sendError(rw http.ResponseWriter, status int) {
    	//Responde con un status
    	rw.WriteHeader(status)
    	//Y devuelve el mensaje
    	fmt.Fprintln(rw, "Resouece Not Found")
    }
    
Esto se puede ver como un convertidor donde los datos de la base de datos, se estarán viendo en el Postman, la func SendData modifica el estatus en el header. La func sendError nos mandará un error en caso de que algo haya fallado 

###Ruteo de las funciones
    
    func main() {
    
    	//models.MigrarUser()
    
    	//Rutas
    	mux := mux.NewRouter()
    
    	//Funcion CORS
    	enableCORS(mux)
    
    	//Responder al cliente
    	mux.HandleFunc("/api/user/", handlers.GetUsers).Methods("GET")
    	mux.HandleFunc("/api/user/{id:[0-9]+}", handlers.GetUser).Methods("GET")
    	mux.HandleFunc("/api/user/", handlers.CreateUser).Methods("POST")
    	mux.HandleFunc("/api/user/{id:[0-9]+}", handlers.UpdateUser).Methods("PUT")
    	mux.HandleFunc("/api/user/{id:[0-9]+}", handlers.DeleteUser).Methods("DELETE")
    
    	//Servidor
    	fmt.Println("Run server: http://localhost:3000")
    	log.Fatal(http.ListenAndServe(":3000", mux))
    }
En el main encontraremos exclusivamente nuestros imports y los ruteos de las funciones crud que realizamos junto con los métodos que le asignamos 

###CORS
    func enableCORS(router *mux.Router) {
    	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    		w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
    	}).Methods(http.MethodOptions)
    	router.Use(middlewareCors)
    }
    
    func middlewareCors(next http.Handler) http.Handler {
    	return http.HandlerFunc(
    		func(w http.ResponseWriter, req *http.Request) {
    			w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
    			w.Header().Set("Access-Control-Allow-Credentials", "true")
    			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
    			next.ServeHTTP(w, req)
    		})
    }
En la sección del CORS encontramos la funcion enableCORS la que le daremos el acceso de nuestro localhost. En la funcion middlewareCORS incluiremos los metodos que el CORS aceptará
