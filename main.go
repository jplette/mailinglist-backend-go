package main

import (
	"fmt"
	"mailinglist-backend-go/config"
	"mailinglist-backend-go/mailgun"
)

func main() {
	fmt.Println(config.Value("MAILGUN_API_KEY"))
	fmt.Println(mailgun.MailingLists())
}
