package main

import (
	"bufio"
	"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"runtime"
	"os"
	"os/exec"
	"errors"
	"github.com/libpoly/libcardutils-go"
	"github.com/qwerty0981/formAssistant"
	flag "github.com/ogier/pflag"
)

func getEntryFromForm(formData, field string) string {
	reg := *regexp.MustCompile(`(?mi)aria-label\=\"` + field + `\" [a-z0-9\.\'\"\-\= ]* name=\"(entry\.[0-9]*)\"`)

	match := reg.FindStringSubmatch(formData)
	if match == nil {
		return ""
	}
	return match[1]
}

func getForm(url string) (string, error) {
	response, err := http.Get(url)

	if err != nil {
		return "", errors.New("Failed to GET the form")
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
         fmt.Println("Error Clearing Screen");
      }
   } else if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
      c := exec.Command("clear");
      c.Stdout = os.Stdout
      if err := c.Run(); err != nil {
         fmt.Println("Error Clearing Screen");
      }
   }
}

func main() {
	supportedEntries := map[string][]string{
		"First Name": {"first name", "firstName"},
		"Last Name": {"last name", "lastName"},
		"Student ID": {"student id", "studentID", "student number", "studentNumber"},
		"Email": {"email", "studentEmail", "student Email"},
	}


	// Setup Flags
	var verbosity *int = flag.IntP("logLevel", "v", 1, "Defines the verbosity of the logging.\n" +
					"\t\t1 - Critial Errors only (default)\n" +
					"\t\t2 - Includes Warnings\n" +
					"\t\t3 - Detailed logging\n" +
					"\t\t4 - Debug logging (May effect performance)")
	flag.Parse()
	// Finished loading flags
	fmt.Println("Verbosity:",verbosity)
	clear()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please enter the url of the Google form you are linking to:")
	formURL, err := reader.ReadString('\n')
	if err != nil {
		//Print failed
	}

	formURL = strings.TrimRight(formURL, "\n\r")
	gForm := form.GoogleForm(strings.TrimSuffix(formURL, "viewform") + "formResponse")

	fmt.Printf("Attempting to GET the form...")

	formData, err := getForm(formURL)
	if err != nil {
		fmt.Printf("Failed!\n")
		fmt.Println("\t-> ", err)
		// Handle Error		
	}

	fmt.Printf("Done\n")

	fmt.Printf("Attempting to find entry IDs...\n")

	for k, v := range supportedEntries {
		fmt.Println("  Locating " + k + "...")

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

	fmt.Println("The form has been configured")
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

		fmt.Println(card)
	}
}
