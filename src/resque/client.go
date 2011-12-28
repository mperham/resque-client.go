package main

import (
  "os"
  "fmt"
  "redis"
  "json"
  "flag"
)


func main() {

  var queue string

  flag.StringVar(&queue, "queue", "", "queue to process")
  flag.Parse()

  if queue == "" {
    flag.PrintDefaults()
    os.Exit(1)
  } else {
    queue = fmt.Sprintf("resque:queue:%s", queue)
  }

  spec := redis.DefaultSpec()
	client, e := redis.NewSynchClientWithSpec(spec)
	if e != nil {
    fmt.Println ("failed to create the client", e);
    os.Exit(1)
  }

  job_data, e := client.Lpop(queue)
  if e != nil {
    fmt.Println("Unable to fetch job", e)
    os.Exit(1)
  }

  fmt.Println(job_data)

  var job map[string] interface{}
  json.Unmarshal(job_data, &job)
  fmt.Println(job)
}
