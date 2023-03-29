package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"somefun/api/gpt"
	"somefun/pkg/log"
	"somefun/pkg/openai"
)

type respBody struct {
	Err  string      `json:"err,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func returnResp(data interface{}, err error, w http.ResponseWriter) {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	resp := &respBody{
		errStr,
		data,
	}
	bs, _ := json.Marshal(resp)
	w.Write(bs)
}

func ChatHandle(w http.ResponseWriter, r *http.Request) {
	var chatRequest *gpt.ChatRequest
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("io read error:", err)
		returnResp(nil, err, w)
		return
	}
	err = json.Unmarshal(bs, &chatRequest)
	if err != nil {
		log.Error("json unmarshal error:", err)
		returnResp(nil, err, w)
		return
	}
	chatId, message, err := openai.GetClient().Chat(context.TODO(), chatRequest.Message, int(chatRequest.ChatId))
	if err != nil {
		log.Error("Chat error: %w", err)
		returnResp(nil, err, w)
		return
	}
	returnResp(gpt.ChatResponse{ChatId: int64(chatId), Message: message}, nil, w)
	return
}

func GenerateImage(w http.ResponseWriter, r *http.Request) {
	var imageReq *gpt.ImageRequest
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("io read error:", err)
		returnResp(nil, err, w)
		return
	}
	err = json.Unmarshal(bs, &imageReq)
	if err != nil {
		log.Error("json unmarshal error:", err)
		returnResp(nil, err, w)
		return
	}

	imageUrl, err := openai.GetClient().GenerateImage(context.Background(), imageReq.Message)
	if err != nil {
		log.Error("GenerateImage error: %w", err)
		returnResp(nil, err, w)
		return
	}
	returnResp(imageUrl, nil, w)
	return
}
