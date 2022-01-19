package handler

import "time"

const (
	// OK 請求成功，適用Read、Update
	// Created 請求成功，適用Create
	// Accepted 此請求已被接受但未做任何處理
	// NoContent server已經處理請求，但未返回任何內容，適用Delete
	// BadRequest server無法理解請求，適用錯誤api參數
	// Unauthorized 未經過token認證
	// Forbidden 無權限訪問，這邊跟Unauthorized不同的是，這邊有token但無權限
	// NotFound	找不到資源
	OK           = 200
	Created      = 201
	Accepted     = 202
	NoContent    = 204
	BadRequest   = 400
	Unauthorized = 401
	Forbidden    = 403
	NotFound     = 404
)

var ResponseFlag = map[int]string{
	OK:           "請求成功",
	Created:      "請求成功",
	Accepted:     "請求成功",
	NoContent:    "請求成功",
	BadRequest:   "請依照API文件重新發起請求",
	Unauthorized: "該請求未經過認證",
	Forbidden:    "你的權限不足以發起該請求",
	NotFound:     "",
}

type Response struct {
	ResultCode int       `json:"result_code"`
	Message    string    `json:"message"`
	Data       string    `json:"data"`
	TimeStamp  time.Time `json:"time_stamp"`
}
