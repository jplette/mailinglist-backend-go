package mailinglist

import (
	"context"
	"mailinglist-backend-go/envcfg"
	"slices"
	"time"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

var domain = envcfg.Value("MAILGUN_DOMAIN")
var apiKey = envcfg.Value("MAILGUN_API_KEY")

type MGMailingList struct {
	*mtypes.MailingList
	Blocked bool `json:"blocked"`
	Hidden  bool `json:"hidden"`
}

func Lists(includeHidden bool) ([]MGMailingList, error) {
	mg := mailgun.NewMailgun(apiKey)
	err := mg.SetAPIBase(mailgun.APIBaseEU)
	if err != nil {
		return nil, err
	}

	listIterator := mg.ListMailingLists(&mailgun.ListOptions{Limit: 100})

	var lists []MGMailingList

	var page []mtypes.MailingList
	// The entire operation should not take longer than 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	for listIterator.Next(ctx, &page) {
		for _, list := range page {
			hideElement := isHidden(&list)
			if includeHidden == true || (includeHidden == false && hideElement == false) {
				lists = append(lists, MGMailingList{&list, isSubscriptable(&list), hideElement})
			}
		}
	}
	return lists, nil
}

func Subscribe(listaddress string, memberaddress string) error {
	mg := mailgun.NewMailgun(apiKey)
	err := mg.SetAPIBase(mailgun.APIBaseEU)
	if err != nil {
		return err
	}

	// The entire operation should not take longer than 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	subcribed := true

	err = mg.CreateMember(ctx, true, listaddress, mtypes.Member{Address: memberaddress, Subscribed: &subcribed})
	return err
}

func Unsubscribe(listaddress string, memberaddress string) error {
	mg := mailgun.NewMailgun(apiKey)
	err := mg.SetAPIBase(mailgun.APIBaseEU)
	if err != nil {
		return err
	}

	// The entire operation should not take longer than 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err = mg.DeleteMember(ctx, memberaddress, listaddress)

	return err
}

func isSubscriptable(list *mtypes.MailingList) bool {
	blocked := envcfg.Values("MAILGUN_BLOCKED_MAILING_LISTS")
	return !slices.Contains(blocked, list.Address)
}

func isHidden(list *mtypes.MailingList) bool {
	hidden := envcfg.Values("MAILGUN_HIDDEN_MAILING_LISTS")
	return slices.Contains(hidden, list.Address)
}
