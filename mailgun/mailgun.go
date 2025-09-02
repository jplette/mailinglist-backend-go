package mailgun

import (
	"context"
	"mailinglist-backend-go/config"
	"time"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

var domain = config.Value("MAILGUN_DOMAIN")
var apiKey = config.Value("MAILGUN_API_KEY")

func MailingLists() []string {
	mg := mailgun.NewMailgun(apiKey)
	err := mg.SetAPIBase(mailgun.APIBaseEU)
	if err != nil {
		return nil
	}

	listIterator := mg.ListMailingLists(&mailgun.ListOptions{Limit: 100})

	var lists []string

	var page []mtypes.MailingList
	// The entire operation should not take longer than 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	for listIterator.Next(ctx, &page) {
		for _, list := range page {
			lists = append(lists, list.Address)
		}
	}
	return lists
}
