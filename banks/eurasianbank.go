package banks

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

// ParseEUBank scrapes data from the eubank.kz website under the bonus program page.
func ParseEUBank() string {
	// Create a new collector
	c := colly.NewCollector(
		colly.AllowedDomains("eubank.kz"),
	)

	var results []string

	// Log requests
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	// Extract data from the 'advantages__item' blocks
	c.OnHTML("ul.advantages", func(e *colly.HTMLElement) {
		e.ForEach("li.advantages__item", func(_ int, el *colly.HTMLElement) {
			category := el.ChildText("div.advantages__title")
			details := el.ChildText("div.advantages__text")
			row := fmt.Sprintf("Category: %s, Details: %s", category, details)
			results = append(results, row)
		})
	})

	// Start scraping the specific URL
	c.Visit("https://eubank.kz/bonus-program/?lang=en")

	return strings.Join(results, "\n")
}
