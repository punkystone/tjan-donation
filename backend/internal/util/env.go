package util

import (
	"errors"
	"os"
)

type Env struct {
	ReceiveChannel            string
	Channel                   string
	LatestTopDonationSavePath string
}

func CheckEnv() (*Env, error) {
	receiveChannel, exists := os.LookupEnv("RECEIVE_CHANNEL")
	if !exists {
		return nil, errors.New("ACCESS_TOKEN environment variable not set")
	}
	channel, exists := os.LookupEnv("CHANNEL")
	if !exists {
		return nil, errors.New("REFRESH_TOKEN environment variable not set")
	}
	latestTopDonationSavePath, exists := os.LookupEnv("LATEST_TOP_DONATION_SAVE_PATH")
	if !exists {
		return nil, errors.New("LATEST_TOP_DONATION_SAVE_PATH environment variable not set")
	}

	env := &Env{
		ReceiveChannel:            receiveChannel,
		Channel:                   channel,
		LatestTopDonationSavePath: latestTopDonationSavePath,
	}
	return env, nil
}
