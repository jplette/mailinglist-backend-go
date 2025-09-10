package mailgun

import (
	"context"
	"mailinglist-backend-go/services/common"
	"mailinglist-backend-go/services/configReader"
	"slices"
	"time"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

var apiKey = configReader.Value("MAILGUN_API_KEY")

type MGMailingList struct {
	*mtypes.MailingList
	Blocked bool `json:"blocked"`
	Hidden  bool `json:"hidden"`
}

// APIMailingList is a simplified model used for API documentation to avoid
// referencing embedded external types that swag cannot inspect (mtypes.MailingList).
// This structure reflects the JSON the API returns for each mailing list entry.
//
// swagger:model APIMailingList
//
// Note: The runtime response encodes MGMailingList which embeds fields from
// mtypes.MailingList. The key fields commonly used are represented here.
// Keep this in sync with the actual JSON output of MGMailingList.
// If additional fields are needed in documentation, add them here accordingly.
// The presence of this type is solely to help Swagger generation.
//
// Fields examples are illustrative; adjust as needed.
//
//nolint:revive // exported for swagger docs
type APIMailingList struct {
	Address     string `json:"address" example:"news@example.com"`
	Name        string `json:"name,omitempty" example:"News"`
	Description string `json:"description,omitempty" example:"General news and updates"`
	Blocked     bool   `json:"blocked"`
	Hidden      bool   `json:"hidden"`
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
			hideElement := isHidden(list.Address)
			if includeHidden == true || (includeHidden == false && hideElement == false) {
				lists = append(lists, MGMailingList{&list, !isSubscriptable(list.Address), hideElement})
			}
		}
	}
	return lists, nil
}

func Subscribe(listAddress string, memberAddress string) error {
	mg := mailgun.NewMailgun(apiKey)
	err := mg.SetAPIBase(mailgun.APIBaseEU)
	if err != nil {
		return err
	}

	if isSubscriptable(listAddress) == false {
		return common.ErrForbidden
	}

	// The entire operation should not take longer than 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	subscribed := true

	err = mg.CreateMember(ctx, true, listAddress, mtypes.Member{Address: memberAddress, Subscribed: &subscribed})
	return err
}

func Unsubscribe(listAddress string, memberAddress string) error {
	mg := mailgun.NewMailgun(apiKey)
	err := mg.SetAPIBase(mailgun.APIBaseEU)
	if err != nil {
		return err
	}

	if isSubscriptable(listAddress) == false {
		return common.ErrForbidden
	}

	// The entire operation should not take longer than 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err = mg.DeleteMember(ctx, memberAddress, listAddress)

	return err
}

func isSubscriptable(list string) bool {
	blocked := configReader.Values("MAILGUN_BLOCKED_MAILING_LISTS")
	return !slices.Contains(blocked, list)
}

func isHidden(list string) bool {
	hidden := configReader.Values("MAILGUN_HIDDEN_MAILING_LISTS")
	return slices.Contains(hidden, list)
}
