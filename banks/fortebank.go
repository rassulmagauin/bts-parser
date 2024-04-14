package banks

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

// ParseForte scrapes data from the forte.kz website and returns the scraped data as a string.
func ParseForte() string {
	// Create a new collector
	c := colly.NewCollector(
		colly.AllowedDomains("forte.kz"),
	)

	var results []string

	// Log requests
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	// Targeting specific divs by class and extracting all text
	c.OnHTML("div.sc-bPLjHf.fHJNIN", func(e *colly.HTMLElement) {
		content := fmt.Sprintf("Targeted Content: %s\nClass: %s", e.Text, e.Attr("class"))
		results = append(results, content)
	})

	c.OnHTML("table.MuiTable-root", func(e *colly.HTMLElement) {
		tableContent := e.Text
		var tableResults []string
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			col1 := el.ChildText("td:nth-child(1)")
			col2 := el.ChildText("td:nth-child(2)")
			row := fmt.Sprintf("Column 1: %s, Column 2: %s", col1, col2)
			tableResults = append(tableResults, row)
		})
		tableFull := fmt.Sprintf("%s\n%s", tableContent, strings.Join(tableResults, "\n"))
		results = append(results, tableFull)
	})

	// Start scraping
	c.Visit("https://forte.kz/blue")

	return strings.Join(results, "\n\n")
}
