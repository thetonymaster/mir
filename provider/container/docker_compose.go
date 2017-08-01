package container

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

type Provider interface {
	Execute(target string, task ...string) error
	Kill() error
	Run() error
	Scale(containers map[string]int) error
}

type Container struct {
	Provider
}

// DockerComposeGenerator generates
type DockerComposeGenerator struct {
	Paths []string
}

// NewDockerComposeGenerator generates new DockerCompose instances
func NewDockerComposeGenerator(paths []string) *DockerComposeGenerator {
	return &DockerComposeGenerator{
		Paths: paths,
	}
}

// New returns a new container provider with docker compose
func (gen DockerComposeGenerator) New(projectName string, args ...string) *Container {
	compose := NewDockerCompose(projectName, args[0], gen.Paths)
	return &Container{
		compose,
	}
}

// DockerCompose runs docker compose tasks
type DockerCompose struct {
	project project.APIProject
	target  string
}

// NewDockerCompose creates a
func NewDockerCompose(projectName string, target string, files []string) *DockerCompose {
	context := &ctx.Context{
		Context: project.Context{
			ComposeFiles: files,
			ProjectName:  projectName,
		},
	}
	project, err := docker.NewProject(context, nil)

	if err != nil {
		log.Fatal(err)
	}

	return &DockerCompose{
		target:  target,
		project: project,
	}
}

func (dc DockerCompose) Scale(containers map[string]int) error {
	return dc.project.Scale(context.TODO(), 300, containers)
}

// Run runs the docker compose files
func (dc DockerCompose) Run() error {
	err := dc.project.Up(context.Background(), options.Up{})
	if err != nil {
		return err
	}
	err = dc.project.Kill(context.Background(), "SIGKILL", dc.target)
	time.Sleep(5 * time.Second)

	return err
}

// Kill kills all the containers related to the project
func (dc DockerCompose) Kill() error {
	return dc.project.Down(context.TODO(), options.Down{
		RemoveImages:  "local",
		RemoveOrphans: true,
		RemoveVolume:  true,
	})
}

// Execute Executes a task in a container
func (dc DockerCompose) Execute(target string, task ...string) error {
	out, errs := dc.runTask(target, task)
	select {
	case err := <-errs:
		return err
	case status := <-out:
		if status != 0 {
			return fmt.Errorf("Container exited with status %d", status)
		}
	}
	return nil
}

func (dc DockerCompose) runTask(target string, task []string) (<-chan int, <-chan error) {
	out := make(chan int, 1)
	errs := make(chan error, 1)

	go func() {
		status, err := dc.project.Run(context.TODO(),
			target,
			task,
			options.Run{
				Detached: false,
			})

		if err != nil {
			errs <- err
		} else {
			out <- status
		}
		close(out)
		close(errs)
	}()

	return out, errs
}
