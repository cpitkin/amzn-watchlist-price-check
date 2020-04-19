// Copyright (c) 2020 Charlie Pitkin
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
)

type Book struct {
	Name  string
	Price string
	Link  string
}

var sugar = zap.NewExample().Sugar()

// Price to the nearest whole number
var maxPrice int64 = 2

func parseList(e *colly.HTMLElement) []Book {

	var newBooks []Book

	e.ForEach("li.g-item-sortable", func(_ int, e *colly.HTMLElement) {

		wholePrice := strings.Split(e.ChildText("span.a-price-whole"), ".")

		if wholePrice[0] != "" {
			numPrice, err := strconv.ParseInt(wholePrice[0], 10, 64)
			if err != nil {
				sugar.Error(err)
			}

			if numPrice <= maxPrice {
				book := Book{
					Name:  e.ChildAttr("a.a-link-normal", "title"),
					Price: e.ChildText("span.a-price > span.a-offscreen"),
					Link:  "https://smile.amazon.com" + e.ChildAttr("a.a-link-normal", "href"),
				}

				newBooks = append(newBooks, book)
			}
		}
	})

	return newBooks
}

func main() {
	var err error
	var allBooks, newBooksList []Book

	defer sugar.Sync()

	wishlistIdsString, err := ioutil.ReadFile("/var/amzn_wishlist_price_check/secrets/wishlistids")
	if err != nil {
		sugar.Error(err)
	}

	wishlistIds := strings.Split(string(wishlistIdsString), ",")

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:70.0) Gecko/20100101 Firefox/70.0"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*.amazon.*",
		Parallelism: 1,
		Delay:       1 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		sugar.Debugf("Visiting: %s", r.URL.String())
	})

	c.OnHTML("ul#g-items", func(e *colly.HTMLElement) {
		newBooksList = parseList(e)

		if len(newBooksList) != 0 {
			allBooks = append(allBooks, newBooksList...)
		}

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

	sendEmail(emailString)
}

func sendEmail(emailString string) {
	sendgridToken, err := ioutil.ReadFile("/var/amzn_wishlist_price_check/secrets/sendgridtoken")
	if err != nil {
		sugar.Error(err)
	}

	fromEmailInput, err := ioutil.ReadFile("/var/amzn_wishlist_price_check/secrets/fromemail")
	if err != nil {
		sugar.Error(err)
	}

	toEmailInput, err := ioutil.ReadFile("/var/amzn_wishlist_price_check/secrets/toemail")
	if err != nil {
		sugar.Error(err)
	}

	var from *mail.Email
	json.Unmarshal(fromEmailInput, &from)

	var tos []*mail.Email
	json.Unmarshal(toEmailInput, &tos)

	m := mail.NewV3Mail()
	p := mail.NewPersonalization()

	m.SetFrom(from)

	p.Subject = "Amazon Watchlist Price Check"

	p.AddTos(tos...)

	c := mail.NewContent("text/html", emailString)
	m.AddContent(c)

	m.AddPersonalizations(p)

	client := sendgrid.NewSendClient(strings.TrimSpace(string(sendgridToken)))
	_, err = client.Send(m)
	if err != nil {
		sugar.Error(err)
	}

	sugar.Infof("Email sent: %s", time.Now())
}
