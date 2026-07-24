# MapReduce

## Sequential

A naive, single-process implementation of MapReduce (word count). Not the real, distributed thing.

```sh
go build -buildmode=plugin -o a.so ./apps/wordcount.go
go run ./sequential
```

Expects a `words/text/` folder with input files and a `words/result/` folder to write into (both must already exist). Reads files from `words/text/`, writes `word count` pairs to `words/result/result.txt`.

## Distributed

Coming soon.
