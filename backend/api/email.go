package api

import (
	"backend/db"
	"backend/types"
	"fmt"
	"net/smtp"
	"os"
	"sync"
)

type ChannelData struct {
	err   error
	email string
}

const smtpHost string = "smtp.gmail.com"
const smtpPort string = "587"
const senderEmail string = "antiqfurn.project@gmail.com"

var senderPassword string = os.Getenv("ANTIQ_FURN_PASS")

func sendEmail(
	c chan ChannelData,
	wg *sync.WaitGroup,
	subscriber types.User,
	auth smtp.Auth,
	listing types.FurnitureListing,
) {
	defer wg.Done()

	linkToListing := fmt.Sprintf("127.0.0.1:5173/market/%s", listing.ListingID.Hex())

	msg := []byte("To: " + subscriber.Email + "\r\n" +
		"Subject: New Furniture Listing" + "\r\n\r\n" +
		"A new furniture listng has been posted for the " + listing.Title +
		" a price of " + fmt.Sprintf("%.2f", listing.Cost) +
		fmt.Sprintf(". Click here to go to the listing: %s", linkToListing),
	)

	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		senderEmail,
		[]string{subscriber.Email},
		msg,
	)

	data := ChannelData{
		err:   err,
		email: subscriber.Email,
	}

	c <- data
}

func emailResults(c chan ChannelData) {
	for data := range c {
		if data.err != nil {
			fmt.Printf("Email result (%s): %s\n", data.email, data.err.Error())
		} else {
			fmt.Printf("Email result (%s): successfully sent email update\n", data.email)
		}
	}
}

/*
Sends an email update to all subscribed users of the recently listed furniture listing
*/
func SendNewListingNotificationEmail(listing types.FurnitureListing) error {
	subscribers, err := db.GetSubscribers()
	if err != nil {
		return err
	}

	// authenticate and connect to SMTP server for gmail
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	var wg sync.WaitGroup
	c := make(chan ChannelData)

	go emailResults(c)
	for i, subscriber := range subscribers {
		fmt.Printf("Subscriber #%d: %s\n", i+1, subscriber.UserID.Hex())

		wg.Add(1)
		go sendEmail(c, &wg, subscriber, auth, listing) // send each email asynchronously
	}

	wg.Wait()
	fmt.Printf("All emails have been sent!\n")

	return nil
}
