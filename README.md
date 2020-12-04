# serviceCheck_golang

Program written in Golang to check status of services.

Services and intervals are stored in a DB. This program check services and alerts a channel on Slack and Telegram if status changes.

Set webhookURLSlack and webhookTelegram variables with your own webhooks to send messages.
