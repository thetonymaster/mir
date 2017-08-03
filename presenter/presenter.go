package presenter

import "fmt"

type Result struct {
	Task  string
	Time  float64
	Error error
}

func PrintResult(results []Result) {
	flag := false
	message := ""
	for _, r := range results {
		if r.Error != nil {
			flag = true
		}
		message = fmt.Sprintf("%s\n%s took %f", message, r.Task, r.Time)
	}

	if flag {
		fmt.Printf("\n\n%s%s\n", "FAIL", message)
	} else {
		fmt.Printf("\n\n%s%s\n", "SUCCESS", message)
	}
}
