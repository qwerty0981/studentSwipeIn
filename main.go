package main

import (
	"bufio"
	"net/http"
	"fmt"
	"time"
	"io/ioutil"
	"regexp"
	"strings"
	"runtime"
	"os"
	"os/exec"
	"errors"
	"github.com/libpoly/libcardutils-go"
	"github.com/qwerty0981/formAssistant"
)

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

func clear() {
   if runtime.GOOS == "windows"{
      c := exec.Command("cmd", "/c", "cls");
      c.Stdout = os.Stdout
      if err := c.Run(); err != nil {
         fmt.Println("Failed to clear screen");
      }
   } else if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
      c := exec.Command("clear");
      c.Stdout = os.Stdout
      if err := c.Run(); err != nil {
         fmt.Println("Failed to clear screen");
      }
   }
}

func makeEmail(fname, lname, num string) string {
	return strings.ToLower(string(fname[0]) + lname + num + "@floridapoly.edu")
}

func main() {
	supportedEntries := map[string][]string{
		"First Name": {"first name", "firstName"},
		"Last Name": {"last name", "lastName"},
		"Student ID": {"student id", "id", "studentID", "student number", "studentNumber"},
		"Email": {"email", "studentEmail", "student Email"},
	}

	errText, infoText := "[ERROR]:", "[INFO] :"

	clear()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please enter the url of the Google form you are linking to:")
	formURL, err := reader.ReadString('\n')
	if err != nil {
		//Print failed
	}

	formURL = strings.TrimRight(formURL, "\n\r")
	formURL = cleanURL(formURL)
	gForm := form.GoogleForm(strings.TrimSuffix(formURL, "viewform") + "formResponse")

	fmt.Printf("%s Attempting to GET the form...", infoText)

	formData, err := getForm(formURL)
	if err != nil {
		fmt.Printf("Failed!\n")
		fmt.Println(errText, err)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	fmt.Printf("Done\n")

	fmt.Println(infoText, "Attempting to find entry IDs...")

	for k, v := range supportedEntries {
		fmt.Println(infoText, "  Locating " + k + "...")

		for _, pattern := range v {
			fmt.Printf("%s \tTrying %s: ", infoText, pattern)
			match := getEntryFromForm(formData, pattern)
			if match != "" {
				fmt.Printf("Found!\n")
				fmt.Printf("%s \tAdding %s as %s entry\n", infoText, match, k)
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

	fmt.Println("Please press enter to start accepting card scans")
	reader.ReadString('\n')

	card := cardutils.New()

	for {
		clear()
		fmt.Println("Please swipe a card:")
		scanIn, _ := reader.ReadString('\n')

		scanIn = strings.TrimRight(scanIn, "\n\r")

		if scanIn == "exit" || scanIn == "quit" {
			fmt.Println("Quiting...")
			break
		}

		card.Swipe(scanIn)
		postData := make(map[string]string)

		for k, _ := range gForm.Endpoints() {
			switch(k){
			case "First Name":
				postData["First Name"] = card.GetFirstName()
			case "Last Name":
				postData["Last Name"] = card.GetLastName()
			case "Email":
				postData["Email"] = makeEmail(card.GetFirstName(), card.GetLastName(), card.GetID())
			case "Student ID":
				postData["Student ID"] = card.GetID()
			default:
				fmt.Println("Unhandled Endpoint:", k)
			}
		}

		gForm.Post(postData)

		fmt.Println("Submitted successfully!")
		time.Sleep(1 * time.Second)
	}
}
