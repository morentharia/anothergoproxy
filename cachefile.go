package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type CacheFile struct {
	Request              *RequestDTO
	Response             *ResponseDTO
	RequestInfoFilename  string
	RequestBodyFilename  string
	ResponseInfoFilename string
	ResponseBodyFilename string
}

func NewCacheFile(req *RequestDTO) *CacheFile {
	hash := req.Hash()
	dataPath := filepath.Join(options.OutputPath, "data")
	return &CacheFile{
		Request:              req,
		RequestInfoFilename:  filepath.Join(dataPath, fmt.Sprintf("%s_req.json", hash)),
		RequestBodyFilename:  filepath.Join(dataPath, fmt.Sprintf("%s_req_body", hash)),
		ResponseInfoFilename: filepath.Join(dataPath, fmt.Sprintf("%s_resp.json", hash)),
		ResponseBodyFilename: filepath.Join(dataPath, fmt.Sprintf("%s_resp_body", hash)),
	}
}

func (c *CacheFile) Save(resp *ResponseDTO) error {
	b, err := json.MarshalIndent(c.Request, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.RequestInfoFilename, b, 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.RequestBodyFilename, c.Request.body, 0644)
	if err != nil {
		return err
	}

	b, err = json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.ResponseInfoFilename, b, 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.ResponseBodyFilename, resp.body, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheFile) Load() (*ResponseDTO, error) {
	data, err := ioutil.ReadFile(c.ResponseInfoFilename)
	if err != nil {
		return nil, err
	}

	r, err := NewResponseDTO(nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}

	r.body, err = ioutil.ReadFile(c.ResponseBodyFilename)
	if err != nil {
		return nil, err
	}

	return r, nil
}
