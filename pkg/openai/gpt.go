package openai

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"somefun/pkg/cache"
	"somefun/pkg/conf"
	"somefun/pkg/log"
	"sync"
)

type Message = openai.ChatCompletionMessage

type Client interface {
	// chat对话
	Chat(ctx context.Context, message string, token int) (int, string, error)
	// 根据message生成image
	GenerateImage(ctx context.Context, message string) (string, error)
	// 返回聊天历史记录
	History(token int) []Message
}

type client struct {
	//for chatid generation
	sync.Mutex
	*openai.Client
	cache.Cache
	concurrencyLimit chan struct{}
	chatIdGeneration int
}

var gptOnce sync.Once
var openaiClient *client

func GetClient() Client {
	gptOnce.Do(func() {
		conf := conf.GetConfig()
		openaiClient = &client{
			sync.Mutex{},
			openai.NewClient(conf.OpenaiKey),
			cache.GetCache(),
			make(chan struct{}, conf.OpenaiMaxConcurrencyLimit), 1,
		}
	})
	return openaiClient
}

func (c *client) Chat(ctx context.Context, message string, chatId int) (int, string, error) {

	c.concurrencyLimit <- struct{}{}
	defer func() {
		<-c.concurrencyLimit
	}()

	i, exist := c.Load(chatId)
	var contextMessage []Message
	if exist {
		contextMessage = i.([]Message)
	} else {
		c.Lock()
		chatId = c.chatIdGeneration
		c.chatIdGeneration++
		log.Info("chatId is ", chatId)
		c.Unlock()
	}
	contextMessage = append(contextMessage, Message{Role: openai.ChatMessageRoleUser, Content: message})
	log.Info("contextMessage is ", contextMessage)

	resp, err := c.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: contextMessage,
		},
	)
	if err != nil {
		return -1, "", fmt.Errorf("openai client createChatCompletion error:%w", err)
	}

	// store chat record
	c.Store(chatId, append(contextMessage, resp.Choices[0].Message))
	//for debug
	c.View()
	return chatId, resp.Choices[0].Message.Content, nil
}

func (c *client) GenerateImage(ctx context.Context, message string) (string, error) {

	c.concurrencyLimit <- struct{}{}
	defer func() {
		<-c.concurrencyLimit
	}()

	reqUrl := openai.ImageRequest{
		Prompt:         message,
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		N:              1,
	}

	respUrl, err := c.CreateImage(ctx, reqUrl)
	if err != nil {
		return "", fmt.Errorf("Image creation error: %w\n", err)
	}

	return respUrl.Data[0].URL, nil
}

func (c *client) History(token int) []Message {
	return nil
}

//ref:  https://github.com/sashabaranov/go-openai
