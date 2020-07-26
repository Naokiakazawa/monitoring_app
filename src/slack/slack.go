package slack

import (
	"bytes"
	"log"
	"net/http"
)

func SlackWebhook(channel, username, text, icon_emoji, webhook string) (err error) {
	jsonStr := `{"channel":"` + channel + `","username":"` + username + `","text":"` + text + `","icon_emoji":"` + icon_emoji + `"}`
	req, err := http.NewRequest(
		"POST",
		webhook,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil{
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil{
		log.Fatal(err)
	}
	defer response.Body.Close()
	return nil
}