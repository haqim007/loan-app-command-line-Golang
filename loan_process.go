package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type requester struct {
	name, id_card, address string
	age                    int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Gunakan camelcase untuk setiap argument \n")
	for {
		fmt.Print("$ ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		err = runCommand(cmdString)
		if err != nil {
			// os.Stderr,
			fmt.Println(err)
		}
	}
}

func runCommand(inputStr string) error {
	inputStr = strings.TrimSuffix(inputStr, "\n")
	arrInputStr := strings.Fields(inputStr)

	if arrInputStr[0] == "exit" {
		os.Exit(0)
	} else if arrInputStr[0] == "create_day_max" {
		if len(arrInputStr) < 2 {
			return errors.New("Peringatan!, Membutuhkan argumen jumlah request")
		} else {

			_, msgDayMax := createDayMax(arrInputStr[1])
			return errors.New(msgDayMax)
		}
	} else {
		return errors.New("Ketik 'help' untuk bantuan")
	}
	return errors.New("")
}

func createDayMax(dayMax string) (int, string) {
	dayMaxAsINt, _ := strconv.Atoi(dayMax)
	return dayMaxAsINt, "Created max request with " + dayMax
}
