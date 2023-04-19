package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"halyklife/internal/models"
	"io"
	"log"
	"net/http"
)

func HandleRequest(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var req models.Request
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// сохраняем запрос
		reqBody, err := json.Marshal(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Received request: %s", reqBody)

		// создаем новый запрос на данных
		reqHeaders := http.Header{}
		for k, v := range req.Headers {
			reqHeaders.Set(k, v)
		}
		reqToService, err := http.NewRequest(req.Method, req.URL, bytes.NewReader(req.Body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reqToService.Header = reqHeaders

		// отправляем запрос к сервису
		client := &http.Client{}
		respFromService, err := client.Do(reqToService)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer respFromService.Body.Close()

		// читаем ответ от сервиса
		respBody, err := io.ReadAll(respFromService.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// сохраняем ответ от сервиса
		respHeaders := map[string]string{}
		for k, v := range respFromService.Header {
			respHeaders[k] = v[0]
		}
		resp := models.Response{
			ID:      "requestId",
			Status:  respFromService.StatusCode,
			Headers: respHeaders,
			Length:  respFromService.ContentLength,
			Body:    respBody,
		}
		respBody, err = json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Received response: %s", respBody)

		// отправляем ответ
		for k, v := range resp.Headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(resp.Status)
		w.Write(resp.Body)

		// сохраняем запрос и ответ
		reqDoc := models.RequestDocument{
			Request: req,
		}
		reqColl := db.Collection("requests")
		res, err := reqColl.InsertOne(ctx, reqDoc)
		if err != nil {
			log.Printf("Error saving request: %v", err)
		}
		log.Printf("Inserted request with ID %s", res.InsertedID)

		respDoc := models.ResponseDocument{
			Response: resp,
		}
		respColl := db.Collection("responses")
		res, err = respColl.InsertOne(ctx, respDoc)
		if err != nil {
			log.Printf("Error saving response: %v", err)
		}
		log.Printf("Inserted response with ID %s", res.InsertedID)
	}
}
