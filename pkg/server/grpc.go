package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"somefun/api/gpt"
	"somefun/pkg/log"
	"somefun/pkg/openai"
)

func NewGptHandler() (gpt.GptServiceServer, error) {
	return &gptService{}, nil
}

type gptService struct{}

func (gs *gptService) Chat(chatServer gpt.GptService_ChatServer) error {
	for {
		chatRequest, err := chatServer.Recv()
		if err != nil {
			if err == io.EOF {
				log.Infof("stream[recv] exit")
				return nil
			}
			log.Errorf("chatServer Recv error:", err)
			return err
		}
		log.Info("chatRequest is:", chatRequest)

		chatId, message, err := openai.GetClient().Chat(context.TODO(), chatRequest.Message, int(chatRequest.ChatId))

		err = chatServer.Send(&gpt.ChatResponse{
			ChatId:  int64(chatId),
			Message: message,
		})
		if err == io.EOF {
			log.Infof("stream[send] exit")
			return nil
		} else if err != nil {
			log.Errorf("chatServer Send error:", err)
			return err
		} else {

		}
	}
	return status.Errorf(codes.Unimplemented, "method Chat not implemented")
}

func (gs *gptService) GenerateImage(ctx context.Context, imageReq *gpt.ImageRequest) (*gpt.ImageResponse, error) {
	imageUrl, err := openai.GetClient().GenerateImage(ctx, imageReq.Message)
	if err != nil {
		return nil, fmt.Errorf("GenerateImage error: %w", err)
	}
	return &gpt.ImageResponse{Message: imageUrl}, nil
}

//type GptServiceServer interface {
//	Chat(GptService_ChatServer) error
//	GenerateImage(context.Context, *ImageRequest) (*ImageResponse, error)
//}
