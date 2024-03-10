package main

import (
	"Rapid/src/database"
	"fmt"
	"os"
	"strings"
)

// Main method for runnning the system
func main() {
	database.InitializeDatabase()
	database.HandleAccountStartup()

	_, user, _, _ := database.GetUserDetails()
	fmt.Printf("Currently Logged in as %s\n", user)

	// Creates a splice from the command line input
	splice := os.Args[1:]
	arguments := strings.Join(splice, " ")
	for {
		// Retreives flags called and then runs the commands inside of them
		fmt.Scanln(&arguments)
		if strings.Compare(arguments, "quit") == 0 {
			fmt.Println("You have chosen to exit the program")
			break
		}
		// Retreives flags called and then runs the commands inside of them
		flags := retrieveFlags(arguments)
		command(flags, user)
		fmt.Print(">> ")
	}
}
