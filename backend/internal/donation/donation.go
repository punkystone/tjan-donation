package donation

import (
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/rs/zerolog/log"
)

var donationPattern = regexp.MustCompile(`([^\s]+) spendet €([^\s]+) tjanL Danke!`)

type Donation struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type Service struct {
	TopDonationChannel        chan Donation
	chat                      *twitch.Client
	receiveChannel            string
	latestTopDonation         *Donation
	latestTopDonationSavePath string
}

func NewDonationService(channel string, receiveChannel string, latestTopDonationSavePath string) *Service {
	chat := twitch.NewAnonymousClient()
	chat.TLS = true
	chat.OnConnect(func() {
		log.Info().Msg("Connected to Twitch!")
	})

	chat.Join(channel)

	return &Service{
		TopDonationChannel:        make(chan Donation),
		chat:                      chat,
		receiveChannel:            receiveChannel,
		latestTopDonation:         nil,
		latestTopDonationSavePath: latestTopDonationSavePath,
	}
}

const reconnectInterval = 10 * time.Second

func (donationService *Service) Start() {
	donation, err := donationService.getLatestTopDonation()
	if err != nil {
		log.Error().Err(err)
	} else {
		donationService.latestTopDonation = donation
		donationService.TopDonationChannel <- *donation
	}

	donationService.chat.OnPrivateMessage(func(message twitch.PrivateMessage) {
		donationService.handleMessage(message)
	})
	for {
		err := donationService.chat.Connect()
		if err != nil {
			log.Error().Err(err)
		}
		time.Sleep(reconnectInterval)
	}
}

func (donationService *Service) handleMessage(message twitch.PrivateMessage) {
	if message.User.Name != "punkystone" {
		return
	}
	matches := donationPattern.FindStringSubmatch(message.Message)
	const requiredGroups = 3
	if len(matches) != requiredGroups {
		return
	}
	donationAmount, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse donation amount")
		return
	}
	if donationService.latestTopDonation != nil && donationAmount <= donationService.latestTopDonation.Amount {
		return
	}
	donation := Donation{matches[1], donationAmount}
	donationService.latestTopDonation = &donation
	donationService.TopDonationChannel <- donation
	err = donationService.saveLatestTopDonation(donation)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save top donation")
	}
}

func (donationService *Service) GetLatestTopDonation() *Donation {
	return donationService.latestTopDonation
}

func (donationService *Service) getLatestTopDonation() (*Donation, error) {
	file, err := os.OpenFile(donationService.latestTopDonationSavePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	donation := &Donation{}
	err = json.NewDecoder(file).Decode(donation)
	if err != nil {
		return nil, err
	}
	return donation, nil
}

func (donationService *Service) saveLatestTopDonation(donation Donation) error {
	file, err := os.OpenFile(donationService.latestTopDonationSavePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	err = enc.Encode(donation)
	if err != nil {
		return err
	}
	return nil
}
