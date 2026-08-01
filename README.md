# MapReduce

MapReduce is a programming model for processing large data sets in parallel: a `Map` function turns input into key/value pairs, and a `Reduce` function combines all values for each key. This repo is a learning implementation, not a match for the [original paper](https://research.google/pubs/mapreduce-simplified-data-processing-on-large-clusters/) — no GFS-backed storage, no data-locality-aware scheduling, and fault tolerance here is just a fixed 10s task timeout instead of the paper's fuller design.

## Sequential

A naive, single-process implementation of MapReduce (word count).

```sh
go build -buildmode=plugin -o a.so ./apps
go run ./cmd/sequential -build=a.so
```

Expects a `words/text/` folder with input files and a `words/result/` folder to write into (both must already exist). Reads files from `words/text/`, writes `word count` pairs to `words/result/result.txt`.

## Distributed

A naive master/worker implementation over RPC. Not real multi-machine distribution — the master and all workers run as separate local processes talking over a Unix socket. One master process hands out map and reduce tasks; one or more worker processes do the work.

```sh
mkdir -p distributed/worker mapresult tmp words/result
go build -buildmode=plugin -o distributed/worker/wordcount.so ./apps

go run ./cmd/distributed/master   # one master
go run ./cmd/distributed/worker   # one or more workers, each in its own terminal
```

Expects `words/text/` (input), `mapresult/` and `tmp/` (intermediate map output), and `words/result/` (final output) to already exist. Reduce output is written per partition to `words/result/mr-<n>`, one file per reduce task. The reduce partition count is fixed at 5.
