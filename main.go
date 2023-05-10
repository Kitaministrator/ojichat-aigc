package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	// openai "github.com/sashabaranov/go-openai"
)

func main() {
	// leave a name here for manually test...
	// しおりこちゃん
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("なまえをいれてね：")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		name := strings.TrimSpace(input)
		if len(name) > 40 {
			fmt.Println("なまえが長すぎるよ！")
			continue
		}

		result, err := generateOjisanMessage(name)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", result)
	}

	// // Load config from file or env
	// err := loadConfig()
	// if err != nil {
	// 	panic(err)
	// }
	// log.Printf("Finished loading custom configs.")

	// config := ClientCfg
	// c := openai.NewClientWithConfig(config)
	// ctx := context.Background()

	// req := openai.ChatCompletionRequest{
	// 	Model:     openai.GPT3Dot5Turbo,
	// 	MaxTokens: 100,
	// 	Messages: []openai.ChatCompletionMessage{
	// 		{
	// 			Role:    openai.ChatMessageRoleUser,
	// 			Content: "Who is Hideo Kojima?",
	// 		},
	// 	},
	// 	Stream: true,
	// }
	// stream, err := c.CreateChatCompletionStream(ctx, req)
	// if err != nil {
	// 	fmt.Printf("ChatCompletionStream error: %v\n", err)
	// 	return
	// }
	// defer stream.Close()

	// fmt.Printf("Stream response: ")
	// for {
	// 	response, err := stream.Recv()
	// 	if errors.Is(err, io.EOF) {
	// 		fmt.Println("\nStream finished")
	// 		return
	// 	}

	// 	if err != nil {
	// 		fmt.Printf("\nStream error: %v\n", err)
	// 		return
	// 	}

	// 	fmt.Printf(response.Choices[0].Delta.Content)
	// }
}

func buildPrompt(name string) (string, error) {
	var result strings.Builder
	for i := 0; i < 3; i++ {
		message, err := generateOjisanMessage(name)
		if err != nil {
			return "", err
		}
		result.WriteString(message)
		if i < 3-1 {
			result.WriteString("\n")
		}
	}
	return result.String(), nil
}
