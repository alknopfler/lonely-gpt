package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func convertAudioToText(audioFilePath, apiKey string) (string, error) {
	log.Println("Converting audio to text")
	file, err := os.Open(audioFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open audio file: %s", err.Error())
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("audio", audioFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %s", err.Error())
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file data: %s", err.Error())
	}
	writer.Close()

	url := "https://api.openai.com/v1/engines/davinci-codex/completions"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %s", err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %s", err.Error())
	}
	defer resp.Body.Close()

	var respBody bytes.Buffer
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %s", err.Error())
	}

	response := respBody.String()
	startIndex := bytes.Index([]byte(response), []byte(`text":"`)) + len(`text":"`)
	endIndex := bytes.Index([]byte(response[startIndex:]), []byte(`"`))
	text := response[startIndex : startIndex+endIndex]

	return string(text), nil
}
