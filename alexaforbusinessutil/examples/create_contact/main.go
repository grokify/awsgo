package main

/*

https://aws.amazon.com/blogs/security/wheres-my-secret-access-key/
https://console.aws.amazon.com/iam/home

1. Create a user
2. Set `Access type` to `Programmatic access`
3. Add `AlexaForBusinessFullAccess` policy
4. Retrieve your `Access key ID` and `Secret access key`
5. Add credentials and `region = us-east-1` to `~/.aws/credentials`

*/

import (
	"fmt"
	"log"

	"github.com/grokify/mogo/fmt/fmtutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/alexaforbusiness"
)

type ContactInputBuilder struct {
	DisplayName string
	FirstName   string
	LastName    string
	PhoneNumber string
}

func (b *ContactInputBuilder) ToCreateContactInput() *alexaforbusiness.CreateContactInput {
	return &alexaforbusiness.CreateContactInput{
		DisplayName: &b.DisplayName,
		FirstName:   &b.FirstName,
		LastName:    &b.LastName,
		PhoneNumber: &b.PhoneNumber}
}

func main() {

	// Create Session
	mySession := session.Must(session.NewSession())

	// https://docs.aws.amazon.com/sdk-for-go/api/service/alexaforbusiness/#New
	svc := alexaforbusiness.New(mySession, aws.NewConfig().WithRegion("us-east-1"))

	if 1 == 0 {
		contact := ContactInputBuilder{
			DisplayName: "John Wang",
			FirstName:   "John",
			LastName:    "Wang",
			PhoneNumber: "+16505550100"}
		resp, err := svc.CreateContact(contact.ToCreateContactInput())
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(resp)
	}

	if 1 == 1 {
		req := alexaforbusiness.SearchContactsInput{}

		key := "DisplayName"
		name := ""
		filter := alexaforbusiness.Filter{
			Key:    &key,
			Values: []*string{}}
		filter.Values = append(filter.Values, &name)

		resp, err := svc.SearchContacts(&req)
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(resp)
	}

	fmt.Println("DONE")
}
