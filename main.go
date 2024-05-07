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