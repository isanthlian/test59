package main

import (
	"log"
	"os"
	"os/exec"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DownloadResponse struct {
	Url string `json:"url"`
}

func main() {
	bot, err := tgbotapi.NewBotAPI("7773439673:AAHWw7iiGT5g6PoYEljzdhRZoZiRpjRQbnc")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	linkRegex := regexp.MustCompile(`https?://(www\.)?instagram\.com/reel/\w+`)

	log.Println("Бот запущен и ожидает сообщения...")

	for update := range updates {
		if update.Message != nil {
			text := update.Message.Text
			log.Println("Получено сообщение:", text)

			if linkRegex.MatchString(text) {
				log.Println("Найдена ссылка на Instagram!")
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "🔍 Обнаружена ссылка! Скачиваю..."))

				filePath, err := downloadVideo(text)
				if err != nil {
					log.Println("Ошибка при загрузке видео:", err)
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Ошибка при загрузке видео."))
					continue
				}

				log.Println("Видео скачано, отправляю...")
				video := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FilePath(filePath))
				bot.Send(video)
				os.Remove(filePath)
			}
		}
	}
}

func downloadVideo(instagramURL string) (string, error) {
	filePath := "video.mp4"
	log.Println("Скачивание видео с помощью yt-dlp:", instagramURL)
	cmd := exec.Command("yt-dlp", "-o", filePath, instagramURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Println("Ошибка при скачивании видео:", err)
		return "", err
	}

	return filePath, nil
}
