---
title: Telegram Bot and Channel
description: Create the Telegram bot and channel used for MiniPaaS monitoring alerts.
weight: 7
---

## Overview

The MiniPaaS role sends monitoring reports and log alerts to Telegram. It needs two values:

- `telegram_bot_token`: the HTTP API token of the bot that posts alerts
- `telegram_chat_id`: the channel or group that receives them

When both are set, the role installs the monitoring script and writes `/etc/telegram.secrets` with `BOT_TOKEN` and `GROUP_ID`. See **[Monitoring](/role/monitoring/)**.

This guide covers creating the bot, creating the channel, and giving the role its credentials.

---

## Prerequisites

- A Telegram account
- MiniPaaS installed with the role; see **[Role Installation](/role/installation/)**

---

## 1. Create the Bot

1. Open Telegram and chat with [@BotFather](https://t.me/BotFather).
2. Send `/newbot` and follow the prompts. Choose a display name and a username ending in `bot`.
3. Copy the token BotFather returns. It looks like `123456:ABC-DEF...`. This is `telegram_bot_token`.

Keep the token private. Anyone with it can post as the bot.

---

## 2. Create the Channel

1. In Telegram, create a channel for alerts. A group works the same way.
2. Add the bot to the channel as an administrator. Open the channel, select Administrators, and add the bot. A bot can post only when the chat grants it permission.
3. Post a message in the channel so the chat has an update the bot can read.

---

## 3. Find the Chat ID

1. Open this URL in a browser, replacing the token:

```
https://api.telegram.org/bot<YourBotToken>/getUpdates
```

2. Find `"chat":{"id":` in the response. The number that follows is `telegram_chat_id`. Copy it exactly.

Channels and groups use negative IDs. The leading minus is part of the value. If the response is empty, post another message in the channel and reload the URL.

---

## 4. Give the Credentials to the Role

Set both variables in `group_vars/`, `host_vars/`, or `inventory.ini`:

```yaml
telegram_bot_token: "123456:ABC-DEF..."
telegram_chat_id: "-1001234567890"
```

The variables default to the `TELEGRAM_TOKEN` and `TELEGRAM_CHAT` environment variables. Store the token in Ansible Vault rather than plain text. See **[Configuration](/role/configuration/)**.

The role activates Telegram when both values are set. Re-run the playbook to install the monitoring script and write `/etc/telegram.secrets`.

---

## 5. Verify

The role installs `/usr/local/bin/system-monitoring.sh` and starts it at boot through `/etc/cron.d/monitoring-telegram`. The script sends a boot report. Watch the channel for the message.

Check the secrets file on a node:

```bash
sudo ls -l /etc/telegram.secrets
```

The file is owned by `root:root` with mode `0600`.

---

## Notes

- A bot posts to a channel only when the channel grants it permission.
- The role reads the chat ID as `telegram_chat_id` and writes it as `GROUP_ID`.
- Change the token in BotFather at any time. Update `telegram_bot_token` and re-run the playbook.
- The monitoring script also forwards container errors. See **[Monitoring](/role/monitoring/)**.

---

## Next Steps

- **[Monitoring](/role/monitoring/)**: alerts, syslog forwarding, and thresholds.
- **[Variables](/role/variables/)**: every role variable.
- **[Role Installation](/role/installation/)**: provision a cluster.
