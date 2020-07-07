package notify

import (
	"context"
	"encoding/json"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"os"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type Build struct {
	Id         string   `json:"id"`
	ProjectId  string   `json:"projectId"`
	Status     string   `json:"status"`
	StartTime  string   `json:"startTime"`
	FinishTime string   `json:"finishTime"`
	LogUrl     string   `json:"logUrl"`
	Images     []string `json:"images"`
	Source     Source
}

type Source struct {
	RepoSource RepoSource `json:"repoSource"`
}

type RepoSource struct {
	RepoName   string
	BranchName string
}

func NotifyBuild(ctx context.Context, m PubSubMessage) error {

	log.Printf("Message is %s", m.Data)

	var result Build
	if err := json.Unmarshal(m.Data, &result); err != nil {
		return err
	}
	log.Printf("Decoded message is %s", result.Id)

	config := Config{
		EmailTo:   os.Getenv("EMAIL_TO"),
		EmailFrom: os.Getenv("EMAIL_FROM")}

	var err error
	var email *Email
	if email, err = CreateEmail(config, result); err != nil {
		return err
	}

	log.Printf("%s", *email)
	return send(config, email)
}

var send = func(config Config, email *Email) error {
	from := mail.NewEmail("Cloud Build", config.EmailFrom)
	subject := email.Subject
	to := mail.NewEmail("Cloud Build", config.EmailTo)
	message := mail.NewSingleEmail(from, subject, to, email.Text, email.Html)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return err
	} else {
		log.Printf("Build notification complete %d, body=%s, headers=%s", response.StatusCode, response.Body, response.Headers)
	}
	return nil
}
