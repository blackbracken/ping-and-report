# ping-and-report
A slack bot to report life & death of servers

# how to use
Call the built binary via cron.

# configure
Put config.yml into a same directory as the binary.

Example:
```config.yml
slack:
  webhookurl: "https://hooks.slack.com/services/xxxxxxxxxxxxxxxxxxxxxxx"
  mention: "<!channel>" # <@USER_ID>
pinged:
  - "xxx.xxx.xxx.xxx"
  - "google.com"
```