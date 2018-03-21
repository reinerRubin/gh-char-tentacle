package tentacle

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
)

// NewGHAuthClient TBD
func NewGHAuthClient(username, password string) (*github.Client, error) {
	r := bufio.NewReader(os.Stdin)

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client := github.NewClient(tp.Client())
	_, _, err := client.Users.Get(context.TODO(), "")

	if _, ok := err.(*github.TwoFactorAuthError); ok {
		fmt.Print("GitHub OTP: ")
		otp, _ := r.ReadString('\n')
		tp.OTP = strings.TrimSpace(otp)
		_, _, err = client.Users.Get(context.TODO(), "")
	}

	if err != nil {
		return nil, err
	}
	return client, nil
}
