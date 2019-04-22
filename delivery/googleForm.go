package delivery

import (
	"bufio"
	"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"os"
	"errors"
	"time"
	"github.com/qwerty0981/formAssistant"
)

type GoogleForm struct{
}

type googleFormDriver struct {
	googleForm form.Form
}

func (gfd googleFormDriver) Input(data map[string]string) error {
	for k, _ := range gfd.googleForm.Endpoints() {
		if _, ok := data[k]; !ok {
			delete(data, k)
		}
	}

	gfd.googleForm.Post(data)

	return nil
}

func (gf GoogleForm) Title() string {
	return "Google Form"
}

func (gf GoogleForm) Configure() (Driver, error) {
	supportedEntries := map[string][]string{
		"First Name": {"first name", "firstName"},
		"Last Name": {"last name", "lastName"},
		"Student ID": {"student id", "id", "studentID", "student number", "studentNumber"},
		"Email": {"email", "studentEmail", "student Email"},
      "PhoneNumber": {"number", "phone", "Phone Number", "phone number", "Phone number"},
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please enter the url of the Google form you are linking to:")
	formURL, err := reader.ReadString('\n')
	if err != nil {
		//Print failed
	}

	formURL = strings.TrimRight(formURL, "\n\r")
	formURL = cleanURL(formURL)
	gForm := form.GoogleForm(strings.TrimSuffix(formURL, "viewform") + "formResponse")

	fmt.Printf("Getting form...")

	formData, err := getForm(formURL)
	if err != nil {
		fmt.Printf("Failed!")
		time.Sleep(1 * time.Second)
		return googleFormDriver{}, err
	}

	fmt.Printf("Done\n")

	fmt.Println("Attempting to find entry IDs...")

	for k, v := range supportedEntries {
		fmt.Println("Locating " + k + "...")

		for _, pattern := range v {
			fmt.Printf("\tTrying %s: ", pattern)
			match := getEntryFromForm(formData, pattern)
			if match != "" {
				fmt.Printf("Found!\n")
				fmt.Printf("\tAdding %s as %s entry\n", match, k)
				gForm.AddEndpoint(k, match)
				break
			} else {
				fmt.Printf("Failed\n")
			}
		}
	}

	fmt.Printf("Form configured with: ")
	var endpoints string
	for endpoint, _ := range gForm.Endpoints() {
		endpoints += endpoint + ", "
	}
	endpoints = strings.TrimSuffix(endpoints, ", ") + "\n"
	fmt.Printf("%s", endpoints)

	return googleFormDriver{gForm}, nil
	// gForm has all endpoints at this point
}



func getForm(url string) (string, error) {
	response, err := http.Get(url)

	if err != nil {
		return "", errors.New("Invalid URL? Make sure a google form was supplied")
	}

	formContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("Failed to read bytes from form GET response")
	}

	return string(formContent), nil
}

func getEntryFromForm(formData, field string) string {
	reg := *regexp.MustCompile(`(?mi)aria-label\=\"` + field + `\" [a-z0-9\.\'\"\-\= ]* name=\"(entry\.[0-9]*)\"`)

	match := reg.FindStringSubmatch(formData)
	if match == nil {
		return ""
	}
	return match[1]
}

func cleanURL(url string) string {
	return strings.TrimSuffix(url, "?usp=sf_link")
}
