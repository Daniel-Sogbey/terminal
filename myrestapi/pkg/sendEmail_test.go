package pkg

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestSendEmail(t *testing.T) {
	err := godotenv.Load("../.env")

	if err != nil {
		t.Error("Test failed")
	}

	ec := NewEmailClient()

	er := &EmailRequest{
		Personalization: []Personalization{
			{
				To: []Email{
					{
						Email: "mathematics06physics@gmail.com",
						Name:  "Daniel",
					},
				},
				Subject: "Hello, World",
			},
		},
		Content: []Content{
			{
				Type:  "text/plain",
				Value: "Hello, From Go Server",
			},
		},
		From: From{
			Email: "gift.alchemy.developer@gmail.com",
			Name:  "Daniel",
		},
		ReplyTo: ReplyTo{
			Email: "gift.alchemy.developer@gmail.com",
			Name:  "Daniel",
		},
	}

	err = ec.sendEmail(er)

	if err != nil {
		t.Errorf("Test failed with err %v", err)
	}
}