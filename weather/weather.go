package weather

import (
	"fmt"
	"io"
	"net/http"
)

func GetWeather(city string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://wttr.in/%s?format=3", city))
	if err != nil {
		return "", fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid http status: %d", resp.StatusCode)
	}
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return string(result), nil
}

var (
	Cities = []string{
		"東京",
		"ソウル",
		"北京",
		"シドニー",
		"パリ",
		"ロンドン",
		"ベルリン",
		"ニューヨーク",
		"ロサンゼルス",
	}
)
