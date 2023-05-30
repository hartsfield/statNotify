package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hartsfield/gmailer"
)

var (
	downDomains []string
	allDomains  []string = []string{
		"https://telesoft.network/",
		"https://tsconsulting.telesoft.network",
		"https://generic.telesoft.network",
		"https://sbvrt.telesoft.network",
		"https://particlestore.telesoft.network",
		"https://mysterygift.org",
		"https://btstrmr.xyz",
		"https://tagmachine.xyz",
	}

	lastNotifyDownStatus time.Time     = time.Now().AddDate(0, -1, 0)
	lastNotifyUpStatus   time.Time     = time.Now().AddDate(0, -1, 0)
	statusCheckRate      time.Duration = 60 * time.Second

	adminEmail string = "johnathanhartsfield@gmail.com"

	normalEmailRate    time.Duration = 24 * time.Hour
	normalEmailSubject string        = "🟢All systems nominal"
	normalEmailBody    string        = "Domains checked:\n" + strings.Join(allDomains, ", \n")

	alertEmailRate    time.Duration = 6 * time.Hour
	alertEmailSubject string        = "🔴CODE RED ‎"
	alertEmailBody    string        = "These domains didnt return status 200:\n" + strings.Join(downDomains, ", \n")
)

func main() {
	ticker := time.NewTicker(statusCheckRate)
	quit := make(chan struct{})
	checkStatus()
	for {
		select {
		case <-ticker.C:
			checkStatus()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func checkStatus() {
	log.Println("checking....")

	if len(downDomains) > 0 {
		notifyDown()
	}
	if len(downDomains) == 0 {
		notifyUp()
	}

	for _, domain := range allDomains {
		resp, err := http.Get(domain)
		if err != nil {
			log.Println(err)
		}

		index := firstIndex(downDomains, domain)
		if resp.StatusCode != 200 && index == -1 {
			downDomains = append(downDomains, domain)
		}
		if resp.StatusCode == 200 && index > -1 {
			downDomains = append(downDomains[:index], downDomains[index+1:]...)
		}
	}
}

func notifyDown() {
	if time.Now().Sub(lastNotifyDownStatus) > alertEmailRate && len(downDomains) > 0 {
		log.Println(alertEmailSubject)
		newMsg(adminEmail, alertEmailSubject, alertEmailBody).Send(func() {
			lastNotifyDownStatus = time.Now()
		})
	}
}

func notifyUp() {
	if time.Now().Sub(lastNotifyUpStatus) > normalEmailRate && len(downDomains) == 0 {
		log.Println(normalEmailSubject)
		newMsg(adminEmail, normalEmailSubject, normalEmailBody).Send(func() {
			lastNotifyUpStatus = time.Now()
		})
	}
}

func firstIndex(ss []string, match string) int {
	for k, s := range ss {
		if s == match {
			return k
		}
	}
	return -1
}

func newMsg(recipient, subject, body string) *gmailer.Message {
	return &gmailer.Message{
		Recipient: recipient,
		Subject:   subject,
		Body:      body,
	}
}
