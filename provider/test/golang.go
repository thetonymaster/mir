package test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jeffail/tunny"
	"github.com/thetonymaster/framework/presenter"
	"github.com/thetonymaster/framework/provider/container"
)

type Golang struct {
	Generator    generator
	Target       string
	pool         *tunny.WorkPool
	Done         chan bool
	Results      chan presenter.Result
	dockerClient *client.Client
	Repository   Repository
}

func NewGolang(generator generator, target string, pool *tunny.WorkPool) *Golang {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	return &Golang{
		Generator:    generator,
		Target:       target,
		pool:         pool,
		dockerClient: cli,
	}
}

func (g *Golang) RunTask(tasks []string) error {
	containers := g.Generator.New(JUnitProject, g.Target)
	time.Sleep(1000 * time.Millisecond)
	containers.Run()

	ctrs, _ := g.dockerClient.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	for ctr := range ctrs {
		for name := range ctrs[ctr].Names {
			if strings.Contains(ctrs[ctr].Names[name], g.Target) {
				d := time.Duration(3 * time.Second)
				g.dockerClient.ContainerStop(context.Background(), ctrs[ctr].ID, &d)
				g.dockerClient.ContainerRemove(context.Background(), ctrs[ctr].ID, types.ContainerRemoveOptions{})
				break
			}
		}
	}
	for _, task := range tasks {
		payload := g.getPayload(containers, g.Target, task)
		g.pool.SendWorkAsync(payload, nil)
		time.Sleep(1000 * time.Millisecond)

	}
	for g.pool.NumPendingAsyncJobs() > 0 {
		time.Sleep(200 * time.Millisecond)
	}
	containers.Kill()
	g.Done <- true
	close(g.Results)
	return nil
}

func (g *Golang) getPayload(containers *container.Container, target, task string) func() {
	return func() {
		start := time.Now()
		err := containers.Execute(target, "go", "test", "./cloudscale", fmt.Sprintf("-run=%s", task), "-v")
		elapsed := time.Since(start)
		result := presenter.Result{
			Task:  task,
			Time:  elapsed.Seconds(),
			Error: err,
		}
		ctrs, _ := g.dockerClient.ContainerList(context.Background(),
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

		logs, _ := g.dockerClient.ContainerLogs(context.Background(), id,
			types.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
				Tail:       "50",
			})
		defer logs.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(logs)
		result.Output = buf.String()

		g.dockerClient.ContainerRemove(context.Background(), id,
			types.ContainerRemoveOptions{})

		r := map[string]interface{}{
			"run_time": result.Time,
		}

		if result.Error != nil {
			r["error"] = result.Error.Error()
			r["output"] = result.Output
		}
		tag := map[string]string{
			"test":      result.Task,
			"framework": "golang",
		}
		err = g.Repository.Save("test_data", tag, r)

		log.Printf("%s took %s\n", task, elapsed)

		g.Results <- result
		if err != nil {
			fmt.Println(err)
		}
	}
}
