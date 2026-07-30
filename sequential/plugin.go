package sequential

import (
	"fmt"
	"os"

	"github.com/ademolahh/map-reduce-impl/shared"
)

func runMapPhase(mps func(string, string) []shared.KV, folder string) (KeyValue, error) {
	dirs, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("failed to read folder %q: %w", folder, err)
	}

	var k KeyValue
	for _, dir := range dirs {
		path := folder + "/" + dir.Name()
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %q: %w", path, err)
		}
		content := string(data)

		res := mps(dir.Name(), content)
		k = append(k, res...)
	}

	return k, nil
}

func runReducePhase(k KeyValue, rdc func(string, []string) string, outPath string) error {
	mr, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %q: %w", outPath, err)
	}
	defer mr.Close()

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
	return nil
}
