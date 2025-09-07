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
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalln(err)
	}

	str := string(b)

	rawRqLine := strings.Split(str, "\r\n")
	reqLine, err := parseRequestLine(rawRqLine[0])

	if err != nil {
		return nil, err
	}

	out := Request{reqLine}

	return &out, nil
}

func parseRequestLine(s string) (RequestLine, error) {
	parts := strings.Split(s, " ")

	fmt.Println(len(parts))
	if len(parts) != 3 {
		return RequestLine{}, errors.New("Request-line does not contain the three specified parts")
	}

	// Verify method
	method := parts[0]
	if !isAlphabeticAndUppercase(method) {
		return RequestLine{}, errors.New("Method is not valid")

	}

	target := parts[1]

	// Verify version
	versionFull := parts[2]
	versionParts := strings.Split(versionFull, "/")
	if versionParts[0] != "HTTP" || versionParts[1] != "1.1" {
		return RequestLine{}, errors.New("Version of HTTP not supported. Only HTTP/1.1 is supported.")
	}

	return RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   versionParts[1],
	}, nil

}

func isAlphabeticAndUppercase(s string) bool {
	return !strings.ContainsFunc(s,
		func(r rune) bool {
			return !unicode.IsLetter(r) || !unicode.IsUpper(r)
		})
}
