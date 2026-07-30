package shared

import (
	"fmt"
	"plugin"
)

func Fetch(program string) (func(string, string) []KV, func(string, []string) string, error) {
	pg, err := plugin.Open(program)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open program: %w", err)
	}

	mps, err := pg.Lookup("Map")
	if err != nil {
		return nil, nil, fmt.Errorf("plugin missing: %w", err)
	}

	mpFunc := mps.(func(string, string) []KV)

	rdc, err := pg.Lookup("Reduce")
	if err != nil {
		return nil, nil, fmt.Errorf("plugin missing: %w", err)
	}

	rdcFunc := rdc.(func(string, []string) string)

	return mpFunc, rdcFunc, nil
}
