# MapReduce

## Sequential

A naive, single-process implementation of MapReduce (word count). Not the real, distributed thing.

```sh
go build -buildmode=plugin -o a.so ./apps/wordcount.go
go run ./sequential
```

Reads files from `words/text/`, writes `word count` pairs to `words/result/result.txt`.

## Distributed

Coming soon.
