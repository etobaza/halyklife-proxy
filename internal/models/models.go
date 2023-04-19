package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body,omitempty"`
}

type Response struct {
	ID      string            `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Length  int64             `json:"length"`
	Body    []byte            `json:"body,omitempty"`
}

type RequestDocument struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Request Request            `bson:"request"`
}

type ResponseDocument struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Response Response           `bson:"response"`
}
