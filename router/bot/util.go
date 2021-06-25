package bot

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

const (
	embURLRegexFragment = `/files/([\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12})`
)

var (
	verificationToken = os.Getenv("BOT_VERIFICATION_TOKEN")
	embURLRegex       = regexp.MustCompile(strings.ReplaceAll("https://q.trap.jp", ".", `\.`) + embURLRegexFragment)
)

func init() {
	if verificationToken == "" {
		panic(errors.New("verification token should not be empty."))
	}
}
