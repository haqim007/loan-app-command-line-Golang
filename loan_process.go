package main

import (
	"bufio"
	"encoding/json"
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
	loanID, idCard, status, date string
	amount                       float64
}

var maxRequest int = 0
var acceptedFrequentRequest int = 0

//input data loan
var dataLoans = make(map[string]map[string]string)

// processed data loan (Accepted or Rejected)
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

	if len(arrInputStr) > 0 {
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
				errorMsg = errors.New("Error!, Need argument idcard, name, province, age and amount of loan  \n example format : add idcard name province age amount \n use underscore instead of whitespace")
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
				_, msg := data.addLoans(&dataLoans)
				errorMsg = errors.New(msg)
				// myJSON, _ := json.MarshalIndent(dataLoans, "", "    ")
				// fmt.Println(string(myJSON))
			}
		} else if arrInputStr[0] == "status" {
			if len(arrInputStr) < 2 {
				errorMsg = errors.New("need Loan ID as argument")
			} else {
				errorMsg = errors.New(getStatusByLoanID(arrInputStr[1]))
			}

		} else if arrInputStr[0] == "show_all" {
			printAllLoans()
		} else if arrInputStr[0] == "show_all_with_status" {
			printAllWithStatus()
		} else if arrInputStr[0] == "installment" {
			if len(arrInputStr) < 3 {
				errorMsg = errors.New("need Loan ID and Multiple as argument \n example format: installment 05111901 03")
			} else {
				multiple, _ := strconv.Atoi(arrInputStr[2])
				installment(arrInputStr[1], multiple)
			}
		} else if arrInputStr[0] == "find_by_rejected_amount" {
			if len(arrInputStr) < 2 {
				errorMsg = errors.New("need loan amount as second argument  \n example format : find_by_rejected_amount (amount)")
			} else {
				findAmountStatus(arrInputStr[1], "Rejected")
			}
		} else if arrInputStr[0] == "find_by_accepted_amount" {
			if len(arrInputStr) < 2 {
				errorMsg = errors.New("need loan amount as second argument \n example format : find_by_accepted_amount (amount)")
			} else {
				findAmountStatus(arrInputStr[1], "Accepted")
			}
		} else {
			errorMsg = errors.New("Command not found \n Commands : \n create_day_max : to set maximum time request per day \n add : to add new loan request \n show_all : to show all requested loan \n show_all_with_status : to show all requested loans with status (accepted/rejected) \n status : to check status of a loanID \n find_by_rejected_amount : to find rejected loanID amount \n find_by_accepted_amount : to find accepted loanID by amount")
		}
	} else {
		errorMsg = errors.New("Command not found \n Commands : \n create_day_max : to set maximum time request per day \n add : to add new loan request \n show_all : to show all requested loan \n show_all_with_status : to show all requested loans with status (accepted/rejected) \n status : to check status of a loanID \n find_by_rejected_amount : to find rejected loanID amount \n find_by_accepted_amount : to find accepted loanID by amount \n exit : to close the app")
	}
	return errorMsg
}

func createDayMax(dayMax string) (int, string) {
	dayMaxAsInt, _ := strconv.Atoi(dayMax)
	return dayMaxAsInt, "Created max request with " + dayMax
}

// return status and msg
func (data Requester) addLoans(dataLoans *map[string]map[string]string) (string, string) {
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
	getCurrentLoan := getTotalAmountByID(data.idCard)
	getCurrentLoanStr := fmt.Sprintf("%.2f", getTotalAmountByID(data.idCard))

	if inArray(data.province, arrProvince) == false {
		msg = "Failed : We only serve DKI JAKARTA, JAWA BARAT, JAWA TIMUR and SUMATRA UTARA area"
	} else if data.age < 17 || data.age > 80 {
		msg = "Failed : The minimum and maximum age to request loan is 17 and 80"
	} else if data.amount < 1000000 {
		msg = "Failed : The minimum amount is 1 million"
	} else if getCurrentLoan+data.amount >= 10000000 {
		msg = "Failed : Maximum loan amount is 10 million. \n Your current loan is " + getCurrentLoanStr
	} else if maxRequest > 0 && acceptedFrequentRequest == maxRequest {
		msg = "Failed : You have reach limit to request loan"
	} else {
		status = "Accepted"
		msg = "Sucess : " + key
		acceptedFrequentRequest = acceptedFrequentRequest + 1
	}

	(*dataLoans)[key] = newData

	loanData := Loans{
		loanID: key,
		idCard: data.idCard,
		status: status,
		date:   dt.Format("020106"),
		amount: data.amount,
	}
	loanData.SaveLoanData()

	return status, msg
}

func getTotalAmountByID(ID string) float64 {
	totalAmount := 0.00
	amount := 0.00
	for _, v := range dataLoans {
		if v["idCard"] == ID && len(v) > 0 {
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
		"status": data.status,
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

func printAllLoans() {
	myJSON, _ := json.MarshalIndent(dataLoans, "", "    ")
	fmt.Println(string(myJSON))
}

func printAllWithStatus() {
	myJSON, _ := json.MarshalIndent(processedDataLoans, "", "    ")
	fmt.Println(string(myJSON))
}

func installment(key string, multiple int) {
	var table strings.Builder
	table.WriteString("Month | DueDate | AdministrationFee | Capital | Total \n")
	dueDate := time.Now()
	_ = dueDate

	amount, _ := strconv.ParseFloat(processedDataLoans[key]["amount"], 64)
	multipleFloat := float64(multiple)
	administrationFee := 100000.00 / multipleFloat
	administrationFeeStr := fmt.Sprintf("%.2f", administrationFee)
	capital := amount / multipleFloat
	capitalStr := fmt.Sprintf("%.2f", capital)
	total := fmt.Sprintf("%.2f", capital+administrationFee)
	for index := 1; index <= multiple; index++ {
		dueDate, _ = time.Parse("020106", processedDataLoans[key]["date"])
		dueDate = dueDate.AddDate(0, index, 0)
		table.WriteString("0" + strconv.Itoa(index) + " |" + dueDate.Format("020106") + " | " + administrationFeeStr + " | " + capitalStr + " | " + total)
		table.WriteString("\n")
	}

	fmt.Println(table.String())
}

func findAmountStatus(amount string, status string) {
	amountFLoat, _ := strconv.ParseFloat(amount, 64)
	existingAmountFloat := 0.00
	_, _ = amountFLoat, existingAmountFloat

	var result strings.Builder
	for key, val := range processedDataLoans {
		existingAmountFloat, _ = strconv.ParseFloat(val["amount"], 64)
		if amountFLoat == existingAmountFloat && val["status"] == status {
			result.WriteString(key + " ")
		}
	}
	if result.String() == "" {
		result.WriteString("Sorry, doesn't found it")
	}

	fmt.Println(result.String())
}
