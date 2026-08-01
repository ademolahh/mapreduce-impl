package worker

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ademolahh/map-reduce-impl/internal/shared"
)

func (w *Worker) reduceFunc(response shared.GetTaskResponse, mapResultFolder string,
	dirs []os.DirEntry, rdc func(string, []string) string) error {

	outPath := "words/result/mr-" + response.Input
	mr, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %q: %w", outPath, err)
	}
	defer mr.Close()

	var k shared.KeyValue
	for _, dir := range dirs {
		name := dir.Name()
		name = strings.TrimSuffix(name, filepath.Ext(name))
		r := strings.Split(name, "-")
		rdcTask := r[len(r)-1]

		if response.Input != string(rdcTask) {
			continue
		}

		data, err := os.ReadFile(mapResultFolder + "/" + dir.Name())
		if err != nil {
			continue
		}

		var result shared.KeyValue
		if err := json.Unmarshal(data, &result); err != nil {
			continue
		}

		k = append(k, result...)
	}

	sort.Sort(k)

	var values []string

	for i := 0; i < len(k)-1; i++ {
		values = append(values, k[i].Value)

		if k[i].Word != k[i+1].Word {
			result := rdc(k[i].Word, values)
			fmt.Fprintf(mr, "%s %s\n", k[i].Word, result)
			values = []string{}
		}
	}

	if len(k) > 0 {
		ps := k.Len() - 1
		values = append(values, k[ps].Value)
		result := rdc(k[ps].Word, values)
		fmt.Fprintf(mr, "%s %s\n", k[ps].Word, result)
	}

	req := shared.CompleteTaskRequest{
		Task:   response.Task,
		Output: response.Id,
	}
	err = w.client.Call("Master.CompleteTask", req, &shared.CompleteTaskResponse{})

	if err != nil {
		return fmt.Errorf("Completion Call Failed: %w", err)

	}

	return nil
}
