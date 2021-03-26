package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

func quotes(username string, text string, noPfp string, Pfp string, color string, adminTitle string, style string, media string) string {
	return makeRequest(username, text, noPfp, Pfp, color, adminTitle, style, media)
}

func makeRequest(username string, text string, noPfp string, Pfp string, color string, adminTitle string, style string, media string) string {
	// there is no reply because of telegram bot api limitations
	// if you have a non-mtproto golang solution - pm me pls (tg: @rf0x1d)
	toSend := url.Values{
		"no_pfp":     {noPfp},
		"pfp":        {Pfp},
		"username":   {username},
		"raw_text":   {text},
		"colour":     {color},
		"admintitle": {adminTitle},
		"style":      {style},
		"token":      {getQuotesToken()},
		"mediaImage": {media},
	}

	resp, err := http.PostForm("https://api.rf0x3d.su/dev/quote", toSend)

	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	log.Println(resp.Status)
	log.Println(string(respBody))
	return string(respBody)
}

func parseQuoteAnswer(ans []byte) QuoteAnswer {
	var answer QuoteAnswer
	json.Unmarshal(ans, &answer)
	return answer
}

func parseQuoteURL(js []byte) string {
	var ans QuoteAnswer
	json.Unmarshal(js, &ans)
	return ans.Success.File
}

// Some shitcode due to 0iq solutions
func getFileID(photos *[]tgbotapi.PhotoSize) string {
	lastSize := 0
	fileID := ""
	for _, photo := range *photos {
		if photo.FileSize > lastSize {
			log.Println(photo.FileSize)
			lastSize = photo.FileSize
			fileID = photo.FileID
		}
	}
	return fileID
}

func getStickerFileID(stickers *tgbotapi.Sticker) string {
	return stickers.FileID
}

func parseEntities(text string, entities *[]tgbotapi.MessageEntity) string {
	textSplitted := strings.Split(text, "")
	log.Println(textSplitted)
	for _, entity := range *entities {
		start := entity.Offset
		end := start + (entity.Length - 1)
		if entity.Type == "code" || entity.Type == "pre" {
			textSplitted[start] = "{c}" + textSplitted[start]
		} else if entity.Type == "mention" || entity.Type == "url" || entity.Type == "bot_command" || entity.Type == "hashtag" || entity.Type == "email" || entity.Type == "text_link" || entity.Type == "phone_number" || entity.Type == "text_mention" {
			textSplitted[start] = "{l}" + textSplitted[start]
		} else if entity.Type == "bold" {
			textSplitted[start] = "{b}" + textSplitted[start]
		} else if entity.Type == "italic" {
			textSplitted[start] = "{i}" + textSplitted[start]
		} else if entity.Type == "strikethrough" {
			textSplitted[start] = "{s}" + textSplitted[start]
		} else if entity.Type == "underline" {
			textSplitted[start] = "{u}" + textSplitted[start]
		} else {
			log.Println("unknown entity!!!")
		}
		textSplitted[end] = textSplitted[end] + "{cl}"
	}
	return strings.Join(textSplitted, "")
}

func saveToFile(URL string) string {
	resp, err := http.Get(URL)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	uid, err := uuid.NewUUID()
	if err != nil {
		log.Panic(err)
	}

	uuid := fmt.Sprintf("quotes/%s.webp", uid)

	file, err := os.Create(uuid)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Panic(err)
	}
	return uuid
}

func removeFile(name string) {
	err := os.Remove(name)
	if err != nil {
		log.Panic(err)
	}
}
