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
	if err := runApp(login, pass); err != nil {
		log.Fatal(err)
	}

}

func runApp(login, pass string) error {
	client, err := tentacle.NewGHAuthClient(login, pass)
	if err != nil {
		return err
	}

	user, _, err := client.Users.Get(context.TODO(), "")
	if err != nil {
		return err
	} else if user.Login == nil {
		return fmt.Errorf("cant determinate user login")
	}

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: reposPerPage},
	}

	for {
		repos, resp, err := client.Repositories.List(context.TODO(),
			*user.Login, opt)
		if err != nil {
			return err
		}

		for _, repo := range repos {
			fmt.Printf("%s\n", *repo.FullName)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil
}
