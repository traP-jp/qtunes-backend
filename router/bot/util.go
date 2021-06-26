package bot

import (
	"errors"
	"os"
	"regexp"
)

var (
	verificationToken = os.Getenv("BOT_VERIFICATION_TOKEN")
)

func init() {
	if verificationToken == "" {
		panic(errors.New("verification token should not be empty."))
	}
}

func removeExt(str string) string {
	rep := regexp.MustCompile(`\.[A-Za-z0-9]{3,5}$`)
	return rep.ReplaceAllString(str, "")
}

func extractFileIDs(text string) []string {
	urlRegexp := regexp.MustCompile(`https://q\.trap\.jp/files/([\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12})`)
	matches := urlRegexp.FindAllStringSubmatch(text, -1)

	res := make([]string, 0, len(matches))
	for _, v := range matches {
		res = append(res, v[1])
	}

	return res
}
