package helpers

import (
	"os"

	"github.com/plivo/plivo-go/v7"
)

func InitPlivo() (*plivo.Client, error) {
	return plivo.NewClient(os.Getenv("PLIVO_AUTH_ID"), os.Getenv("PLIVO_AUTH_TOKEN"), &plivo.ClientOptions{})
}

func SendSMS(client *plivo.Client, destination string, message string) (*plivo.MessageCreateResponseBody, error) {
	return client.Messages.Create(
		plivo.MessageCreateParams{
			Src:  os.Getenv("SMS_SENDER_ID"),
			Dst:  destination,
			Text: message,
		},
	)
}
