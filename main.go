package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	openai "github.com/sashabaranov/go-openai"
	// ojichat "github.com/greymd/ojichat/generator"
)

func main() {
	// Load config from file or env
	err := loadConfig()
	if err != nil {
		panic(err)
	}
	log.Printf("Finished loading custom configs.")

	config := ClientCfg
	c := openai.NewClientWithConfig(config)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 100,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Who is Hideo Kojima?",
			},
		},
		Stream: true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Printf("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		fmt.Printf(response.Choices[0].Delta.Content)
	}
}
