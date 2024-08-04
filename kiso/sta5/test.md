# Basic認証の検証

## 検証の前提

1. テスト対象のサーバが起動できる。

```bash
# サーバ起動 (Terminal 1)
# テスト項目により BASIC_AUTH_USER_ID/BASIC_AUTH_PASSWORD を変更して起動する
$ BASIC_AUTH_USER_ID=user BASIC_AUTH_PASSWORD=pass go run main.go

# ヘルスチェックへのアクセス (Terminal 2)
$ curl -v -X GET \
-H "User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36" \
"http://localhost:8080/healthz"

# 標準出力へのログ表示 (Terminal 1)
{"Timestamp":"2024-08-04T22:07:19.266757+09:00","Latency":0,"Path":"/healthz","OS":"Windows","Status":200}
```

2. テスト対象のサーバが提供しているパスとBasic認証の有無は以下の通りである。

- `/healthz`: Basic認証なし
- `/todos`: Basic認証なし
- `/api/todos`: Basic認証あり
- `/api/do-panic`: Basic認証あり

3. 認証ミドルウェアは、ミドルウェアの中で一番最後に評価される。これにより、メイン処理の直前に認証処理が行われる。

4. 以下の条件を満たすユーザID/パスワードをサーバに指定した場合、起動せずにエラー終了する。

- 空文字(ユーザID/パスワード)
- コロンを含む文字列(ユーザID)

5. 事前にいくつかデータを登録しておく。

```bash
$ curl -v -X POST \
-H "Content-Type: application/json" \
-d '{"subject": "test", "description": "test"}' "http://localhost:8080/todos"
```

## 検証内容

Markdownでは表現しづらいため、以下のスプレッドシート(要権限)に記録する。

https://docs.google.com/spreadsheets/d/11VtZYNc7xawtvzFL-5idTUP3WT_06aHSutnnbsvB3jU/edit?gid=0#gid=0