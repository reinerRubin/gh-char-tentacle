package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	tentacle "github.com/reinerRubin/gh-char-tentacle"
)

const (
	reposPerPage = 300
)

func main() {
	progName := os.Args[0]
	args := os.Args[1:]
	if len(args) != 2 {
		log.Fatalf("usage: %s login pass", progName)
	}

	login, pass := args[0], args[1]
	client, err := tentacle.NewGHClient(login, pass)
	if err != nil {
		log.Fatal(err)
	}

	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil || user.Login == nil {
		log.Fatal(err)
	}

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: reposPerPage},
	}

	for {
		repos, resp, err := client.Repositories.List(
			context.Background(),
			*user.Login, opt)
		if err != nil {
			panic(err)
		}

		for _, repo := range repos {
			fmt.Printf("%s\n", *repo.FullName)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}
}
