package server

import "time"

type ServiceResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Count   int         `json:"count"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
}

type addInput struct {
	Mac  string    `json:"mac"`
	Time time.Time `json:"time"`
}

type listInput struct {
	Offset int `json:"offset"`
}
