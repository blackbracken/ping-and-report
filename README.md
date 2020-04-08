# ping-and-report

## ビルド
- `git clone https://github.com/blackbracken/ping-and-report.git`
- `cd ping-and-report`
- `go build`

## 使い方
1. config.ymlを記述しバイナリの同じディレクトリに置きます.
2. 必要に応じてコマンドを実行します.
 - このツールはインスタントなコマンドのみ提供するため, 常時監視等はcronなどを用い行ってください.

## コマンド

| コマンド | 機能 |
| ---- | ---- |
| `ping-and-report measure` | 計測を行い結果をslackに通知します |

## 設定
```config.yml
slack:
  webhookurl: "https://hooks.slack.com/services/xxxxxxxxxxxxxxxxxxxxxxx"
  mention: "<!channel>" # <@USER_ID>
destinations:
  - "xxx.xxx.xxx.xxx"
  - "google.com"
```
