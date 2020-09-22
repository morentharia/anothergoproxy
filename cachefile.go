package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
)

type CacheFile struct{}

var _ ReqRespCacheI = &CacheFile{}

func NewCacheFile() (*CacheFile, error) {
	return &CacheFile{}, nil
}

func (c *CacheFile) Load(req *RequestDTO) (*ResponseDTO, error) {
	hash := req.Hash()
	ResponseInfoFilename := filepath.Join(options.CachePath(), fmt.Sprintf("%s_resp.json", hash))
	ResponseBodyFilename := filepath.Join(options.CachePath(), fmt.Sprintf("%s_resp_body", hash))

	data, err := ioutil.ReadFile(ResponseInfoFilename)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r, err := NewResponseDTO(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = json.Unmarshal(data, r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r.body, err = ioutil.ReadFile(ResponseBodyFilename)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// logrus.Info("dlkjfdlkfj")

	return r, nil
}

func (c *CacheFile) Store(req *RequestDTO, resp *ResponseDTO) error {
	hash := req.Hash()
	RequestInfoFilename := filepath.Join(options.CachePath(), fmt.Sprintf("%s_req.json", hash))
	RequestBodyFilename := filepath.Join(options.CachePath(), fmt.Sprintf("%s_req_body", hash))
	ResponseInfoFilename := filepath.Join(options.CachePath(), fmt.Sprintf("%s_resp.json", hash))
	ResponseBodyFilename := filepath.Join(options.CachePath(), fmt.Sprintf("%s_resp_body", hash))

	b, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return errors.WithStack(err)
	}

	err = ioutil.WriteFile(RequestInfoFilename, b, 0644)
	if err != nil {
		return errors.WithStack(err)
	}

	err = ioutil.WriteFile(RequestBodyFilename, req.body, 0644)
	if err != nil {
		return errors.WithStack(err)
	}

	b, err = json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return errors.WithStack(err)
	}

	err = ioutil.WriteFile(ResponseInfoFilename, b, 0644)
	if err != nil {
		return errors.WithStack(err)
	}

	err = ioutil.WriteFile(ResponseBodyFilename, resp.body, 0644)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
