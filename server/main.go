package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mattermost/mattermost/server/public/plugin"
)

func main() {
	if len(os.Args) > 1 {
		if err := demo(os.Args[1:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	plugin.ClientMain(&Plugin{})
}

// demo renders a meme to a file without a running server, which is how you
// preview a template, font or metadata change:
//
//	go run ./server -out demo.jpg 'memes. memes everywhere'
func demo(argv []string) error {
	flags := flag.NewFlagSet("demo", flag.ContinueOnError)
	out := flags.String("out", "", "path to write the rendered meme to")
	if err := flags.Parse(argv); err != nil {
		return err
	}
	if *out == "" {
		return fmt.Errorf("usage: -out <file.jpg> '<meme input>'")
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("expected exactly one input string, got %d", flags.NArg())
	}

	_, data, err := renderMemeJPEG(flags.Arg(0))
	if err != nil {
		return err
	}
	return os.WriteFile(*out, data, 0o644)
}
