package main

import "github.com/ViBiOh/k8s-cleaner/pkg/job"

type services struct {
	job job.Service
}

func newServices(config configuration, clients clients) services {
	var output services

	output.job = job.New(config.job, clients.k8s)

	return output
}
