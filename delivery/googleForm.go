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
	"github.com/qwerty0981/FormAssistant"
)

type GoogleForm struct{}

type googleFormDriver struct {
	googleForm form.Form
}

func (gfd googleFormDriver) Input(data map[string]string) error {	
   formData, _ := getForm(gfd.googleForm.GetURL())

   values := getHiddenValuesFromForm(formData)

   for k,v := range values {
      data[k] = v
   }

   for k, _ := range data {
		if _, ok := gfd.googleForm.Endpoints()[k]; !ok {
			delete(data, k)
		}
	}

   err := gfd.googleForm.Post(data)
   if err != nil {
      fmt.Println(err)
   }

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

   values := getHiddenValuesFromForm(formData)

   for k, _ := range values {
      gForm.AddEndpoint(k, k)
   }
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

func getHiddenValuesFromForm(formData string) map[string]string {
   reg := *regexp.MustCompile(`<input type="hidden" name="([a-zA-Z0-9]*)" value="([\[\]a-z\,0-9\-\&\;\n]*)">`)

   matches := reg.FindAllStringSubmatch(formData, -1)
   if matches == nil {
      fmt.Printf("Failed to find any hidden values!")
      panic(1)
   }

   results := make(map[string]string)

   for _, match := range matches {
      results[match[1]] = strings.ReplaceAll(match[2], "&quot;", "\"")
   }

   return results
}

func cleanURL(url string) string {
	return strings.TrimSuffix(url, "?usp=sf_link")
}
