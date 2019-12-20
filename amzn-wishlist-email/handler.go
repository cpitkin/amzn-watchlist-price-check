package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

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
		fetchedTime string `json:"fetchedTime"`
		books       string `json:"books"`
	}

	var em EmailBooks
	json.Unmarshal(req.Body, &em)

	from := mail.NewEmail("Charlie Pitkin", "charlie.pitkin@gmail.com")
	subject := "Amazon Price Check"

	tos := []*mail.Email{
		mail.NewEmail("Celeste Lempke", "celeste.lempke@gmail.com"),
		mail.NewEmail("Charlie Pitkin", "charlie.pitkin@gmail.com"),
	}

	plainTextContent := em.fetchedTime

	htmlContent := em.books

	message := mail.NewSingleEmail(from, subject, tos, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sendgridToken)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return handler.Response{mail.GetRequestBody(m)}, err
}
