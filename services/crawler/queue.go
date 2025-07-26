package main

import "fmt"

type Job struct {
	URL   string
	Depth int
}

type JobQueue struct {
	queue []Job
}

// for printing
func (j Job) String() string {
	return fmt.Sprintf("Job{URL: %s, Depth: %d}", j.URL, j.Depth)
}

// enqueue a job
func (q *JobQueue) Enqueue(job Job) {
	q.queue = append(q.queue, job)
}

// dequeue a job
func (q *JobQueue) Dequeue() Job {
	if len(q.queue) == 0 {
		return Job{}
	}
	job := q.queue[0]
	q.queue = q.queue[1:]
	return job
}

// simple is empty
func (q *JobQueue) IsEmpty() bool {
	return len(q.queue) == 0
}
