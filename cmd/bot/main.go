package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"discord-bot/config"
	"discord-bot/internal/discord"
)

func main() {
	cfg := config.LoadConfig()

	bot, err := discord.NewBot(cfg)
	if err != nil {
		log.Fatal("Error creating bot:", err)
	}

	err = bot.Start()
	if err != nil {
		log.Fatal("Error starting bot:", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	bot.Session.Close()
}
