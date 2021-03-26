package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	stringify "github.com/vicanso/go-stringify"
)

const (
	// APIEndpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint is the endpoint for downloading a file from Telegram.
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"
)

func main() {
	creators := strings.Join(getBotCreators(), ", ")
	log.Printf("QuotesAPI Bot by %s", creators)
	log.Printf("Loading Golang QuotesAPI Bot Alpha...")
	bot, err := tgbotapi.NewBotAPI(getBotToken())
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)
	sendMessage(bot, getReportsChat(), "Started successfully!", 0)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		switch update.Message.Command() {
		case "start":
			sendMessage(bot, update.Message.Chat.ID, "Hi there!", 0)
		case "dmsg":
			data := stringify.String(update.Message, replacer)
			sendMessage(bot, update.Message.Chat.ID, data, update.Message.MessageID)
		case "q":
			if update.Message.ReplyToMessage == nil {
				sendMessage(bot, update.Message.Chat.ID, "Reply not found", 0)
				continue
			}
			if update.Message.ReplyToMessage.Text == "" && update.Message.ReplyToMessage.Photo == nil && update.Message.ReplyToMessage.Sticker == nil {
				sendMessage(bot, update.Message.Chat.ID, "No text or media in reply found", 0)
				continue
			}
			message, err := sendMessage(bot, update.Message.Chat.ID, "Processing quote...", update.Message.MessageID)
			if err != nil {
				log.Panic(err)
			}
			username := ""
			colors := [8]string{
				"#fb6169",
				"#85de85",
				"#f3bc5c",
				"#65bdf3",
				"#b48bf2",
				"#ff5694",
				"#62d4e3",
				"#faa357",
			}
			num1 := update.Message.ReplyToMessage.From.ID % 7
			num2 := [7]int{0, 7, 4, 1, 6, 3, 5}
			color := colors[num2[num1]]
			log.Println(update.Message.ReplyToMessage.Entities)
			text := update.Message.ReplyToMessage.Text
			if update.Message.ReplyToMessage.Entities != nil {
				text = parseEntities(text, update.Message.ReplyToMessage.Entities)
			}
			noPFP := "True"
			PFP := ""
			args := strings.Split(update.Message.Text, " ")
			style := "desk"
			adminTitle := ""
			mediaUrl := ""
			if len(args) > 1 {
				style = args[len(args)-1]
			}
			// golang are best lang
			if update.Message.ReplyToMessage.From.LastName == "" {
				username = update.Message.ReplyToMessage.From.FirstName
			} else {
				username = update.Message.ReplyToMessage.From.FirstName + " " + update.Message.ReplyToMessage.From.LastName
			}
			// smart things haunted him but he was faster
			if update.Message.ReplyToMessage.Photo != nil {
				fileID := getFileID(update.Message.ReplyToMessage.Photo)
				mediaUrl, err = bot.GetFileDirectURL(fileID)
				if err != nil {
					log.Panic(err)
				}
			} else if update.Message.ReplyToMessage.Sticker != nil {
				fileID := getStickerFileID(update.Message.ReplyToMessage.Sticker)
				mediaUrl, err = bot.GetFileDirectURL(fileID)
				if err != nil {
					log.Panic(err)
				}
			}
			photo := tgbotapi.NewUserProfilePhotos(
				update.Message.ReplyToMessage.From.ID,
			)
			photos, err := bot.GetUserProfilePhotos(photo)
			if err != nil {
				log.Panic(err)
			}
			if photos.TotalCount == 0 {
				PFP = ""
			} else {
				ProfilePic, err := bot.GetFileDirectURL(photos.Photos[0][0].FileID)
				if err != nil {
					log.Panic(err)
				}
				PFP = ProfilePic
				noPFP = "False"
			}
			values := url.Values{}
			values.Set("chat_id", fmt.Sprintf("%v", update.Message.Chat.ID))
			values.Set("user_id", fmt.Sprintf("%v", update.Message.ReplyToMessage.From.ID))
			member, err := bot.MakeRequest(
				"getChatMember",
				values,
			)
			if err != nil {
				log.Panic(err)
			} else {
				chatMember := parseAnswer(member.Result)
				if chatMember["custom_title"] != nil {
					// admintitle
					adminTitle = chatMember["custom_title"].(string)
				}
			}
			quoteText := quotes(username, text, noPFP, PFP, color, adminTitle, style, mediaUrl)
			quoteAnswer := parseQuoteAnswer([]byte(quoteText))
			if quoteAnswer.TokenInvalid != "" {
				editMessage(bot, update.Message.Chat.ID, message.MessageID, "API Token invalid!")
				continue
			} else if quoteAnswer.InvalidTemplate != "" {
				editMessage(bot, update.Message.Chat.ID, message.MessageID, "Invalid Template!")
				continue
			} else if quoteAnswer.AccessDenied != "" {
				editMessage(bot, update.Message.Chat.ID, message.MessageID, "Template Access Denied!")
				continue
			} else {
				editMessage(bot, update.Message.Chat.ID, message.MessageID, "Sending quote...")
				quoteURL := parseQuoteURL([]byte(quoteText))
				quotePath := saveToFile(quoteURL)
				deleteMsg := tgbotapi.NewDeleteMessage(
					message.Chat.ID,
					message.MessageID,
				)
				bot.Send(deleteMsg)
				log.Println(quotePath)
				sendFile := tgbotapi.NewStickerUpload(
					update.Message.Chat.ID,
					quotePath,
				)
				sendFile.ReplyToMessageID = update.Message.MessageID
				bot.Send(sendFile)
				removeFile(quotePath)
			}
		}
	}
}
