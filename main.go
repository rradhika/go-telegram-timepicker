package main

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tp "github.com/rradhika/timepicker/timepicker"
)

var myBot *tgbotapi.BotAPI
var updChannel tgbotapi.UpdatesChannel

func main() {

	// db, err := datastore.NewDB()
	bot, err := tgbotapi.NewBotAPI("YOUR-BOT-KEY")

	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updChannel, err = bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatalln(err)
	}

	for update := range updChannel {

		switch {

		case update.Message != nil && update.Message.Text == `/timepicker`:

			msg := fmt.Sprintln("Pilih Jam Mulai")
			setMsg := "Set Jam Mulai"

			keyboard := tp.CreateTimepicker(setMsg, tp.JamPertama, tp.MenitPertama)

			reply := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
			reply.ReplyMarkup = &keyboard
			bot.Send(reply)

		case update.CallbackQuery != nil:
			data := update.CallbackQuery.Data
			setMsg := "Set Jam Mulai"
			switch {
			case strings.Contains(data, "UPDATE-JAM"):
				resData := tp.SeparateCallbackData(data)

				var msg = tgbotapi.NewEditMessageText(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
					update.CallbackQuery.Message.Text,
				)
				var msgMarkup = tgbotapi.NewEditMessageReplyMarkup(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
					tp.CreateTimepicker(setMsg, resData[1], resData[2]),
				)
				msg.ReplyMarkup = msgMarkup.ReplyMarkup
				bot.Send(msg)

			case strings.Contains(data, "SET-JAM"):
				bot.AnswerCallbackQuery(
					tgbotapi.CallbackConfig{
						CallbackQueryID: update.CallbackQuery.ID,
						Text:            "Jam Telah Dipilih",
					},
				)
				resData := tp.SeparateCallbackData(data)
				msg := fmt.Sprintf("Anda Memilih Jam: %s:%s", resData[1], resData[2])
				reply := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msg)
				bot.Send(reply)
			}

		}
	}

}
