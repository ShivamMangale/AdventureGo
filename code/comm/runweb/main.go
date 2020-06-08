package main

import (
	"AdventureGo/code/greet"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func displayError(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func main() {
	fmt.Println("Started")
	filename := flag.String("filename", "storydata.json", "To change the json source file for the story data used. \nThe file should be of the format .json .\n(default is storydata.json)")
	port := flag.Int("port", 3000, "Assign a different port for the web app to run on. \n(default is 3030)")
	templatetype := flag.Int("template", 1, "Choose the layout of webpage you want to use. (default:1)")
	flag.Parse()
	fmt.Println("The story file chosen is", *filename)
	fmt.Println("The port chosen is", *port)
	fmt.Println("The template chosen is", *templatetype)

	fileToRead, err := os.Open(*filename)
	if err != nil {
		displayError(fmt.Sprint("The file", *filename, "cannot be read"))
	}

	story, err := greet.SourceToJSON(fileToRead)

	// fmt.Printf("%+v\n", story)

	// h := greet.NewHandler(story, greet.ChooseTempl(*templatetype), greet.ChoosePathFn(newPathFn))
	h := greet.NewHandler(story, greet.ChooseTempl(*templatetype))
	fmt.Println("Starting the server at:", *port)
	// addr := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))

	fmt.Println("Finished")
}

func newPathFn(r *http.Request) string {
	pat := r.URL.Path
	curpath := strings.TrimSpace(pat)
	if curpath == "/story" || curpath == "/story/" {
		curpath = "/story/intro"
	}
	curpath = curpath[len("/story/"):]

	return curpath
}
