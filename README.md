# ojichat-aigc
Ojichat powered by OpenAI, now ojisan has become knowledgeable!

## 何が起こったんですか？
最先端のLLMサービスを使用して、今のおじさんはいろんな分野の話題ができます！（いったいなんのために…）

## Ojichatとは
greymdさんが創立した「おじさん構文」を生成するプログラムです。元ページは： https://github.com/greymd/ojichat 

## 使い方
1. 本コードをダウンロードして
2. config-template.json ファイルの名前を config.json にリネームして
3. config.json を開いて、特に要件がなければ、"authtoken" の行だけに OpenAI API keyを入力し、他のはそのままで
4. `go run .` コマンドを実行して
5. 話相手の名前を入力して、入力キーを押して確認
6. 次は話題にしたいものを入力し、複数のものを区切る場合はスペースを使用して、入力キーを押して確認
7. prompt（AIへの指令）一度チェックして、問題ないなら入力キーて進む
8. その後、結果が表示されます

## 注意
1. このソフトを利用するには、OpenAI ID が必要です。
2. このソフト自体は無料ですが、OpenAI の GPT モデルによって生成されるコンテンツは　OpenAI によって課金されます。
3. 元のコードには、==**セクハラ**==に関する内容が僅かに含まれることがあり、==**OpenAI アカウントが停止されるの可能性が少なくとも存在します**==。どうかご注意ください。

## Config.json
- maxtokens: 簡単に言うとGPTの処理量というもの、ユーザーの輸入とモデルの max_tokens はユーザーの入力(prompt) およびモデルの出力(completion) の文字の総量を制限します。
        日本語の仮名や漢字のtokensは1.2〜3程度です。
        OpenAIはモデルのタイプと毎回処理したtokensの量で料金を計算する。現時点で（令和5年5月10日）、instruct能力を備えた最も低価格なモデルはGPT-3.5-turboです、1000 tokens（約700文字相当）で0.002 ドールになっています。
        なお、GPT-3.5-turbo の最大の処理量は2048となります。
        高く設定されているからといって、必ず制限までの長さの回答が生成されるわけではありません。ただし、低く設定すると、モデルの回答が途中で切断される可能性があります。
        本プログラムの prompt + completion は通常600〜700程度です。
        
- authtoken：これはOpenAIがユーザーを識別するために使用する資格情報であり、OpenAIアカウントを登録した後、https://platform.openai.com/account/api-keys にアクセスして作成することができます。
        これは個人情報であり、金にかかるものです。他人に漏らさないようにしてください。
        本プログラムは一切の情報の収集機能がついていません。
- baseurl：ここで設定されたアドレスは、OpenAIのAPIアドレスと置き換えられます、一般的には不要です。
- orgid：OpenAIでの組織認証コードです、一般的には不要です。
- emptymessageslimit：おじさんにもわかんない！
- proxyurl: プロキシの設定です、インターネットに直接に接続された実行環境では不要です。
- temperature/topp/frequencypenalty/presencepenalty: モデル関連の設定です、おじさんが面白いかおかしいかに関わるです、詳細は https://platform.openai.com/docs/api-reference/completions/create まで
