package main

func init() {
	initializers.LoadENV()
	initializers.ConnectDatabase()
}
func main() {
	initializers.DB.AutoMigrate(&models.Patient{})
	initializers.DB.AutoMigrate(&models.User{})

}
