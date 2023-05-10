package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
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

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("なまえ：")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		name := strings.TrimSpace(input)
		if len(name) > 40 {
			fmt.Println("なまえが長すぎるよ、、、") // なまえが長すぎとおじさんぽっくない
			continue
		}

		fmt.Print("テーマをいれてね (スペースで区切って複数の単語を入力できます)：")
		themeInput, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		theme := strings.TrimSpace(themeInput)

		prompt, err := buildOjisanPrompt(name, theme)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		fmt.Println("This prompt will be sent to OpenAI:")
		fmt.Println(prompt)
		fmt.Println("Press enter to continue...")
		reader.ReadString('\n')

		req := openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			MaxTokens: int(GptCfg.MaxTokens),
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature:      GptCfg.Temperature,
			TopP:             GptCfg.TopP,
			FrequencyPenalty: GptCfg.FrequencyPenalty,
			PresencePenalty:  GptCfg.PresencePenalty,
			Stream:           true,
		}
		stream, err := c.CreateChatCompletionStream(ctx, req)
		if err != nil {
			fmt.Printf("ChatCompletionStream error: %v\n", err)
			return
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}

			if err != nil {
				fmt.Printf("\nStream error: %v\n", err)
				return
			}

			fmt.Printf(response.Choices[0].Delta.Content)
		}
	}

}
