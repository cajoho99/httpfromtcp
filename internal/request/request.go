package request

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	parseStatus int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	Initialised int = iota
	Done
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalln(err)
	}

	str := string(b)

	reqLine, numBytesRead, err := parseRequestLine(str)

	if err != nil {
		return nil, err
	}

	out := Request{reqLine}

	return &out, nil
}

func (r *Request) parse(data []byte) (int, error) {

}

func parseRequestLine(str string) (RequestLine, int, error) {
	rawRqLine := strings.Split(str, "\r\n")
	if len(rawRqLine) == 1 {
		return RequestLine{}, 0, nil
	}

	s := rawRqLine[0]
	numBytesRead := len(s)
	parts := strings.Split(s, " ")

	fmt.Println(len(parts))
	if len(parts) != 3 {
		return RequestLine{}, numBytesRead, errors.New("Request-line does not contain the three specified parts")
	}

	// Verify method
	method := parts[0]
	if !isAlphabeticAndUppercase(method) {
		return RequestLine{}, numBytesRead, errors.New("Method is not valid")

	}

	target := parts[1]

	// Verify version
	versionFull := parts[2]
	versionParts := strings.Split(versionFull, "/")
	if versionParts[0] != "HTTP" || versionParts[1] != "1.1" {
		return RequestLine{}, numBytesRead, errors.New("Version of HTTP not supported. Only HTTP/1.1 is supported.")
	}

	return RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   versionParts[1],
	}, numBytesRead, nil

}

func isAlphabeticAndUppercase(s string) bool {
	return !strings.ContainsFunc(s,
		func(r rune) bool {
			return !unicode.IsLetter(r) || !unicode.IsUpper(r)
		})
}
