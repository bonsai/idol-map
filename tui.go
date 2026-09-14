package idolmap

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// RunTUI starts a dependency-free terminal interface for the local RAG dataset.
func RunTUI(d Dataset) {
	in := bufio.NewScanner(os.Stdin)
	fmt.Println("idol-map TUI")
	fmt.Println("commands: search <query>, show <subject>, axes, quit")

	for {
		fmt.Print("> ")
		if !in.Scan() {
			return
		}
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		cmd := strings.ToLower(parts[0])
		arg := ""
		if len(parts) == 2 {
			arg = strings.TrimSpace(parts[1])
		}

		switch cmd {
		case "search", "s":
			for i, name := range d.Search(arg, 10) {
				fmt.Printf("%d. %s\n", i+1, name)
			}
		case "show":
			if arg == "" {
				fmt.Println("usage: show <subject>")
				continue
			}
			fmt.Print(d.Explain(arg))
		case "axes":
			for key, label := range d.Axes {
				fmt.Printf("%-12s %s\n", key, label)
			}
		case "quit", "q", "exit":
			return
		default:
			if n, err := strconv.Atoi(cmd); err == nil && n > 0 {
			fmt.Printf("selection %d\n", n)
			} else {
				fmt.Println("unknown command")
			}
		}
	}
}
