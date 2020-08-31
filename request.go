package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/sirupsen/logrus"
)

type RequestDTO struct {
	*http.Request
	body []byte
}

func NewRequestDTO(req *http.Request) *RequestDTO {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logrus.WithError(err).Error("can't read body")
		return &RequestDTO{req, []byte("")}
	}
	req.Body.Close()
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return &RequestDTO{req, body}
}

func (req RequestDTO) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Method     string
		Host       string
		RequestURI string
		URL        *url.URL
		Header     http.Header
	}{
		req.Method,
		req.Host,
		req.RequestURI,
		req.URL,
		req.Header.Clone(),
	})
}

func (req RequestDTO) RawString() string {
	dump, err := httputil.DumpRequest(req.Request, true)
	if err != nil {
		logrus.WithError(err).Warn("DumpRequest")
		return ""
	}
	return string(dump)
}

func (req RequestDTO) Hash() string {
	data := fmt.Sprintf(
		"%s %s %s",
		req.Method, req.URL.String(), string(req.body),
	)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))[:10]
}
