// Copyright (c) 2019 Charlie Pitkin
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package function

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	handler "github.com/openfaas-incubator/go-function-sdk"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailBooks struct {
	FetchedTime string `json:"fetchedTime"`
	Books       string `json:"books"`
}

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	var err error

	sendgridToken, err := ioutil.ReadFile("/var/openfaas/secrets/sendgridtoken")
	if err != nil {
		return handler.Response{}, err
	}

	fromEmailInput, err := ioutil.ReadFile("/var/openfaas/secrets/fromemail")
	if err != nil {
		return handler.Response{}, err
	}

	toEmailInput, err := ioutil.ReadFile("/var/openfaas/secrets/toemail")
	if err != nil {
		return handler.Response{}, err
	}

	var from *mail.Email
	json.Unmarshal(fromEmailInput, &from)

	var tos []*mail.Email
	json.Unmarshal(toEmailInput, &tos)

	var em EmailBooks
	json.Unmarshal(req.Body, &em)

	m := mail.NewV3Mail()
	p := mail.NewPersonalization()

	m.SetFrom(from)

	p.Subject = "Amazon Watchlist Price Check"

	p.AddTos(tos...)

	c := mail.NewContent("text/plain", em.FetchedTime)
	m.AddContent(c)

	c = mail.NewContent("text/html", em.Books)
	m.AddContent(c)

	m.AddPersonalizations(p)

	client := sendgrid.NewSendClient(strings.TrimSpace(string(sendgridToken)))
	res, err := client.Send(m)
	if err != nil {
		return handler.Response{}, err
	}

	return handler.Response{
		Header: map[string][]string{
			"Status": []string{string(res.StatusCode)},
		},
		Body: []byte(res.Body),
	}, err
}
