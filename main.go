package main

import (
	"FilesToTarBot/bot"
	"FilesToTarBot/config"
	"context"
	"fmt"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	// Load the bot
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("FileToTar version " + config.Version)
	err := telegram.BotFromEnvironment(ctx,
		telegram.Options{
			SessionStorage: &session.FileStorage{Path: "session"},
			UpdateHandler:  bot.Dispatcher,
		},
		bot.RunBot,
		telegram.RunUntilCanceled)
	if err != nil {
		log.Fatalln(err)
	}
}
