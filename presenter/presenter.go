package presenter

import (
	"fmt"
	"strings"
)

type Result struct {
	Task   string
	Time   float64
	Error  error
	Output string
}

type Repository interface {
	Save(table string, data map[string]interface{}) error
	Close() error
}

type Presenter struct {
	Repository Repository
}

func NewPresenter(repository Repository) *Presenter {
	return &Presenter{
		Repository: repository,
	}
}

func (p *Presenter) PrintResult(results []Result, realTime float64) {
	flag := false
	message := ""
	var total float64
	for _, r := range results {
		f := ""

		if r.Error != nil {
			flag = true
			logs := strings.Split(r.Output, "\n")
			for i := range logs {
				if strings.Contains(logs[i], "Failed tests:") {
					for index := i + 1; index < len(logs); index++ {
						if strings.Contains(logs[index], "Tests run:") {
							break
						}
						f = fmt.Sprintf("%s\n%s", f, logs[index])
					}
					break
				}

			}
		}
		message = fmt.Sprintf("%s\n%s took %f%s", message, r.Task, r.Time, f)
		total += r.Time
	}
	avg := total / float64(len(results))
	if flag {
		fmt.Printf("\n\n%s%s\nTOTAL: %f\nAVERAGE: %f\nREAL TIME: %f\n", "FAIL", message, total, avg, realTime)
	} else {
		fmt.Printf("\n\n%s%s\nTOTAL: %f\nAVERAGE: %f\nREAL TIME: %f\n", "SUCCESS", message, total, avg, realTime)
	}

	p.Save(realTime, avg, total)
}

func (p *Presenter) Save(realTime, average, total float64) error {
	r := map[string]interface{}{
		"total_time": total,
		"average":    average,
		"real_time":  realTime,
	}

	err := p.Repository.Save("results", r)
	if err != nil {
		return err
	}

	return nil
}
