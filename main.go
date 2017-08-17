package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jeffail/tunny"
	"github.com/thetonymaster/framework/configuration"
	"github.com/thetonymaster/framework/presenter"
	"github.com/thetonymaster/framework/provider/container"
	"github.com/thetonymaster/framework/provider/test"
	"github.com/thetonymaster/framework/repository"
)

func main() {

	args := os.Args[1:]

	conf, err := configuration.Read(args[0])
	if err != nil {
		log.Fatal(err)
	}

	repo, err := repository.NewInfluxDBClient("http://localhost:8086")
	if err != nil {
		log.Fatal(err)
	}
	p := presenter.NewPresenter(repo)

	pool, _ := tunny.CreatePool(conf.Containers.Limit, func(f interface{}) interface{} {
		input, _ := f.(func())
		input()
		return nil
	}).Open()
	defer pool.Close()

	for framework, configuration := range conf.Tests {
		runTests(framework, &configuration, conf, pool, p)
	}

}

func runTests(framework string, cfb *configuration.TestConfiguration,
	conf *configuration.Configuration, pool *tunny.WorkPool, p *presenter.Presenter) {
	done := make(chan bool, 1)
	results := make(chan presenter.Result, 100)
	var realTime float64

	switch framework {
	case "junit":

		dir, _ := filepath.Abs(filepath.Dir(conf.Tests["junit"].Path))
		containerProvider := container.NewDockerComposeGenerator([]string{conf.Tests["junit"].Path})
		jUnitTestProvider := test.NewJUnit(containerProvider, conf.Tests["junit"].Target, pool)
		tasks := jUnitTestProvider.GetFiles(dir + "/src/test/")
		jUnitTestProvider.Done = done
		jUnitTestProvider.Results = results
		start := time.Now()
		jUnitTestProvider.RunTask(tasks)
		elapsed := time.Since(start)
		realTime = elapsed.Seconds()
	}

	<-done

	var res []presenter.Result
	for r := range results {
		res = append(res, r)
	}
	p.PrintResult(res, realTime)
}
