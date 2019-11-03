package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Requester struct for loan requester.
type Requester struct {
	idCard, name, province string
	age                    int
	amount                 float64
}

//Loans data accepted or rejected
type Loans struct {
	loanID, idCard, status string
	amount                 float64
}

var maxRequest int = 0
var dataLoans = make(map[string]map[string]string)
var processedDataLoans = make(map[string]map[string]string)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("use snake case (ex : my_name) for each string argument \n")

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

	errorMsg := errors.New("")
	_ = errorMsg

	if arrInputStr[0] == "exit" {
		os.Exit(0)
	} else if arrInputStr[0] == "create_day_max" {
		if len(arrInputStr) < 2 {
			errorMsg = errors.New("Error!, Need max request argument \n 'create_day_max maxRequest' ")
		} else {
			intDayMax, msgDayMax := createDayMax(arrInputStr[1])
			maxRequest = intDayMax
			errorMsg = errors.New(msgDayMax)
		}
	} else if arrInputStr[0] == "add" {
		if len(arrInputStr) < 6 {
			errorMsg = errors.New("Error!, Need argument idcard, name, province, age and amount of loan  \n 'add idCard name province age amount_of_loan' ")
		} else {
			age, _ := strconv.Atoi(arrInputStr[4])
			amount, _ := strconv.ParseFloat(arrInputStr[5], 64)
			data := Requester{
				idCard:   arrInputStr[1],
				name:     arrInputStr[2],
				province: arrInputStr[3],
				age:      age,
				amount:   amount,
			}
			data.addLoans(&dataLoans)

			// myJSON, _ := json.MarshalIndent(dataLoans, "", "    ")
			// fmt.Println(string(myJSON))
		}
	} else if arrInputStr[0] == "status" {
		if len(arrInputStr) < 2 {
			errorMsg = errors.New("need Loan ID as argument")
		} else {
			errorMsg = errors.New(getStatusByLoanID(arrInputStr[1]))
		}

	} else {
		errorMsg = errors.New("Ketik 'help' untuk bantuan")
	}
	return errorMsg
}

func createDayMax(dayMax string) (int, string) {
	dayMaxAsInt, _ := strconv.Atoi(dayMax)
	return dayMaxAsInt, "Created max request with " + dayMax
}

func (data Requester) addLoans(dataLoans *map[string]map[string]string) {
	dt := time.Now()
	newData := map[string]string{
		"idCard":   data.idCard,
		"name":     data.name,
		"province": data.province,
		"age":      strconv.Itoa(data.age),
		"amount":   fmt.Sprintf("%.2f", data.amount),
	}
	status := "Rejected"
	msg := ""
	_ = msg
	_ = status
	lenDataLoans := len(*dataLoans)
	key := dt.Format("020106") + "0" + strconv.Itoa(lenDataLoans+1)
	arrProvince := []string{
		"DKI_JAKARTA",
		"JAWA_BARAT",
		"JAWA_TIMUR",
		"SUMATRA_UTARA",
	}

	if data.amount < 1000000 {
		msg = "The minimum amount is 1 million"
	} else if getTotalAmountByID(key) >= 10000000 {
		msg = "You have reach maximum loan amount"
	} else if data.age < 17 || data.age > 80 {
		msg = "The minimum and maximum age to request loan is 17 and 80"
	} else if inArray(data.province, arrProvince) == false {
		msg = "We only serve DKI JAKARTA, JAWA BARAT, JAWA TIMUR and SUMATRA UTARA area"
	} else {
		status = "accepted"
		msg = "sucess : " + key
	}

	(*dataLoans)[key] = newData

	loanData := Loans{
		loanID: key,
		idCard: data.idCard,
		status: status,
		amount: data.amount,
	}
	loanData.SaveLoanData()

}

func getTotalAmountByID(ID string) float64 {
	totalAmount := 0.00
	amount := 0.00
	for k, v := range dataLoans {
		if k == ID && len(v) > 0 {
			amount, _ = strconv.ParseFloat(v["amount"], 64)
			totalAmount = totalAmount + amount
		}
	}

	return totalAmount
}

func inArray(val string, array []string) bool {
	result := false
	for _, v := range array {
		if v == val {
			result = true
		}
	}
	return result
}

//SaveLoanData used for for saving loan data into map
func (data Loans) SaveLoanData() {
	processedDataLoans[data.loanID] = map[string]string{
		"idCard": data.idCard,
		"amount": fmt.Sprintf("%.2f", data.amount),
	}
}

func getStatusByLoanID(loanID string) string {
	result := ""
	_ = result
	for key, val := range processedDataLoans {
		if key == loanID {
			result = "Loan ID " + loanID + " is " + val["status"]
		}
	}

	return result
}
