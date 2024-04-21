package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/dinizgab/golang-webserver/internal/request"
	"github.com/dinizgab/golang-webserver/internal/response"
	server "github.com/dinizgab/golang-webserver/internal/server"
)

var dirPath string

func main() {
    server := server.New("127.0.0.1", 4221)

    server.AddHandler("/index", getIndex)

    server.Serve()
}

func getIndex(req *request.Request) *response.Response {
    response := response.New(200, map[string]string{}, "Hello, from index!")
    
    return response
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		os.Exit(1)
	}
	request := string(buffer[:n])

	headersMap := map[string]string{}
	headers, body := strings.Split(request, "\r\n\r\n")[0], strings.Split(request, "\r\n\r\n")[1]

	headerLines := strings.Split(headers, "\r\n")
	for i := 1; i < len(headerLines); i++ {
		if headerLines[i] == "" {
			break
		}
		headerParts := strings.Split(headerLines[i], ": ")
		headerName, headerValue := headerParts[0], headerParts[1]

		headersMap[headerName] = headerValue
	}

	splitHeader := strings.Split(headerLines[0], " ")
	method, path := splitHeader[0], splitHeader[1]

	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.HasPrefix(path, "/echo") {
		result := strings.Split(path, "/echo/")[1]

		response := "HTTP/1.1 200 OK\r\n"
		response += "Content-Type: text/plain\r\n"
		response += fmt.Sprintf("Content-Length: %d\r\n\r\n", len(result))
		response += fmt.Sprintf("%s\r\n", result)

		conn.Write([]byte(response))
	} else if strings.HasPrefix(path, "/user-agent") {
		userAgent := headersMap["User-Agent"]

		response := "HTTP/1.1 200 OK\r\n"
		response += "Content-Type: text/plain\r\n"
		response += fmt.Sprintf("Content-Length: %d\r\n\r\n", len(userAgent))
		response += fmt.Sprintf("%s\r\n", userAgent)

		conn.Write([]byte(response))
	} else if strings.HasPrefix(path, "/files") && method == "POST" {
		filename := strings.Split(path, "/files/")[1]
		fullFilePath := dirPath + filename

		result, err := postFileContent(fullFilePath, body)
		if err != nil {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}

		response := "HTTP/1.1 201 OK\r\n"
		response += "Content-Type: application/octet-stream\r\n"
		response += fmt.Sprintf("Content-Length: %d\r\n\r\n", len(result))
		response += fmt.Sprintf("%s\r\n", result)

		conn.Write([]byte(response))
	} else if strings.HasPrefix(path, "/files") {
		filename := strings.Split(path, "/files/")[1]
		fullFilePath := dirPath + filename

		fileContent, err := getFileContent(fullFilePath)
		if err != nil {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}

		response := "HTTP/1.1 200 OK\r\n"
		response += "Content-Type: application/octet-stream\r\n"
		response += fmt.Sprintf("Content-Length: %d\r\n\r\n", len(fileContent))
		response += fmt.Sprintf("%s\r\n", fileContent)

		conn.Write([]byte(response))
	} else if path != "/" {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
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
