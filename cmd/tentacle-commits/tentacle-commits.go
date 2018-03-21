package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	tentacle "github.com/reinerRubin/gh-char-tentacle"
)

func main() {
	progName := os.Args[0]

	args := os.Args[1:]
	if len(args) != 3 {
		log.Fatalf("usage: %s login pass repo", progName)
	}

	login, pass, repo := args[0], args[1], args[2]

	client, err := tentacle.NewGHAuthClient(context.TODO(), login, pass)
	if err != nil {
		log.Fatal(err)
	}

	statSource, err := tentacle.NewGHCommitsStatSource(client, repo)
	if err != nil {
		log.Fatal(err)
	}

	stats := statSource.Source()
	result := make(chan tentacle.CharStat)

	go func() {
		result <- tentacle.NewStatMerger(stats).RunMergeStats()
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		statSource.Run()
	}()

	wg.Wait()

	sortedStat := (<-result).SortedStat()

	if len(sortedStat) == 0 {
		fmt.Println("something wrong :(")
		return
	}

	fmt.Print(sortedStat.TerminalGraph(terminalWidth()))
}

func terminalWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	tokens := strings.Split(strings.Trim(string(out), "\n"), " ")
	if len(tokens) != 2 {
		log.Fatal(err)
	}

	widthStr := tokens[1]
	width, err := strconv.Atoi(widthStr)
	if err != nil {
		log.Fatal(err)
	}

	return width
}
