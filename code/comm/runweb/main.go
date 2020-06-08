package main

import (
	"AdventureGo/code/greet"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func displayError(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func main() {
	fmt.Println("Started")
	filename := flag.String("filename", "storydata.json", "To change the json source file for the story data used. \nThe file should be of the format .json .\n(default is storydata.json)")
	port := flag.Int("port", 3000, "Assign a different port for the web app to run on. \n(default is 3030)")
	flag.Parse()
	fmt.Println("The story file chosen is", *filename)
	fmt.Println("The port chosen is", *port)

	fileToRead, err := os.Open(*filename)
	if err != nil {
		displayError(fmt.Sprint("The file", *filename, "cannot be read"))
	}

	story, err := greet.SourceToJSON(fileToRead)

	// fmt.Printf("%+v\n", story)

	h := greet.NewHandler(story)
	fmt.Println("Starting the server at:", *port)
	// addr := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))

	fmt.Println("Finished")
}
