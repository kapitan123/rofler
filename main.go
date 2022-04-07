package main

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

//const postCommand, startCommand, topCommand = "/post", "/start", "/top"
//const topRofler, topWorker, topTaste = "rofler", "worker", "taste"
//var lenPostCommand, lenStartCommand = len(postCommand), len(startCommand)
//const botTag = "@TelegroflerBot"
//var lenBotTag = len(botTag)
//const telegramTokenEnv = "TELEGRAM_BOT_TOKEN"
//var telegramToken = os.Getenv(telegramTokenEnv)
//const telegramApiSendMessage = "/sendMessage"

func main() {
	//https://vm.tiktok.com/ZSdM14SCd
	// regexp to match
	// ^https:\/\/vm\.tiktok\.com\/.*

	DownloadTikTok("https://vm.tiktok.com/ZSdM14SCd")
}

func DownloadTikTok(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	fileName := GenerateRandName()
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, res.Body)

	return err
}

func GenerateRandName() string {
	randName := "tiktok-"

	rand.Seed(time.Now().UnixNano())
	ints := rand.Perm(5)

	for _, v := range ints {
		randName += strconv.Itoa(v)
	}

	return randName + ".mp4"
}
