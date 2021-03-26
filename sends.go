package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string, replyTo int) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	if replyTo != 0 {
		msg.ReplyToMessageID = replyTo
	}
	return bot.Send(msg)
}

func editMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int, text string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewEditMessageText(
		chatID,
		messageID,
		text,
	)
	return bot.Send(msg)
}
