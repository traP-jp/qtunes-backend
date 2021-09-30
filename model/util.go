package model

import (
	"net/url"
	"regexp"
)

const SoundChannelID = "8bd9e07a-2c6a-49e6-9961-4f88e83b4918"

var (
	BaseURL, _ = url.Parse("https://q.trap.jp/api/v3")
)

func RemoveExtensions(str string) string {
	rep := regexp.MustCompile(`\.[A-Za-z0-9]{3,5}$`)
	return rep.ReplaceAllString(str, "")
}

func ExtractFileIDs(text string) []string {
	urlRegexp := regexp.MustCompile(`https://q\.trap\.jp/files/([\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12})`)
	matches := urlRegexp.FindAllStringSubmatch(text, -1)

	res := make([]string, 0, len(matches))
	for _, v := range matches {
		res = append(res, v[1])
	}

	return res
}
