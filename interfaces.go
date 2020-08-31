package main

type ReqRespCacheI interface {
	Load(*RequestDTO) (*ResponseDTO, error)
	Store(*RequestDTO, *ResponseDTO) error
}
