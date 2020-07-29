package slack

import (
	"bytes"
	"net/http"

	"app/tools"
)

func SlackWebhook(channel, username, text, icon_emoji, webhook string) (err error) {
	jsonStr := `{"channel":"` + channel + `","username":"` + username + `","text":"` + text + `","icon_emoji":"` + icon_emoji + `"}`
	req, err := http.NewRequest(
		"POST",
		webhook,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	tools.FailOnError(err)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	tools.FailOnError(err)
	defer response.Body.Close()
	return nil
}