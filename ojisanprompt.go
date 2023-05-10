// ojisan.go
package main

import (
	"math/rand"

	"github.com/greymd/ojichat/generator"
)

func generateOjisanMessage(name string) (string, error) {
	randomEmojiNum := rand.Intn(5)         // Random int from 0 to 4
	randomPunctuationLevel := rand.Intn(4) // Random int from 0 to 3

	ojisanConfig := generator.Config{
		TargetName:       name,
		EmojiNum:         randomEmojiNum,
		PunctuationLevel: randomPunctuationLevel,
	}
	return generator.Start(ojisanConfig)
}
