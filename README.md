# back end

## Requirements

docker
make

## ビルド

```bash
$ make up
```

## ビルド(バックグラウンド)

```bash
$ make up-d
```

## docker のログ確認

```bash
$ make logs
```

## コンテナの停止

```bash
$ make down
```

## フォーマット

```
$ make lint
```

## OAuth

1. `GET /oauth/generate/code`
2. `https://q.trap.jp/api/v3/oauth2/authorize?code_challenge={{codeChallenge}}&code_challenge_method={{codeChallengeMethod}}&client_id={{clientID}}&response_type={{responseType}}`
   → リダイレクト先 URL の`code`をコピー
3. `GET /oauth/callback?code={{code}}`
