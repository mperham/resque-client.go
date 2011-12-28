package main

import (
  "os"
  "fmt"
  "redis"
  "json"
  "strings"
  "strconv"
)

type ClientSpec struct {
  RedisLocation string
  Queue string
  WorkerCount int
}

type Client struct {
	spec *ClientSpec
	conn redis.Client
}

const (
	DefaultLocation    = "localhost:6379"
	DefaultQueue       = "test_queue"
  DefaultWorkerCount = 4
)

func DefaultSpec() *ClientSpec {
  return &ClientSpec{
    DefaultLocation,
    DefaultQueue,
    DefaultWorkerCount,
  }
}

type Job struct {
	Payload map[string] interface{}
}

func (c Client) Next() *Job {
	fmt.Println("Fetching job from", c.spec.Queue)
	job_data, e := c.conn.Lpop(c.spec.Queue)
  if e != nil {
    fmt.Println("Unable to fetch job", e)
    os.Exit(1)
  }
	fmt.Println(job_data)
	job := new(Job)
  json.Unmarshal(job_data, &job.Payload)
	return job
}

func NewClient(cspec *ClientSpec) Client{
  if cspec == nil {
    cspec = DefaultSpec()
  }

  cspec.Queue = fmt.Sprintf("resque:queue:%s", cspec.Queue)

  spec := redis.DefaultSpec()
  spec.Host(strings.Split(cspec.RedisLocation, ":")[0])
  i, err := strconv.Atoi(strings.Split(cspec.RedisLocation, ":")[1])
  if err != nil {
    fmt.Println("Invalid port", cspec.RedisLocation, err)
    os.Exit(5)
  }
  spec.Port(i)
	client, e := redis.NewSynchClientWithSpec(spec)
	if e != nil {
    fmt.Println ("failed to create the client", e);
    os.Exit(1)
  }
	return Client{
		cspec,
		client,
	}
}

func main() {
	client := NewClient(nil)
  job := client.Next()
  fmt.Println(job.Payload)
}
