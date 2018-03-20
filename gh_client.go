package tentacle

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
)

// NewGHClient TBD
func NewGHClient(username, password string) (*github.Client, error) {
	r := bufio.NewReader(os.Stdin)

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client := github.NewClient(tp.Client())
	ctx := context.Background()
	_, _, err := client.Users.Get(ctx, "")

	if _, ok := err.(*github.TwoFactorAuthError); ok {
		fmt.Print("GitHub OTP: ")
		otp, _ := r.ReadString('\n')
		tp.OTP = strings.TrimSpace(otp)
		_, _, err = client.Users.Get(ctx, "")
	}

	if err != nil {
		return nil, err
	}

	return client, nil
}