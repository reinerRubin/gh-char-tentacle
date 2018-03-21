package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

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
		log.Fatalf("usage: %s username pass", progName)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	username, password := args[0], args[1]

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ready := make(chan struct{})

	go func() {
		defer func() {
			ready <- struct{}{}
		}()

		runApp(ctx, appOpt{
			username: username,
			password: password,
		})
	}()

	select {
	case <-ready:
	case <-c:
		// explicit cancel call
		cancel()
	}
}

type appOpt struct {
	username, password string
}

func runApp(ctx context.Context, appOpt appOpt) {
	client, err := tentacle.NewGHAuthClient(ctx, appOpt.username, appOpt.password)
	if err != nil {
		log.Fatal(err)
	}

	user, _, err := client.Users.Get(ctx, "")
	if err != nil || user.Login == nil {
		log.Fatal(err)
	}

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: reposPerPage},
	}

	for {
		repos, resp, err := client.Repositories.List(
			ctx,
			*user.Login, opt)
		if err != nil {
			panic(err)
		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		for _, repo := range repos {
			if repo.FullName == nil {
				continue
			}

			fmt.Printf("%s\n", *repo.FullName)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}
}
