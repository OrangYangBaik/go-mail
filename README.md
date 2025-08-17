# Priorimail 

Priorimail is a batch program that helps you stay on top of your important emails without constantly checking your inbox. It automatically fetches your Gmail messages, filters them based on your preferences, and forwards the matching ones directly to your Telegram chat using your own bot.

With Priorimail, you can:

- Reduce email clutter
- Get instant Telegram notifications for priority emails
- Customize preferences per user (subject, sender)
- Use your own Telegram bot for privacy

# Features

- Batch Email Fetching: Periodically checks for new messages.
- User Preferences:
  - Filter by email subject
  - Filter by sender address
- Telegram Integration: Forwards filtered emails to Telegram using your personal bot token and chat ID.

# How It Works

1. Priorimail connects to Gmail API to fetch unread emails.

2. It checks each email against the user’s filter criteria (subject & sender).

3. If a match is found, the email content is forwarded to Telegram using the user’s bot token & chat ID.

4. Emails outside the preference criteria are ignored.

# Installation
Requirements

- Frontend: React (Priorimail FE)

- Backend: Go services (go-preference, go-auth, go-mail)

- Database: PostgreSQL (for storing users and preferences)

- Gmail API credentials (credentials.json)

- Telegram Bot API token & chat ID

- LM Studio for local LLM or open API from LLM provider

# Steps
1. Pull the Frontend:
```
git clone https://github.com/OrangYangBaik/priorimail-FE.git
cd priorimail-FE
npm install
npm run dev
```

2. Pull the Backend Repositories:
- Preferences Service:
```
git clone https://github.com/OrangYangBaik/go-preference.git
cd go-preference
go mod tidy
```
- Auth Service:
```
git clone https://github.com/OrangYangBaik/go-auth.git
cd go-auth
go mod tidy
```
- Mail Batch:
```
git clone https://github.com/OrangYangBaik/go-mail.git
cd go-mail
go mod tidy
```

3. Run the preference and auth service then run the batch (can be done manually or scheduler)
