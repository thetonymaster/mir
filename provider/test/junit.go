package test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jeffail/tunny"
	"github.com/thetonymaster/framework/presenter"
	"github.com/thetonymaster/framework/provider/container"
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
	Generator    generator
	Target       string
	pool         *tunny.WorkPool
	Done         chan bool
	Results      chan presenter.Result
	dockerClient *client.Client
}

const JUnitProject = "junit"

// NewJUnit creates a new instance of a JUnit task manager
func NewJUnit(generator generator, target string, pool *tunny.WorkPool) *JUnit {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	return &JUnit{
		Generator:    generator,
		Target:       target,
		pool:         pool,
		dockerClient: cli,
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
		payload := junit.getPayload(containers, junit.Target, task)
		junit.pool.SendWorkAsync(payload, nil)
		time.Sleep(200 * time.Millisecond)
	}
	for junit.pool.NumPendingAsyncJobs() > 0 {
		time.Sleep(1 * time.Second)
	}
	containers.Kill()
	junit.Done <- true
	close(junit.Results)
	return nil
}

func random(min, max int64) int64 {
	rand.Seed(time.Now().Unix())
	return rand.Int63n(max-min) + min
}

func (junit *JUnit) getPayload(containers *container.Container, target, task string) func() {
	return func() {
		start := time.Now()
		err := containers.Execute(target, "mvn", "surefire:test", fmt.Sprintf("-Dtest=%s", task))
		elapsed := time.Since(start)
		result := presenter.Result{
			Task:  task,
			Time:  elapsed.Seconds(),
			Error: err,
		}
		ctrs, _ := junit.dockerClient.ContainerList(context.Background(),
			types.ContainerListOptions{
				All: true,
			})

		id := ""
		for _, ctr := range ctrs {
			if strings.Contains(ctr.Command, task) {
				id = ctr.ID
				break
			}
		}

		logs, _ := junit.dockerClient.ContainerLogs(context.Background(), id,
			types.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
				Tail:       "50",
			})
		defer logs.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(logs)
		result.Output = buf.String()

		junit.dockerClient.ContainerRemove(context.Background(), id,
			types.ContainerRemoveOptions{})

		junit.Results <- result
		log.Printf("%s took %s\n", task, elapsed)
		if err != nil {
			fmt.Println(err)
		}
	}
}
