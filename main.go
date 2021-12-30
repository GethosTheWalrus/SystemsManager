package main

import (
	"ThePooReview/controllers"
	"ThePooReview/middleware"
	"ThePooReview/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

const (
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	dbuser := "app"
	dbpassword := "Apple123!123"
	dbserver := "127.0.0.1"
	dbport := "3306"

	db, err = gorm.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbserver+":"+dbport+")/db?parseTime=true")
	db.LogMode(true)

	if err != nil {
		log.Println("Could not connect to database: " + err.Error())
	} else {
		log.Println("Connected to database")
	}

	controllers.Db = db
	middleware.Db = db

	// auto migrate models to the DB
	db.AutoMigrate(&models.Session{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.HostSystem{})
	db.AutoMigrate(&models.Cpu{})
	db.AutoMigrate(&models.Memory{})
	db.AutoMigrate(&models.Network{})
	db.AutoMigrate(&models.General{})

	// listen for agent connections on a background thread
	go listenForAgents()

	// listen for HTTP API requests
	handleRequests()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to HomePage!")
}

func handleRequests() {
	log.Println("Starting development server at http://127.0.0.1:10000/")
	log.Println("Quit the server with CONTROL-C.")

	// creates a new instance of a mux router
	router := mux.NewRouter().StrictSlash(true)
	secureRouter := router.PathPrefix("/secure").Subrouter()
	openRouter := router.PathPrefix("/open").Subrouter()

	// assign middleware to the routers
	router.Use(middleware.HeaderMiddleware)
	secureRouter.Use(middleware.AuthenticatorMiddleware)

	/* Login */
	openRouter.HandleFunc("/User/Login", controllers.AuthenticateUser).Methods("POST")

	/* Users */
	openRouter.HandleFunc("/User", controllers.CreateUser).Methods("POST")
	openRouter.HandleFunc("/User", controllers.GetUsers).Methods("GET")
	openRouter.HandleFunc("/User/{userid}", controllers.GetUsers).Methods("GET")
	openRouter.HandleFunc("/User/{userid}", controllers.UpdateUser).Methods("PUT")
	openRouter.HandleFunc("/User/{userid}", controllers.DeleteUser).Methods("DELETE")

	/* Systems */
	openRouter.HandleFunc("/System", controllers.GetSystems).Methods("GET")

	/* Sessions */
	secureRouter.HandleFunc("/Session", controllers.CreateSession).Methods("POST")
	secureRouter.HandleFunc("/Session/{userid}", controllers.GetSession).Methods("GET")

	openRouter.HandleFunc("/", homePage)

	log.Fatal(http.ListenAndServe(":10000", router))
}

func listenForAgents() {

	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)

	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes.
	defer l.Close()

	fmt.Println("Listening for agent connections on " + CONN_HOST + ":" + CONN_PORT)

	for {

		// Listen for an incoming connection.
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			// os.Exit(1)
		}

		// Handle connections in a new goroutine.
		go handleAgentConnection(conn)

	}

}

func handleAgentConnection(con net.Conn) {

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	// Read the incoming connection into the buffer.
	_, err := con.Read(buf)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	// decode payload data
	system := models.HostSystem{}
	newSystem := models.HostSystem{}

	trimmedBuf := bytes.Trim(buf, "\x00")

	e := json.Unmarshal(trimmedBuf, &system)
	if e != nil {
		fmt.Println("invalid system data", e)
		return
	}

	// insert or update the system in the database
	db.
		Joins("INNER JOIN networks ON host_systems.network_id = networks.id").
		First(&newSystem, "hostname = ?", system.Network.Hostname)

	fmt.Println(newSystem.Id > 0)

	if newSystem.Id > 0 {

		var cpu models.Cpu
		var memory models.Memory
		var network models.Network
		var general models.General

		db.First(&cpu, "id = ?", newSystem.CpuId)
		db.First(&memory, "id = ?", newSystem.MemoryId)
		db.First(&network, "id = ?", newSystem.NetworkId)
		db.First(&general, "id = ?", newSystem.GeneralId)

		cpu = models.Cpu{
			Id:     cpu.Id,
			Idle:   system.Cpu.Idle,
			System: system.Cpu.System,
			Total:  system.Cpu.Total,
			User:   system.Cpu.User,
		}

		memory = models.Memory{
			Id:     memory.Id,
			Cached: system.Memory.Cached,
			Free:   system.Memory.Free,
			Total:  system.Memory.Total,
			Used:   system.Memory.Used,
		}

		network = models.Network{
			Id:          network.Id,
			Hostname:    system.Network.Hostname,
			PreferredIp: system.Network.PreferredIp,
		}

		general = models.General{
			Id:              general.Id,
			OperatingSystem: system.General.OperatingSystem,
			LastSeen:        system.General.LastSeen,
		}

		// cpu stats not implemented on some platforms
		if cpu.Idle+cpu.System+cpu.Total+cpu.User > 0 {
			db.Save(&cpu)
		}

		db.Save(&memory)
		db.Save(&network)
		db.Save(&general)

	} else {

		db.Create(&system)

	}

	// Send a response back to person contacting us.
	con.Write([]byte("Message received."))

	// Close the connection when you're done with it.
	defer con.Close()

}
