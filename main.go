package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hartsfield/gmailer"
)

var (
	// configuration
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

	logFilePath = os.Getenv("statLogPath")

	lastNotifyDownStatus time.Time     = time.Now().AddDate(0, -1, 0)
	lastNotifyUpStatus   time.Time     = time.Now().AddDate(0, -1, 0)
	statusCheckRate      time.Duration = 60 * time.Second

	startUpText = `statNotify inititated... checking for service outages` +
		` every ` + statusCheckRate.Abs().String()

	adminEmail string = os.Getenv("statAdminEmail")

	normalEmailRate    time.Duration = 24 * time.Hour
	normalEmailSubject string        = ";^) All systems nominal"

	alertEmailRate    time.Duration = 6 * time.Hour
	alertEmailSubject string        = "! :^0 CODE RED!"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Admin Email:", adminEmail)
	log.Println(startUpText)

	printConsoleMessage()

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

	if len(downDomains) > 0 {
		notifyDown()
	}
	if len(downDomains) == 0 {
		notifyUp()
	}
}

func notifyDown() {
	if time.Now().Sub(lastNotifyDownStatus) > alertEmailRate && len(downDomains) > 0 {
		alertEmailBody := "These domains didnt return status 200:\n" + strings.Join(downDomains, ", \n")
		log.Println(strings.ReplaceAll(alertEmailBody, "\n", " "))
		newMsg(adminEmail, alertEmailSubject, alertEmailBody).Send(func() {
			lastNotifyDownStatus = time.Now()
		})
	}
}

func notifyUp() {
	if time.Now().Sub(lastNotifyUpStatus) > normalEmailRate && len(downDomains) == 0 {
		normalEmailBody := "Domains checked:\n" + strings.Join(allDomains, ", \n")
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

func printConsoleMessage() {
	fmt.Println("Today is: ", time.Now().Format(time.UnixDate))
	fmt.Println(startUpText)
	fmt.Println("log file:", logFilePath)
	fmt.Println("Admin Email:", adminEmail)
}
