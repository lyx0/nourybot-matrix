package commands

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/shkh/lastfm-go/lastfm"
)

func LastFmUserRecent(user string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", ErrInternalServerError
	}

	apiKey := os.Getenv("LAST_FM_API_KEY")
	apiSecret := os.Getenv("LAST_FM_SECRET")

	api := lastfm.New(apiKey, apiSecret)
	result, _ := api.User.GetRecentTracks(lastfm.P{"user": user}) //discarding error

	var reply string
	for i, track := range result.Tracks {
		// The 0th result is the most recent one since it goes from most recent
		// to least recent.
		if i == 0 {
			reply = fmt.Sprintf("Most recently played track for user %v: %v - %v", user, track.Artist.Name, track.Name)
			return reply, nil
		}
	}
	return "", ErrInternalServerError
}
