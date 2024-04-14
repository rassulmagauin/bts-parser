package banks

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

// ParseJusan scrapes data from the jusan.kz website and returns the scraped data as a string.
func ParseJusan() string {
	// Create a new collector
	c := colly.NewCollector(
		colly.AllowedDomains("jusan.kz"),
	)

	var results []string

	// Log requests
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	// Targeting the specific divs by class for the product details and extracting relevant table data
	c.OnHTML("div.product-tariff_product_tariff___TgPH", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			col1 := el.ChildText("td:nth-child(1)")
			col2 := el.ChildText("td:nth-child(2)")
			row := fmt.Sprintf("Detail: %s, Value: %s", col1, col2)
			results = append(results, row)
		})
	})

	// Start scraping
	c.Visit("https://jusan.kz/card/jusan/tariff")

	return strings.Join(results, "\n")
}
