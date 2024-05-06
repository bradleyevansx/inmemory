package main


func main() {
	server := NewServer(":3000")
	db, err := ConnectDB()
	if err != nil {
		return
	}
	repository := NewProfileService(db)
	controller := NewAPIController(server, repository)
	controller.RegisterRoutes()
	server.Run()
}

//dbservice which interacts with db

//controller which recieves requests and then talks to dbservice