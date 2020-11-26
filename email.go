package jitter

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	gomail "gopkg.in/mail.v2"
)

//EmailPackage struct that holds all the required info coming on from the contact form
type EmailPackage struct {
	Name       string
	Sender     string // email address
	Subject    string
	Message    string
	CarbonCopy bool
}

var contactEmail = "grantjliu21@gmail.com"

//ContactEmailHandler Entry point for Contact email
func ContactEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var eP EmailPackage
	err := json.NewDecoder(r.Body).Decode(&eP)
	fmt.Printf("%+v\n", eP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isEmailValid(eP.Sender) {
		http.Error(w, errors.New("Email is invalid. This message should not be seen as emails are checked clientside").Error(), 422)
		return
	}

	if !SendContactMail(eP) {
		http.Error(w, errors.New("Error sending email").Error(), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func isEmailValid(e string) bool {
	var validEmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(e) < 3 || len(e) > 254 {
		return false
	}
	return validEmailRegex.MatchString(e)
}

//SendContactMail Short Comment for now
func SendContactMail(pack EmailPackage) bool {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "NerveNorms@gmail.com")

	// Set E-Mail receivers
	if pack.CarbonCopy {
		m.SetHeader("To", contactEmail, pack.Sender)
	} else {
		m.SetHeader("To", contactEmail)
	}

	// Set E-Mail subject
	m.SetHeader("Subject", "[NerveNorms] "+pack.Subject)

	// Set reply to address
	m.SetAddressHeader("Reply-To", pack.Sender, pack.Name)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", "From: "+pack.Name+"("+pack.Sender+")\n\n"+pack.Message)

	// Settings for SMTP server
	secret := "projects/NerveNorms/secrets/automated-password/versions/latest"
	password, err := getPassword(secret)
	if err == nil {
		panic(err)
	}
	d := gomail.NewDialer("smtp.gmail.com", 587, "NerveNorms@gmail.com", password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return true
}

func getPassword(name string) (string, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	return string(result.Payload.Data), nil
}
