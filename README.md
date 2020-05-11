# ping-and-report

## ビルド
- `git clone https://github.com/blackbracken/ping-and-report.git`
- `cd ping-and-report`
- `go build`

## 使い方
- config.ymlを記述しバイナリと同じディレクトリに置きます.
- Linuxのみ: 非特権なpingをUDP経由で飛ばすために`sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"`を実行します.
  - ref. [go-ping#note-on-linux-support](https://github.com/sparrc/go-ping#note-on-linux-support)
- 必要に応じてコマンドを実行します.
  - このツールはインスタントなコマンドのみ提供するため, 常時監視等はcronなどを用い行ってください.

## コマンド

| コマンド | 機能 |
| ---- | ---- |
| `ping-and-report ping` | destinationsに指定された各アドレスに対してpingを飛ばし、結果をslackに通知します |
| `ping-and-report stats` | 各アドレスの統計情報をslackに表示します |

## 設定(config.yml)
```config.yml
slack:
  webhookurl: "https://hooks.slack.com/services/xxxxxxxxxxxxxxxxxxxxxxx"
destinations:
  - "xxx.xxx.xxx.xxx"
  - "google.com"
message:
  server_up: |-
    <@SLACK_USER_ID> :signal_strength: The server $address$ is currently up!
      Available: $available_percent$
  server_down: |-
    <@SLACK_USER_ID> :warning: The server $address$ is currently down!
      Available: $available_percent$
  server_stats: |-
    Stats of The server $address$:
      Available: $available_percent$
      UpTime: $up_time$
      TotalMeasuredTime: $total_measured_time$
```
