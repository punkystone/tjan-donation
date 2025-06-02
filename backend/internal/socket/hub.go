package socket

import (
	"fmt"

	"strconv"
	"strings"
	"tjan-donation/internal/donation"
)

type Hub struct {
	clients         map[*Client]bool
	Register        chan *Client
	unregister      chan *Client
	donationService *donation.Service
}

func NewHub(donationService *donation.Service) *Hub {
	return &Hub{
		Register:        make(chan *Client),
		unregister:      make(chan *Client),
		clients:         make(map[*Client]bool),
		donationService: donationService,
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.clients[client] = true
			donation := hub.donationService.GetLatestTopDonation()
			if donation == nil {
				client.Send <- []byte("- €")
			} else {
				client.Send <- fmt.Appendf(nil, "%s - %s€", donation.Name, formatFloat(donation.Amount))
			}
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.Send)
			}
		case donation := <-hub.donationService.TopDonationChannel:
			for client := range hub.clients {
				if donation == nil {
					client.Send <- []byte("- €")
				} else {
					client.Send <- fmt.Appendf(nil, "%s - %s€", donation.Name, formatFloat(donation.Amount))
				}
			}
		}
	}
}

func formatFloat(value float64) string {
	rounded := strconv.FormatFloat(value, 'f', 2, 64)
	if strings.HasSuffix(rounded, ".00") {
		return strings.Split(rounded, ".")[0]
	}
	return rounded
}
