package main

type ApiResult struct {
	Code    string                   `json:"code"`
	Message string                   `json:"message"`
	Body    []map[string]interface{} `json:"body"`
}
