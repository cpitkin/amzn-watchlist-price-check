// Copyright (c) 2019 Charlie Pitkin
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Book struct {
	Name  string
	Price string
	Link  string
}

func main() {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:68.0) Gecko/20100101 Firefox/68.0"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  ".*amazon.*",
		Parallelism: 1,
		Delay:       1 * time.Second,
	})

	allBooks := []Book{}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	c.OnHTML("#g-items", func(e *colly.HTMLElement) {
		e.ForEach("li.a-spacing-none", func(_ int, e *colly.HTMLElement) {

			wholePrice := strings.Split(e.ChildText("span.a-price-whole"), ".")

			numPrice, err := strconv.ParseInt(wholePrice[0], 10, 0)
			if err != nil {
				panic(err)
			}

			if numPrice <= 3 {
				book := Book{
					Name:  e.ChildAttr("a", "title"),
					Price: e.ChildText("span.a-price > span.a-offscreen"),
					Link:  "https://smile.amazon.com" + e.ChildAttr("a", "href"),
				}

				fmt.Printf("\nName: %v\n", book.Name)
				fmt.Printf("Price: %v\n", book.Price)
				fmt.Printf("Link: %v\n\n", book.Link)

				allBooks = append(allBooks, book)
			}
		})

	})

	c.Visit("https://smile.amazon.com/hz/wishlist/ls/1DJLN9PNW8R59")
}
