package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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

const buffSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {

	r := &Request{
		parseStatus: Initialised,
	}
	buf := make([]byte, buffSize)
	readToIndex := 0

	for r.parseStatus != Done {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		bytesRead, err := reader.Read(buf[readToIndex:])

		if err != nil {

			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}

		readToIndex += bytesRead

		bytesParsed, err := r.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[bytesParsed:])
		readToIndex -= bytesParsed

	}

	return r, nil
}

func (r *Request) parse(data []byte) (int, error) {

	if r.parseStatus == Initialised {
		reqLine, linesRead, err := parseRequestLine(data)

		if err != nil {
			return -1, err
		}

		if linesRead == 0 {
			return 0, nil
		}

		r.RequestLine = reqLine
		r.parseStatus = Done
		return linesRead, nil
	}

	if r.parseStatus == Done {
		return -1, errors.New("error: trying to read data in a done state")
	}

	return -1, errors.New("error: unknown state")
}

func parseRequestLine(b []byte) (RequestLine, int, error) {
	rawRqLine := bytes.Split(b, []byte("\r\n"))
	if len(rawRqLine) == 1 {
		return RequestLine{}, 0, nil
	}

	s := string(rawRqLine[0])
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
