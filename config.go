package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	openai "github.com/sashabaranov/go-openai"
)

type openaiConfig struct {
	AuthToken          string `json:"authtoken"`
	BaseURL            string `json:"baseurl"`
	OrgID              string `json:"orgid"`
	EmptyMessagesLimit uint   `json:"emptymessageslimit"`
	ProxyURL           string `json:"proxyurl"`
	MaxTokens          uint   `json:"maxtokens"`
}

type gptConfig struct {
	MaxTokens uint `json:"maxtokens"`
}

var openaiCfg openaiConfig
var ClientCfg openai.ClientConfig
var GptCfg gptConfig

func loadConfig() error {
	// Try reading local config file first
	err := loadConfigFromFile()
	if err != nil {
		// If failed to load config from file, try to load from environment variables
		loadConfigFromEnv()
	}
	return nil
}

func loadConfigFromFile() error {
	file, err := os.Open("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("config file not found")
		}
		return err
	}
	defer file.Close()

	// Decode config.json
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&openaiCfg)
	if err != nil {
		return err
	}

	setValue()

	return nil
}

func loadConfigFromEnv() {

	// Read config from env
	openaiCfg.AuthToken = os.Getenv("OPENAI_AUTH_TOKEN")

	openaiCfg.OrgID = os.Getenv("OPENAI_MAX_TOKENS")
	openaiCfg.BaseURL = os.Getenv("OPENAI_BASE_URL")
	openaiCfg.OrgID = os.Getenv("OPENAI_ORG_ID")
	limit, err := strconv.ParseUint(os.Getenv("OPENAI_MTY_MSG_LIM"), 10, 64)
	if err != nil {
		openaiCfg.EmptyMessagesLimit = 300
		log.Printf("Failed to parse OPENAI_MTY_MSG_LIM: %s. Using default value of %d.\n", err, openaiCfg.EmptyMessagesLimit)
	} else {
		openaiCfg.EmptyMessagesLimit = uint(limit)
	}
	openaiCfg.ProxyURL = os.Getenv("OPENAI_PROXY_URL")
	tokens, err := strconv.ParseUint(os.Getenv("OPENAI_MAX_TOKENS"), 10, 64)
	if err != nil {
		GptCfg.MaxTokens = 1000
		log.Printf("Failed to parse OPENAI_MAX_TOKENS: %s. Using default value of %d.\n", err, GptCfg.MaxTokens)
	} else {
		GptCfg.MaxTokens = uint(tokens)
	}
	setValue()
}

func setValue() {
	ClientCfg = openai.DefaultConfig(openaiCfg.AuthToken)
	// Only update BaseURL if it's not empty
	if openaiCfg.BaseURL != "" {
		ClientCfg.BaseURL = openaiCfg.BaseURL
	}

	// Only update OrgID if it's not empty
	if openaiCfg.OrgID != "" {
		ClientCfg.OrgID = openaiCfg.OrgID
	}

	// Only update EmptyMessagesLimit if it's not zero
	if openaiCfg.EmptyMessagesLimit != 0 {
		ClientCfg.EmptyMessagesLimit = openaiCfg.EmptyMessagesLimit
	}

	// Set up the proxy if the proxy URL is provided
	if openaiCfg.ProxyURL != "" {
		proxyUrl, err := url.Parse(openaiCfg.ProxyURL)
		if err != nil {
			log.Printf("Failed to parse proxy URL: %s. Skipping proxy setup.\n", err)
		} else {
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
			ClientCfg.HTTPClient = &http.Client{
				Transport: transport,
			}
		}
	}

	// Only update MaxTokens if it's not zero
	if openaiCfg.MaxTokens > 0 {
		GptCfg.MaxTokens = openaiCfg.MaxTokens
	} else {
		// Get recommendation value if setting is invalid
		GptCfg.MaxTokens = 1000
		log.Printf("Invalid MaxTokens setting, using default value of: %d.\n", GptCfg.MaxTokens)
	}

}
