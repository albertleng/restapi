package main

//TODO: Implement unit testing
import "log"

func main() {

	a := App{}
	a.Initialize()
	log.Println("Marvel REST Api starts...")
	a.Run(":8080")
}
