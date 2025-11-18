package utils

import (
	"bufio"
	"strings"
)

func SplitByLine(s string) <-chan struct {
	Index   int
	Content string
} {
	ch := make(chan struct {
		Index   int
		Content string
	})

	go func() {
		defer close(ch)

		scanner := bufio.NewScanner(strings.NewReader(s))
		index := 0

		for scanner.Scan() {
			ch <- struct {
				Index   int
				Content string
			}{
				Index:   index,
				Content: scanner.Text(),
			}
			index++
		}
	}()

	return ch
}
