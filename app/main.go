package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/nlopes/slack"
)

type Tokens struct {
	Github string
	Slack string
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
	github_token, err := ioutil.ReadFile("config/github_access_token")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tokens.Github = string(github_token)

	// Githubアクセストークン読み込み
	slack_token, err := ioutil.ReadFile("config/slack_api_token")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tokens.Slack = string(slack_token)

	// 起動
	run(&tokens)

	os.Exit(0)
}
