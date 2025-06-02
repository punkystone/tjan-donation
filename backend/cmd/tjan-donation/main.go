package main

import (
	"tjan-donation/internal/donation"
	"tjan-donation/internal/server"
	"tjan-donation/internal/socket"
	"tjan-donation/internal/util"
)

func main() {
	env, err := util.CheckEnv()
	if err != nil {
		panic(err)
	}
	donationService := donation.NewDonationService(env.Channel, env.ReceiveChannel, env.LatestTopDonationSavePath)
	go donationService.Start()
	hub := socket.NewHub(donationService)
	go hub.Run()
	err = server.StartServer(hub)
	if err != nil {
		panic(err)
	}
}
