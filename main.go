package main

import (
	"fmt"
	"mailinglist-backend-go/config"
)

func main() {
	fmt.Println(config.Value("MAILGUN_API_KEY"))
	fmt.Println(config.Values("MAILGUN_MAILING_LISTS"))
}
