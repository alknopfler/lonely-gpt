package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	log.Println("Serving home page")
	tmpl.Execute(w, nil)
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Converting audio to text")
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}

	audioFile, _, err := r.FormFile("audio")
	if err != nil {
		log.Println("Failed to retrieve audio file:", err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
	defer audioFile.Close()

	audioData, err := ioutil.ReadAll(audioFile)
	if err != nil {
		log.Println("Failed to read audio file data:", err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	audioFilePath := "audio.wav"
	err = ioutil.WriteFile(audioFilePath, audioData, 0644)
	if err != nil {
		log.Println("Failed to save audio file:", err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	text, err := convertAudioToText(audioFilePath, os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Println("Failed to convert audio to text:", err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Texto convertido: %s", text)

	// Remove the temporary audio file
	os.Remove(audioFilePath)
}
