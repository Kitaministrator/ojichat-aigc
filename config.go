package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"

	openai "github.com/sashabaranov/go-openai"
)

type openaiConfig struct {
	AuthToken          string `json:"authtoken"`
	BaseURL            string `json:"baseurl"`
	OrgID              string `json:"orgid"`
	EmptyMessagesLimit uint   `json:"emptymessageslimit"`
}

var openaiCfg openaiConfig
var ClientCfg openai.ClientConfig

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
	openaiCfg.BaseURL = os.Getenv("OPENAI_BASE_URL")
	openaiCfg.OrgID = os.Getenv("OPENAI_ORG_ID")
	limit, err := strconv.ParseUint(os.Getenv("OPENAI_MTY_MSG_LIM"), 10, 64)
	if err != nil {
		openaiCfg.EmptyMessagesLimit = 300
		log.Printf("Failed to parse OPENAI_MTY_MSG_LIM: %s. Using default value of %d.\n", err, openaiCfg.EmptyMessagesLimit)
	} else {
		openaiCfg.EmptyMessagesLimit = uint(limit)
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
}

/// external file config.json
/// path: ./config.json
/// structure:
/*
	{
		"authtoken": "your-token",
		"baseurl": "your-private-domain-api",
		"orgid": "your-organization-id",
		"emptymessageslimit": 300
	}
*/

/// list of env variables
/*
	OPENAI_AUTH_TOKEN
	OPENAI_BASE_URL
	OPENAI_ORG_ID
	OPENAI_MTY_MSG_LIM
*/
