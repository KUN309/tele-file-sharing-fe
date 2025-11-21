package main

import (
	"fmt"
	"log"

	"fe-file-sharing/internal/api"
	"fe-file-sharing/internal/bot"
	"fe-file-sharing/internal/config"
)

//! Trong quá trình (WIP)
func main() {
    config.Load("") // Load once
	fmt.Println("Tải config thành công")

	newClient := api.NewClient(
		config.C.BackendAPIBase,
		config.C.TelegramBotToken,
	)
	fmt.Println("Tạo client mới thành công")
    
	newBot, err := bot.NewBot(
        config.C.TelegramBotToken,
		newClient,
    )
	if err != nil {
		log.Fatalf("Tạo bot mới thất bại: %v", err)
    }
	fmt.Println("Tạo bot mới thành công")

    newBot.Start()
	fmt.Println("Bot đã bắt đầu chạy")
}