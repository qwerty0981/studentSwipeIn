package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"strconv"
	"github.com/libpoly/libcardutils-go"
	"./delivery"
)

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
	//errText, infoText := "[ERROR]:", "[INFO] :"

	// Do delivery selection and setup here
	initializers := []delivery.DriverInitializer{
		delivery.GoogleForm{},
		delivery.CsvDriverInitializer{},
	}
	reader := bufio.NewReader(os.Stdin)

	var option int

	for {
		clear()

		fmt.Println("How would you like to output the card swipes?\n")

		for i, init := range initializers {
			fmt.Printf("%d) %s\n", i+1, init.Title())
		}
		fmt.Println("\nPlease select an option:")

		textOption, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Failed to read string:", err)
			time.Sleep(time.Second)
			os.Exit(1)
		}

		option, err = strconv.Atoi(strings.TrimRight(textOption, "\n\r"))

		if err == nil {
			if option >= 1 && option <= len(initializers) {
				break;
			}
		}

		fmt.Println(err)

		fmt.Println("Error invalid item selected")
		time.Sleep(2 * time.Second)
	}

	clear()

	driver, err := initializers[option-1].Configure()

	if err != nil {
		fmt.Println("Invalid configuration settings.")
		os.Exit(1)
	}

	fmt.Println("\n\nPlease press enter to start accepting card swipes.")
	reader.ReadString('\n')

	card := cardutils.New()

	CardData := make(map[string]string)
	for {
		clear()
		fmt.Println("Please enter the section that your are in or type 'quit' to exit:")
		section, _ := reader.ReadString('\n')

		section = strings.TrimRight(section, "\n\r")

		if section == "exit" || section == "quit" {
			fmt.Println("Quiting...")
			break
		}

		fmt.Println("Please swipe a card:")
		scanIn, _ := reader.ReadString('\n')

		scanIn = strings.TrimRight(scanIn, "\n\r")



		card.Swipe(scanIn)

		CardData["First Name"] = card.GetFirstName()
		CardData["Last Name"] = card.GetLastName()
		CardData["Email"] = makeEmail(card.GetFirstName(), card.GetLastName(), card.GetID())
		CardData["Student ID"] = card.GetID()
		CardData["Section"] = section

		driver.Input(CardData)

		fmt.Println("Submitted successfully!")
		time.Sleep(1 * time.Second)
	}
}
