# back end

## Requirements
docker
make
bash

## ビルド
```
$ make up
```

## コンテナの停止
```
$ make down
```

## OAuth
1. `GET /oauth/generate/code`
2. `https://q.trap.jp/api/v3/oauth2/authorize?code_challenge={{codeChallenge}}&code_challenge_method={{codeChallengeMethod}}&client_id={{clientID}}&response_type={{responseType}}`
    →リダイレクト先URLの`code`をコピー
3. `GET /oauth/callback?code={{code}}`

## フォーマット
```
$ make lint
```
