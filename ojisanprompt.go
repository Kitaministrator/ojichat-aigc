// ojisan.go
package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/greymd/ojichat/generator"
)

func buildOjisanPrompt(name string, theme string) (string, error) {
	// 簡単な、明確なpromptを構築する
	var ojisanPrompt strings.Builder
	fmt.Fprintf(&ojisanPrompt, "おじさん構文で「%s」という人物に挨拶して、相手の名前と愛称を付けて", name)
	if len(strings.Fields(theme)) > 0 {
		ojisanPrompt.WriteString("、加えて指定の要素も含めて。")
		ojisanPrompt.WriteString("\n")
		ojisanPrompt.WriteString("要素：")
		ojisanPrompt.WriteString(theme)
		ojisanPrompt.WriteString("\n")
	} else {
		ojisanPrompt.WriteString("。\n")
	}
	ojisanPrompt.WriteString("追加資料として、こちらはおじさん構文の3つの例です：")
	ojisanPrompt.WriteString("\n")
	// モデルが特徴を捉えるために、追加の情報を収集してデータを構築する必要があります、三回分なら十分におじさんの方向に行くはず
	for i := 0; i < 3; i++ {
		message, err := generateOjisanMessage(name)
		if err != nil {
			return "", err
		}
		ojisanPrompt.WriteString(fmt.Sprintf("%d. ", i+1))
		ojisanPrompt.WriteString(message)
		ojisanPrompt.WriteString("\n")
	}
	ojisanPrompt.WriteString("出力：")
	return ojisanPrompt.String(), nil
}

func generateOjisanMessage(name string) (string, error) {

	// 絵文字/顔文字の最大連続数と、句読点の使用頻度をランダムに決める
	randomEmojiNum := rand.Intn(10)        // Random int from 0 to 9
	randomPunctuationLevel := rand.Intn(4) // Random int from 0 to 3

	// なまえ、絵文字/顔文字の最大連続数、句読点の使用頻度を設定して、おじさん構文を生成する
	ojisanConfig := generator.Config{
		TargetName:       name,
		EmojiNum:         randomEmojiNum,
		PunctuationLevel: randomPunctuationLevel,
	}

	return generator.Start(ojisanConfig)
}
