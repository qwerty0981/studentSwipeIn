package main

import (
	"bufio"
	"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
	"runtime"
	"os"
	"os/exec"
	"errors"
	flag "github.com/ogier/pflag"
)

func getEntryFrom(url, field string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", errors.New("Failed to GET the form")
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("Failed to read the response body")
	}

	input := string(contents)
	reg := *regexp.MustCompile(`(?mi)aria-label\=\"` + field + `\" [a-z0-9\.\'\"\-\= ]* name=\"(entry\.[0-9]*)\"`)

	match := reg.FindStringSubmatch(input)
	return match[1], nil
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
	// Setup Flags
	var verbosity *int = flag.IntP("logLevel", "v", 1, "Defines the verbosity of the logging.\n" +
					"\t\t1 - Critial Errors only (default)\n" +
					"\t\t2 - Includes Warnings\n" +
					"\t\t3 - Detailed logging\n" +
					"\t\t4 - Debug logging (May effect performance)")
	flag.Parse()
	// Finished loading flags

	clear()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please enter the url of the Google form you are linking to:")
	formURL, err := reader.ReadSting("\n")
	if err != nil {
		//Print failed
	}

	formURL = strings.TrimRight(formURL, "\n\r")

	fmt.Printf("Attempting to GET the form...")

	formData, err := getForm(formURL)

	fmt.Println("Received logging level of:", *verbosity)
}
