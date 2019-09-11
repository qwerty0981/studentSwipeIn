package delivery

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"sort"
)

var DEFAULT_FILENAME string = "cardSwipes.csv"

type CsvDriver struct {
	filename string
}

type CsvDriverInitializer struct {}

func (csv CsvDriver) Input(input map[string]string) error {
	// Check if file exists
	var file *os.File
	if !fileExists(csv.filename) {
		// Create a new .csv file
		
		// Generate header string
		header := make([]string, 0)
		for k, _ := range input {
			header = append(header, k)
		}

		sort.Strings(header)

		headerString := strings.Join(header, ",")

		// Create file
		f, err := os.Create(csv.filename)
		if err != nil {
			fmt.Println("Failed to create file: " + err.Error())
			panic(1)
		}

		// Write header
		f.WriteString(headerString + "\n")

		f.Sync()

		// Store the file handle for later writing
		file = f
	} else {
		// If the file exists open the handle
		fmt.Println("Opening existing file")
		f, err := os.OpenFile(csv.filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			fmt.Println("Failed to open file: " + err.Error())
			panic(1)
		}

		file = f
	}

	// Convert input data into csv line
	keyOrder := make([]string, 0)

	for k, _ := range input {
		keyOrder = append(keyOrder, k)
	}

	sort.Strings(keyOrder)

	data := ""
	for _, v := range keyOrder {
		data += input[v] + ","
	}

	data = strings.TrimRight(data, ",")

	_, err := file.WriteString(data + "\n")
	if err != nil {
		fmt.Println("Error writing to file: " + err.Error())
		panic(1)
	}

	file.Sync()

	file.Close()

	return nil
}

func (csv CsvDriverInitializer) Title() string {
	return "Local CSV"
}

func (csv CsvDriverInitializer) Configure() (Driver, error) {
	// Get the file that the user wants to write to

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("What file do you want to write to?[%s]\n", DEFAULT_FILENAME)

	filename, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed to read fron STDIN!")
		panic(1)
	}

	filename = strings.TrimRight(filename, "\n\r")
	
	if filename == "" {
		fmt.Printf("Using the default filename.")
		filename = DEFAULT_FILENAME
	}

	return CsvDriver{
		filename,
	}, nil
}

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}