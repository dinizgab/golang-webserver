package main

import (
	"os"

	"github.com/dinizgab/golang-webserver/internal/request"
	"github.com/dinizgab/golang-webserver/internal/response"
	server "github.com/dinizgab/golang-webserver/internal/server"
)

var dirPath string

func main() {
	server := server.New("127.0.0.1", 4221)

	server.AddHandler("GET", "/index", getIndex)
	server.AddHandler("GET", "/file", getFile)

	server.Serve()
}

func getIndex(req *request.Request) *response.Response {
	headers := map[string]string{
		"Content-Type": "text/html",
	}
	response := response.New(200, headers, "<h1>Hello, World!</h1>")

	return response
}

func getFile(req *request.Request) *response.Response {
	headers := map[string]string{
		"Content-Type": "text/html",
	}
	response := response.New(200, headers, "<h1>Hello from getFile!</h1>")

	return response

}

func getFileContent(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileStats, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	fileContent := make([]byte, fileStats.Size())
	n, err := file.Read(fileContent)
	if err != nil {
		return "", err
	}

	return string(fileContent[:n]), nil
}

func postFileContent(path, content string) (string, error) {
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	contentBytes := []byte(content)
	n, err := file.Write(contentBytes)
	if err != nil {
		return "", err
	}

	return string(contentBytes[:n]), nil
}
