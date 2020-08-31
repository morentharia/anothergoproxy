package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type ResponseDTO struct {
	*http.Response
	body []byte
}

func NewResponseDTO(r *http.Response) (*ResponseDTO, error) {
	if r == nil {
		return &ResponseDTO{
			Response: &http.Response{},
			body:     []byte(""),
		}, nil
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	err = r.Body.Close()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return &ResponseDTO{
		Response: r,
		body:     body,
	}, nil
}

func (r *ResponseDTO) HttpResponse() *http.Response {
	resp := &http.Response{}
	resp.Request = r.Request
	resp.TransferEncoding = r.TransferEncoding
	resp.Header = make(http.Header)
	resp.Header = r.Header.Clone()
	resp.StatusCode = r.StatusCode
	resp.Status = http.StatusText(r.StatusCode)
	resp.ContentLength = int64(len(r.body))
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(r.body))
	return resp
}

func (resp ResponseDTO) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Status int
		Header http.Header
	}{
		resp.StatusCode,
		resp.Header.Clone(),
	})
}

func (resp *ResponseDTO) UnmarshalJSON(b []byte) error {
	var data struct {
		Status int
		Header http.Header
	}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	resp.Response.StatusCode = data.Status
	resp.Response.Header = data.Header

	return nil
}
