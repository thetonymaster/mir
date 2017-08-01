package test

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jeffail/tunny"
	"github.com/thetonymaster/test_framework/provider/container"
)

type provider interface {
	Run() error
	Execute(target string, task ...string) error
	Scale(containers map[string]int) error
	Kill() error
}

type generator interface {
	New(projectName string, args ...string) *container.Container
}

// JUnit runs the JUnit tests
type JUnit struct {
	Generator generator
	Target    string
	pool      *tunny.WorkPool
}

const JUnitProject = "junit"

// NewJUnit creates a new instance of a JUnit task manager
func NewJUnit(generator generator, target string, pool *tunny.WorkPool) *JUnit {
	return &JUnit{
		Generator: generator,
		Target:    target,
		pool:      pool,
	}
}

func (junit JUnit) GetFiles(searchDir string) []string {
	fileList := []string{}
	pattern := "(.+?)((Tests.java))"

	filepath.Walk(searchDir, func(filePath string, f os.FileInfo, err error) error {
		match, _ := regexp.MatchString(pattern, filePath)
		if match {
			name := strings.TrimSuffix(path.Base(filePath), filepath.Ext(filePath))
			fileList = append(fileList, name)
		}
		return nil
	})
	return fileList
}

func (junit *JUnit) RunTask(tasks []string) error {
	containers := junit.Generator.New(JUnitProject, junit.Target)
	containers.Run()

	for _, task := range tasks {
		payload := getPayload(containers, junit.Target, task)
		junit.pool.SendWorkAsync(payload, nil)
		time.Sleep(1 * time.Second)

	}
	for junit.pool.NumPendingAsyncJobs() > 0 {
		time.Sleep(1 * time.Second)
	}
	containers.Kill()
	return nil
}

func getPayload(containers *container.Container, target, task string) func() {
	return func() {
		time.Sleep(3 * time.Second)
		start := time.Now()
		taskLocal := task
		err := containers.Execute(target, "./mvnw", "test", fmt.Sprintf("-Dtest=%s", task))
		elapsed := time.Since(start)
		log.Printf("%s took %s\n", taskLocal, elapsed)
		if err != nil {
			fmt.Println(err)
		}
	}
}
