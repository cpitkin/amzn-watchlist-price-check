// Copyright (c) 2019 Charlie Pitkin
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	handler "github.com/openfaas-incubator/go-function-sdk"

	"github.com/gocolly/colly"
)

type Book struct {
	Name  string
	Price string
	Link  string
}

// Collect all the books below a set price
var allBooks = []Book{}

// Price to the nearest whole number
var maxPrice int64 = 2

func parseList(e *colly.HTMLElement) {
	e.ForEach("li.a-spacing-none", func(_ int, e *colly.HTMLElement) {

		wholePrice := strings.Split(e.ChildText("span.a-price-whole"), ".")

		numPrice, err := strconv.ParseInt(wholePrice[0], 10, 64)
		if err != nil {
			panic(err)
		}

		if numPrice <= maxPrice {
			book := Book{
				Name:  e.ChildAttr("a.a-link-normal", "title"),
				Price: e.ChildText("span.a-price > span.a-offscreen"),
				Link:  "https://smile.amazon.com" + e.ChildAttr("a.a-link-normal", "href"),
			}

			allBooks = append(allBooks, book)
		}
	})
}

func Handle(req handler.Request) (handler.Response, error) {
	wishlistIdsString, err := ioutil.ReadFile("/var/openfaas/secrets/wishlistIds")
	if err != nil {
		return handler.Response{}, err
	}

	zapierWebhook, err := ioutil.ReadFile("/var/openfaas/secrets/zapierWebhook")
	if err != nil {
		return handler.Response{}, err
	}

	wishlistIds := strings.Split(string(wishlistIdsString), ",")

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:68.0) Gecko/20100101 Firefox/68.0"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*.amazon.*",
		Parallelism: 1,
		Delay:       2 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	c.OnHTML("#g-items", func(e *colly.HTMLElement) {
		parseList(e)
		seeMoreLink := e.ChildAttr("a.wl-see-more", "href")
		if seeMoreLink != "" {
			c.Visit("https://smile.amazon.com" + seeMoreLink)
		}
	})

	for _, id := range wishlistIds {
		c.Visit("https://smile.amazon.com/hz/wishlist/ls/" + id)
	}

	emailString := ""

	for _, book := range allBooks {
		emailString = emailString + "<br><b>Name:</b> " + book.Name + "<br><b>Price:</b> " + book.Price + "<br><a href=\"" + book.Link + "\">Buy Now</a><br>"
	}

	timeString := time.Now()

	fetchedTime := "Books fetched on " + string(timeString.Format(time.UnixDate)+"<br>")

	jsonValues := map[string]string{"fetchedTime": fetchedTime, "books": emailString}

	jsonValue, _ := json.Marshal(jsonValues)

	res, err := http.Post(string(zapierWebhook), "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return handler.Response{}, err
	}

	return handler.Response{
		Body: []byte(res.Status),
	}, err
}
