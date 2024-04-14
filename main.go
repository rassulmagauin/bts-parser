package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/rassulmagauin/webscraper/banks"
	"github.com/rassulmagauin/webscraper/db"
	"github.com/rassulmagauin/webscraper/gpt"

	"github.com/joho/godotenv"
)

func helper(offers []gpt.Offer, bankId int, cardType string) []gpt.Offer {
	for i := range offers {
		offers[i].CardType = cardType
		offers[i].BankID = bankId
	}
	return offers
}

func Updater(gptClient gpt.GPTClient, driver *db.DbDriver, ctx context.Context) {
	jusan := banks.ParseJusan()
	forte := banks.ParseForte()
	eurasian := banks.ParseEUBank()
	halyk := banks.Halyk
	bereke := banks.Bereke
	bereke_offers, _ := gptClient.AnalyzeOffers(5, "allin", bereke)
	halyk_offers, _ := gptClient.AnalyzeOffers(4, "gold", halyk)
	jusan_offers, _ := gptClient.AnalyzeOffers(3, "standart", jusan)
	eurasian_offers, _ := gptClient.AnalyzeOffers(1, "metal", eurasian)
	forte_offers, _ := gptClient.AnalyzeOffers(2, "blue", forte)

	jusan_offers = helper(jusan_offers, 3, "standart")
	eurasian_offers = helper(eurasian_offers, 1, "metal")
	forte_offers = helper(forte_offers, 2, "blue")
	halyk_offers = helper(halyk_offers, 4, "gold")
	bereke_offers = helper(bereke_offers, 5, "allin")

	driver.UpdateOrCreateOffers(ctx, jusan_offers)
	driver.UpdateOrCreateOffers(ctx, eurasian_offers)
	driver.UpdateOrCreateOffers(ctx, forte_offers)
	driver.UpdateOrCreateOffers(ctx, halyk_offers)
	driver.UpdateOrCreateOffers(ctx, bereke_offers)
	//fmt.Println(jusan_offers)
}

func main() {
	godotenv.Load()
	api := os.Getenv("API_KEY")
	driver := db.NewDBDriver()
	defer driver.Close()
	if err := driver.Connect(); err != nil {
		log.Fatal(err)
	}
	if api == "" {
		log.Fatalln("You haven't API key")
	}
	ctx := context.Background()
	gptClient := gpt.NewClient(api, ctx)
	for {
		Updater(gptClient, driver, ctx)
		time.Sleep(24 * time.Hour) // Sleep for 24 hours
	}
}
