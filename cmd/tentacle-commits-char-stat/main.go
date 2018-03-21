package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	tentacle "github.com/reinerRubin/gh-char-tentacle"
)

func main() {
	progName := os.Args[0]

	args := os.Args[1:]
	if len(args) != 3 {
		log.Fatalf("usage: %s login pass repo", progName)
	}

	login, pass, repo := args[0], args[1], args[2]
	if err := runApp(login, pass, repo); err != nil {
		log.Fatal(err)
	}
}

func runApp(login, pass, repo string) error {
	client, err := tentacle.NewGHAuthClient(login, pass)
	if err != nil {
		return err
	}

	statSource, err := tentacle.NewGHCommitsStatSource(client, repo)
	if err != nil {
		return err
	}

	mergedStat := make(chan tentacle.CharStat)
	go func() {
		mergedStat <- tentacle.NewStatMerger(statSource.StatChannel()).Run()
	}()

	if err := statSource.Run(); err != nil {
		return err
	}

	sortedStat := (<-mergedStat).SortedStat()

	if len(sortedStat) == 0 {
		return fmt.Errorf("cant load stat for %s", repo)
	}

	width, err := terminalWidth()
	if err != nil {
		return nil
	}
	fmt.Print(sortedStat.TextGraph(width))

	return nil
}

func terminalWidth() (int, error) {
	// idk about windows
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()

	if err != nil {
		return 0, err
	}

	tokens := strings.Split(strings.Trim(string(out), "\n"), " ")
	if len(tokens) != 2 {
		return 0, fmt.Errorf("cant determinate console size")
	}

	widthStr := tokens[1]
	width, err := strconv.Atoi(widthStr)
	if err != nil {
		return 0, err
	}

	return width, nil
}
