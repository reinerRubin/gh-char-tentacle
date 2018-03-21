package tentacle

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
)

const (
	commitsPerPage = 400
)

// GHCommitsStatSource TBD
type GHCommitsStatSource struct {
	sinkStat  chan CharStat
	GHClient  *github.Client
	repoOwner string
	repoName  string
}

// NewGHCommitsStatSource TBD
func NewGHCommitsStatSource(client *github.Client, repoFullName string) (*GHCommitsStatSource, error) {
	tokens := strings.Split(repoFullName, "/")
	if len(tokens) != 2 {
		return nil, fmt.Errorf("repo name format: owner/repoName but %s", repoFullName)
	}

	return &GHCommitsStatSource{
		sinkStat:  make(chan CharStat),
		GHClient:  client,
		repoOwner: tokens[0],
		repoName:  tokens[1],
	}, nil
}

// StatChannel TBD
func (ss *GHCommitsStatSource) StatChannel() <-chan CharStat {
	return ss.sinkStat
}

// Run TBD
func (ss *GHCommitsStatSource) Run() error {
	defer close(ss.sinkStat)
	opt := &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: commitsPerPage},
	}

	for {
		commits, resp, err := ss.GHClient.Repositories.ListCommits(
			context.Background(),
			ss.repoOwner, ss.repoName, opt)
		if err != nil {
			return err
		}

		for _, commit := range commits {
			if commit.Commit == nil || commit.Commit.Message == nil {
				continue
			}

			ss.sinkStat <- NewCharStat(*commit.Commit.Message)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil
}
