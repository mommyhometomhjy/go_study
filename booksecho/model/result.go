package model

type Result struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Book   `json:"result"`
}
