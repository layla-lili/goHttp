package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/ascii-art", handleAsciiArt)

	//	http.ListenAndServe(":8080", nil)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("Templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("r.Method", r.Method)
	if r.Method == http.MethodPost {
		// Get the form values from the request
		text := r.FormValue("text")
		bannerFilename := r.FormValue("banner")

		// Read the banner file into memory
		bannerFile, err := os.Open(bannerFilename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(bannerFile)
		defer bannerFile.Close()

		var bannerLines []string
		scanner := bufio.NewScanner(bannerFile)
		for scanner.Scan() {
			bannerLines = append(bannerLines, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Generate the ASCII art
		asciiArt := generateASCIIArt(text, bannerLines)
		fmt.Println("asciiArt", asciiArt)
		// Render the response HTML with the ASCII art
		err = tmpl.Execute(w, map[string]interface{}{
			"ASCIIArt": template.HTML(asciiArt),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleAsciiArt(w http.ResponseWriter, r *http.Request) {

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func generateASCIIArt(text string, bannerLines []string) string {
	var asciiArt string

	for i, line := range bannerLines {
		if i < len(text) {
			for j, char := range line {
				if j < len(text) {
					asciiArt += string(text[j])
				} else {
					asciiArt += string(char)
				}
			}
		} else {
			asciiArt += line
		}
		asciiArt += "\n"
	}

	return "asciiArt"
}
