package main

type ApiResult struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Body    ApiBody `json:"body"`
}

type ApiBody struct {
	RowCount int         `json:"rowCount"`
	Result   interface{} `json:"result"`
}

type Table struct {
	TableName string `json:"tableName"`
}

type TableDesc struct {
	Field   string  `json:"field"`
	Type    string  `json:"type"`
	Null    string  `json:"null"`
	Key     *string `json:"key"`
	Default *string `json:"default"`
	Extra   *string `json:"extra"`
}
