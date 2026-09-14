package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	idolmap "github.com/bonsai/idol-map"
)

func usage() {
	fmt.Println(`idol-map — JSON/RAG/TUI CLI

Usage:
  idol-map search <query>
  idol-map show <subject>
  idol-map axes
  idol-map map --x <axis> --y <axis>
  idol-map tui

Options:
  --file PATH    dataset JSON (default: data/meta-6axis.json)`)
}

func main() {
	file := flag.String("file", "data/meta-6axis.json", "dataset JSON")
	flag.Parse()

	if flag.NArg() == 0 {
		usage()
		return
	}

	d, err := idolmap.LoadJSON(*file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "load:", err)
		os.Exit(1)
	}

	cmd := strings.ToLower(flag.Arg(0))

	switch cmd {
	case "search", "s":
		if flag.NArg() < 2 {
			fmt.Fprintln(os.Stderr, "usage: idol-map search <query>")
			os.Exit(2)
		}
		query := strings.Join(flag.Args()[1:], " ")
		for i, name := range d.Search(query, 10) {
			fmt.Printf("%d. %s\n", i+1, name)
		}

	case "show":
		if flag.NArg() < 2 {
			fmt.Fprintln(os.Stderr, "usage: idol-map show <subject>")
			os.Exit(2)
		}
		fmt.Print(d.Explain(strings.Join(flag.Args()[1:], " ")))

	case "axes":
		for key, label := range d.Axes {
			fmt.Printf("%-12s %s\n", key, label)
		}

	case "map":
		fs := flag.NewFlagSet("map", flag.ExitOnError)
		x := fs.String("x", "music", "X axis")
		y := fs.String("y", "expression", "Y axis")
		_ = fs.Parse(flag.Args()[1:])
		if _, ok := d.Axes[*x]; !ok {
			fmt.Fprintln(os.Stderr, "unknown X axis:", *x)
			os.Exit(2)
		}
		if _, ok := d.Axes[*y]; !ok {
			fmt.Fprintln(os.Stderr, "unknown Y axis:", *y)
			os.Exit(2)
		}
		fmt.Printf("%s × %s\n\n", d.Axes[*x], d.Axes[*y])
		for name, s := range d.Subjects {
			fmt.Printf("%-28s %.2f %.2f\n", name, s[*x], s[*y])
		}

	case "tui":
		idolmap.RunTUI(d)

	default:
		usage()
		os.Exit(2)
	}
}
