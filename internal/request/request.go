package request

import (
	"bufio"
	"errors"
	"strings"
)

type Request struct {
    Path string
    Method string
    Headers map[string]string
    Body string
}

func New(req string) (*Request, error) {
    request := &Request{}

    splittedRequest := strings.Split(req, "\r\n\r\n")
    headers, body := splittedRequest[0], splittedRequest[1]
    err := request.parseHeaders(headers)
    if err != nil {
        return nil, err
    }
    
    request.Body = body

    return request, nil
}

func (req *Request) parseHeaders(headers string) error {
    req.Headers = map[string]string{}

    scanner := bufio.NewScanner(strings.NewReader(headers))
    scanner.Split(bufio.ScanLines)

    scanner.Scan()
    startLine := scanner.Text()
    splittedStartLine := strings.Split(startLine, " ")
    req.Method = splittedStartLine[0]
    req.Path = splittedStartLine[1]

    for scanner.Scan() {
        line := scanner.Text()

        if line == "" {
            continue
        }

        parts := strings.Split(line, ": ")
        if len(parts) < 2 {
            return errors.New("Invalid header line")
        }

        req.Headers[parts[0]] = parts[1] 
    }

    return nil
}
