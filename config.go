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
	AuthToken          string  `json:"authtoken"`
	BaseURL            string  `json:"baseurl"`
	OrgID              string  `json:"orgid"`
	EmptyMessagesLimit uint    `json:"emptymessageslimit"`
	ProxyURL           string  `json:"proxyurl"`
	MaxTokens          uint    `json:"maxtokens"`
	Temperature        float32 `json:"temperature"`
	TopP               float32 `json:"topp"`
	FrequencyPenalty   float32 `json:"frequencypenalty"`
	PresencePenalty    float32 `json:"presencepenalty"`
}

type gptConfig struct {
	MaxTokens        uint    `json:"maxtokens"`
	Temperature      float32 `json:"temperature"`
	TopP             float32 `json:"topp"`
	FrequencyPenalty float32 `json:"frequencypenalty"`
	PresencePenalty  float32 `json:"presencepenalty"`
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
	temp, err := strconv.ParseFloat(os.Getenv("OPENAI_TEMPERATURE"), 32)
	if err != nil {
		GptCfg.Temperature = 1.0
		log.Printf("Failed to parse OPENAI_TEMPERATURE: %s. Using default value of %.1f.\n", err, GptCfg.Temperature)
	} else {
		GptCfg.Temperature = float32(temp)
	}

	topP, err := strconv.ParseFloat(os.Getenv("OPENAI_TOP_P"), 32)
	if err != nil {
		GptCfg.TopP = 1.0
		log.Printf("Failed to parse OPENAI_TOP_P: %s. Using default value of %.1f.\n", err, GptCfg.TopP)
	} else {
		GptCfg.TopP = float32(topP)
	}

	frequencyPenalty, err := strconv.ParseFloat(os.Getenv("OPENAI_FREQUENCY_PENALTY"), 32)
	if err != nil {
		GptCfg.FrequencyPenalty = 0.0
		log.Printf("Failed to parse OPENAI_FREQUENCY_PENALTY: %s. Using default value of %.1f.\n", err, GptCfg.FrequencyPenalty)
	} else {
		GptCfg.FrequencyPenalty = float32(frequencyPenalty)
	}

	presencePenalty, err := strconv.ParseFloat(os.Getenv("OPENAI_PRESENCE_PENALTY"), 32)
	if err != nil {
		GptCfg.PresencePenalty = 0.0
		log.Printf("Failed to parse OPENAI_PRESENCE_PENALTY: %s. Using default value of %.1f.\n", err, GptCfg.PresencePenalty)
	} else {
		GptCfg.PresencePenalty = float32(presencePenalty)
	}
	setValue()
}

func setValue() {
	// Set client parameters
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

	// Set GPT parameters
	if openaiCfg.MaxTokens > 0 && openaiCfg.MaxTokens < 2049 {
		GptCfg.MaxTokens = openaiCfg.MaxTokens
	} else {
		// Get recommendation max tokens value if setting is invalid
		GptCfg.MaxTokens = 1000
		log.Printf("Invalid MaxTokens setting, using default value of: %d.\n", GptCfg.MaxTokens)
	}

	if openaiCfg.Temperature > 0 && openaiCfg.Temperature < 1.0 {
		GptCfg.Temperature = openaiCfg.Temperature
	} else {
		GptCfg.Temperature = 1.0
		log.Printf("Invalid Temperature setting, using default value of: %.1f.\n", GptCfg.Temperature)
	}

	if openaiCfg.TopP > 0 && openaiCfg.TopP < 1.0 {
		GptCfg.TopP = openaiCfg.TopP
	} else {
		GptCfg.TopP = 1.0
		log.Printf("Invalid TopP setting, using default value of: %.1f.\n", GptCfg.TopP)
	}

	if openaiCfg.FrequencyPenalty > 0 && openaiCfg.FrequencyPenalty < 2.0 {
		GptCfg.FrequencyPenalty = openaiCfg.FrequencyPenalty
	} else {
		GptCfg.FrequencyPenalty = 0.0
		log.Printf("Invalid FrequencyPenalty setting, using default value of: %.1f.\n", GptCfg.FrequencyPenalty)
	}

	if openaiCfg.PresencePenalty > 0 && openaiCfg.PresencePenalty < 2.0 {
		GptCfg.PresencePenalty = openaiCfg.PresencePenalty
	} else {
		GptCfg.PresencePenalty = 0.0
		log.Printf("Invalid PresencePenalty setting, using default value of: %.1f.\n", GptCfg.PresencePenalty)
	}
}
