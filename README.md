# koron-go/jobrun

[![GoDoc](https://godoc.org/github.com/koron-go/jobrun?status.svg)](https://godoc.org/github.com/koron-go/jobrun)
[![Actions/Go](https://github.com/koron-go/jobrun/workflows/Go/badge.svg)](https://github.com/koron-go/jobrun/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/jobrun)](https://goreportcard.com/report/github.com/koron-go/jobrun)

Package `jobrun` provides utilities to execute multiple jobs in serial or parallel based on `context.Context`.

`jobrun` は複数のジョブを `context.Context` ベースで直列や並列に実行するのを簡単にするパッケージです。

## Getting started

Install or update:

```console
$ go get -u github.com/koron-go/jobrun
```

## Examples

Run jobs in sequentially (serial).  This runs `job1` then `job2`.

```go
job1 := jobrun.RunFunc(func(ctx context.Context) error {
    // TODO: do something
    return nil
})
job2 := jobrun.RunFunc(func(ctx context.Context) error {
    // TODO: do something
    return nil
})

err := jobrun.Serial{job1, job2}.Run(context.Background())

// ...(snip)...
```

If you want to run both jobs in simultaneously (parallel):

```go
err := jobrun.Parallel{job1, job2}.Run(context.Background())
```
