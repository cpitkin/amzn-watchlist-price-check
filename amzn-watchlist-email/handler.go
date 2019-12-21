// Copyright (c) 2019 Charlie Pitkin
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package function

import (
	"encoding/json"
	"io/ioutil"

	handler "github.com/openfaas-incubator/go-function-sdk"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	var err error

	sendgridToken, err := ioutil.ReadFile("/var/openfaas/secrets/sendgridtoken")
	if err != nil {
		return handler.Response{}, err
	}

	type EmailBooks struct {
		FetchedTime string `json:"fetchedTime"`
		Books       string `json:"books"`
	}

	var em EmailBooks
	json.Unmarshal(req.Body, &em)

	m := mail.NewV3Mail()
	p := mail.NewPersonalization()

	from := mail.NewEmail("Charlie Pitkin", "charlie.pitkin@gmail.com")
	m.SetFrom(from)
	p.Subject = "Amazon Price Check"

	tos := []*mail.Email{
		mail.NewEmail("Celeste Lempke", "celeste.lempke@gmail.com"),
		mail.NewEmail("Charlie Pitkin", "charlie.pitkin@gmail.com"),
	}
	p.AddTos(tos...)

	c := mail.NewContent("text/plain", em.FetchedTime)
	m.AddContent(c)

	c = mail.NewContent("text/html", em.Books)
	m.AddContent(c)

	client := sendgrid.NewSendClient(string(sendgridToken))
	res, err := client.Send(m)
	if err != nil {
		return handler.Response{}, err
	}

	return handler.Response{
		Header: map[string][]string{
			"Status": []string{string(res.StatusCode)},
		},
	}, err
}
