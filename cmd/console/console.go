package console

import (
	"net/rpc"
	"io"
	"os"
	"os/signal"
	"fmt"
	"strings"
	"github.com/peterh/liner"
	"regexp"
)

var (
	passwordRegexp = regexp.MustCompile(`personal.[nus]`)
	onlyWhitespace = regexp.MustCompile(`^\s*$`)
	exit           = regexp.MustCompile(`^\s*exit\s*;*\s*$`)

	cmdnames = []string{"account", "createspace", "recharge", "setattr", "ttx", "test", "setvrf"}
)
var exec map[string]cmdfunc = map[string]cmdfunc{
}
type Console struct {
	name     string
	client   *rpc.Client
	prompt   string
	prompter UserPrompter
	histPath string
	history  []string
	printer  io.Writer
	dataDir  string
}

func (c *Console) execCmd(input string) {
	list := strings.Split(input, " ")
	args := []string{}
	for _, s := range list {
		if s != "" {
			args = append(args, s)
		}
	}

	cmd := args[0]
	args = args[1:]

	if v, ok := exec[cmd]; ok {
		v(c, args)
	} else {
		fmt.Println("command not found")
	}
}


func (c *Console) Interactive() {
	var (
		prompt    = c.prompt
		indents   = 0
		input     = ""
		scheduler = make(chan string)
	)

	go func() {
		for {
			line, err := c.prompter.PromptInput(<-scheduler)
			if err != nil {
				if err == liner.ErrPromptAborted { // ctrl-C
					prompt, indents, input = c.prompt, 0, ""
					scheduler <- ""
					continue
				}
				close(scheduler)
				return
			}
			scheduler <- line
		}
	}()
	abort := make(chan os.Signal, 1)
	signal.Notify(abort, os.Interrupt)

	for {
		scheduler <- prompt
		select {
		case <-abort:
			fmt.Fprintln(c.printer, "caught interrupt, exiting")
			return

		case line, ok := <-scheduler:
			if !ok || (indents <= 0 && exit.MatchString(line)) {
				return
			}
			if onlyWhitespace.MatchString(line) {
				continue
			}
			input += line
			c.execCmd(input)

			indents = countIndents(input)
			if indents <= 0 {
				prompt = c.prompt
			} else {
				prompt = strings.Repeat(".", indents*3) + " "
			}
			input = ""
		}
	}
}