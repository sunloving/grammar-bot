package spell

import (
	"fmt"
	"gamma-rho-bot/telegram"
	"log"
)

type listener struct {
	telegramClient telegram.BotAPIClient
	updatesOffset  int64
	chatsIds       map[int64]struct{}
	error          chan error
}

func (l *listener) start(message chan chatMessage) {
	for {
		log.Print("getUpdates started...")
		updates, err := l.getUpdates()
		if err != nil {
			l.error <- fmt.Errorf("can't get updates from telegram: %s", err.Error())
			continue
		}

		log.Print("getUpdates finished")

		if len(updates) == 0 {
			continue
		}

		l.updatesOffset = updates[len(updates)-1].Id + 1

		for _, update := range updates {
			if l.chatsIds != nil {
				if _, ok := l.chatsIds[update.Message.Chat.Id]; !ok {
					continue
				}
			}

			message <- chatMessage{
				id:     update.Message.Id,
				chatId: update.Message.Chat.Id,
				text:   update.Message.Text,
			}
		}
	}
}

func (l *listener) getUpdates() ([]telegram.Update, error) {
	return l.telegramClient.GetUpdates(
		l.updatesOffset,
		100,
		60,
		[]string{"messages"},
	)
}
