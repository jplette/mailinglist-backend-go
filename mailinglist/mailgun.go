package mailinglist

import (
	"context"
	"mailinglist-backend-go/envcfg"
	"time"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

var domain = envcfg.Value("MAILGUN_DOMAIN")
var apiKey = envcfg.Value("MAILGUN_API_KEY")

func Lists() ([]mtypes.MailingList, error) {
	mg := mailgun.NewMailgun(apiKey)
	err := mg.SetAPIBase(mailgun.APIBaseEU)
	if err != nil {
		return nil, err
	}

	listIterator := mg.ListMailingLists(&mailgun.ListOptions{Limit: 100})

	var lists []mtypes.MailingList

	var page []mtypes.MailingList
	// The entire operation should not take longer than 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	for listIterator.Next(ctx, &page) {
		for _, list := range page {
			lists = append(lists, list)
		}
	}
	return lists, nil
}
