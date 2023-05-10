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
		ojisanPrompt.WriteString("、加えて指定の要素を含めて")
		ojisanPrompt.WriteString("\n")
		ojisanPrompt.WriteString("要素：")
		ojisanPrompt.WriteString(theme)
	}
	ojisanPrompt.WriteString("\n")
	ojisanPrompt.WriteString("追加資料として、こちらはおじさん構文の例文：")
	ojisanPrompt.WriteString("\n")
	// 追加資料を構築する、三回分なら十分におじさんの方向に行くはず
	for i := 0; i < 3; i++ {
		message, err := generateOjisanMessage(name)
		if err != nil {
			return "", err
		}
		ojisanPrompt.WriteString(message)
		if i < 3-1 {
			ojisanPrompt.WriteString("\n")
		}
	}
	// log.Println(ojisanPrompt.String())

	// by now, the prompt would be like this:

	// テーマが入れる場合：
	// おじさん構文で「（&name）」という人物に挨拶して、加えて以下のテーマを含めて
	// （&theme）
	// 追加資料として、こちらはおじさん構文の例文：
	// なんらかのおじさん発言１
	// なんらかのおじさん発言２
	// なんらかのおじさん発言３ (end)

	// テーマが入れない場合：
	// おじさん構文で「（&name）」という人物に挨拶して
	// 例文として：
	// 追加資料として、こちらはおじさん構文の例文：
	// なんらかのおじさん発言１
	// なんらかのおじさん発言２
	// なんらかのおじさん発言３ (end)

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
