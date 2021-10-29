// Contains functions that set up the email struct.
package src

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Email struct {
	To      string
	Cc      string
	Subject string
	Message string
}

// Takes the arguments needed to create an email, and turns them into an Email struct
func CreateEmail(ContentFile string, RecipientsFile string, Cc string) Email {
	var email Email

	// get recipients list
	recipients_file, read_recipients_err := os.Open(RecipientsFile)
	content_file, read_content_err := os.Open(ContentFile)
	if read_recipients_err != nil {
		fmt.Println("Recipient file opening error: ", read_recipients_err)
	}
	if read_content_err != nil {
		fmt.Println("Content file opening error: ", read_content_err)
	}
	defer recipients_file.Close()
	defer content_file.Close()

	// create empty list to hold recipients
	recipientsSlice := make([]string, 0, 999)

	// add mailing list to slice
	recipients_scanner := bufio.NewScanner(recipients_file)
	for recipients_scanner.Scan() {
		recipientsSlice = append(recipientsSlice, recipients_scanner.Text())
	}

	// convert mailing list to string and add to `To` field
	// also add Cc if available
	email.To = strings.Join(recipientsSlice, ", ")
	email.Cc = Cc

	// create empty slice to hold content
	contentSlice := make([]string, 0, 9999)

	// read text file and get message and subject
	content_scanner := bufio.NewScanner(content_file)
	for content_scanner.Scan() {
		contentSlice = append(contentSlice, content_scanner.Text())
	}
	// set first line of email as subject
	email.Subject = striphtml(contentSlice[0])
	email.Message = strings.Join(contentSlice, "\r\n")

	return email
}

// Strips HTML from a string.
func striphtml(in string) string {
	// regex to match html tag
	const pattern = `(<\/?[a-zA-A]+?[^>]*\/?>)*`
	r := regexp.MustCompile(pattern)
	groups := r.FindAllString(in, -1)
	for _, group := range groups {
		if strings.TrimSpace(group) != "" {
			in = strings.ReplaceAll(in, group, "")
		}
	}
	return in
}