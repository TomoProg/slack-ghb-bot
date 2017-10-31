package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/nlopes/slack"
)

type Tokens struct {
	Github string
	Slack  string
}

func run(tokens *Tokens) int {
	api := slack.New(tokens.Slack)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				log.Print("Hello Event")

			case *slack.MessageEvent:
				log.Printf("Message: %v\n", ev)
				rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", ev.Channel))

				// 認証
				ctx := context.Background()
				ts := oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: string(tokens.Github)},
				)
				tc := oauth2.NewClient(ctx, ts)

				client := github.NewClient(tc)

				// list all repositories for the authenticated user
				repos, _, err := client.Repositories.List(ctx, "", nil)
				if err != nil {
					fmt.Println("%v", err)
					os.Exit(1)
				}

				fmt.Println("%v", repos)

			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1

			}
		}
	}
}

func main() {
	tokens := Tokens{}

	// Githubアクセストークン読み込み
	tokens.Github = os.Getenv("GITHUB_TOKEN")

	// Slackアクセストークン読み込み
	tokens.Slack = os.Getenv("SLACK_TOKEN")

	// 起動
	run(&tokens)

	os.Exit(0)
}
