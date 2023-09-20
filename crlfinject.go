// crlfinject.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
)

func main() {

	red := color.New(color.FgRed)
	green := color.New(color.FgGreen).Add(color.Underline)

	onlyPositive := flag.Bool("p", false, "show only positive results.")
	payloadList := "payloads.txt"
	website := flag.String("site", "", "site to test.")
	flag.Parse()

	file, err := os.Open(payloadList)
	if err != nil {
		fmt.Println("Error : ", err)
		os.Exit(1)
	}
	defer file.Close()

	if *website == "" {
		fmt.Printf("choose a website to check\n")
		os.Exit(1)
	}
	response, err := http.Get(*website)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	fmt.Printf("Site :  %s\n", *website)

	client := &http.Client{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		payloads := scanner.Text()
		result := *website + payloads

		req, err := http.NewRequest(http.MethodGet, result, nil)
		if err != nil {
			fmt.Println("Error : ", err)
			os.Exit(1)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error : ", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}

		setCookie := response.Header.Get("Set-Cookie")

		if *onlyPositive {
			if setCookie != "" {
				green.Printf("[", payloads, "]     HIT")
			}
		} else {
			if setCookie != "" {
				green.Printf("[", payloads, "]     HIT")
			} else {
				red.Println("[", payloads, "]    :(")
			}
		}
	}
}
