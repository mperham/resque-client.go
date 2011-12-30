package main

import (
  "fmt"
  "redis"
  "json"
  "log"
  "strings"
  "strconv"
  "time"
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
  log.Println("Fetching job from", c.spec.Queue)
  job_data, e := c.conn.Lpop(c.spec.Queue)
  if e != nil {
    log.Panicln("Unable to fetch job", e)
  }
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
    log.Panicln("Invalid port", cspec.RedisLocation, err)
  }
  spec.Port(i)
  client, e := redis.NewSynchClientWithSpec(spec)
  if e != nil {
    log.Panicln ("failed to create the client", e);
  }
  return Client{
    cspec,
    client,
  }
}

func (client *Client) Process() {
  job := client.Next()
  if len(job.Payload) == 0 {
    return
  }
  klass := job.Payload["class"]
  if klass == nil {
    log.Panicln("No class found in payload", job)
  }
  class, _ := klass.(string)
  function := workerFunctions[class]
  if function == nil {
    log.Panicln("No function registered for job", klass)
  }
  args := job.Payload["args"]
  argsArray, _ := args.([] interface {})
  argsMap := argsArray[0].(map[string] interface {})
  function(argsMap)
}

var workerFunctions = make(map[string] func(map[string] interface {}))

func Register(name string, processor func(map[string] interface {})) {
  workerFunctions[name] = processor
}

func myProcessor(args map[string] interface{}) {
  log.Println(args)
}

func main() {
  Register("Job", myProcessor)
  client := NewClient(nil)
  for {
    client.Process()
    time.Sleep(1000000000)
  }
}
