package main

import (
	"context"
	"log"
	"os"
	"regexp"

	l "github.com/aws/aws-lambda-go/lambda"
	resty "github.com/go-resty/resty/v2"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/walkersumida/sls-slack-event-subscriber-template/event/action/strings/ja"
)

type DeeplResponse struct {
	Translations []struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	} `json:"translations"`
}

func Handler(ctx context.Context, input *slackevents.AppMentionEvent) {
	cli := slack.New(os.Getenv("SLACK_ACCESS_TOKEN"))
	text := removeMention(input.Text)
	isJa := ja.IsJA(text)
	targetLang := "JA"
	if isJa {
		targetLang = "EN"
	}

	client := resty.New()
	resp, err := trans(client, text, targetLang)
	if err != nil {
		log.Printf("failed to trans: %s", err)
	}
	result := resp.Result().(*DeeplResponse)

	_, _, err = cli.PostMessage(
		input.Channel,
		slack.MsgOptionText(
			result.Translations[0].Text, false,
		),
	)
	if err != nil {
		log.Printf("failed to post message: %s", err)
	}
}

func removeMention(msg string) string {
	re := regexp.MustCompile("<@.+?>")
	return re.ReplaceAllString(msg, "")
}

func trans(client *resty.Client, text string, targetLang string) (*resty.Response, error) {
	resp, err := client.
		SetHeader("Authorization", "DeepL-Auth-Key "+os.Getenv("DEEPL_API_KEY")).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetContentLength(true).
		SetFormData(map[string]string{
			"text":        text,
			"target_lang": targetLang,
		}).
		R().
		SetResult(&DeeplResponse{}).
		Post("https://api-free.deepl.com/v2/translate")
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func main() {
	l.Start(Handler)
}
