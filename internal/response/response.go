package response

import "fmt"

type Response struct {
    StatusCode int
    StatusMessage string
    Headers map[string]string
    Body string
}

var statusMessages = map[int]string{
    200: "OK",
    400: "Bad Request",
    404: "Not Found",
    500: "Internal Server Error",
}

func New(statusCode int, headers map[string]string, body string) *Response {
    return &Response{
        StatusCode: statusCode,
        StatusMessage: statusMessages[statusCode],
        Headers: headers,
        Body: body,
    }
}

func (r *Response) String() string {
    response := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.StatusCode, r.StatusMessage)
    for key, value := range r.Headers {
        response += fmt.Sprintf("%s: %s\r\n", key, value)
    }
    response += fmt.Sprintf("\r\n%s", r.Body)

    return response
}
