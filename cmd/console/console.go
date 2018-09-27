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
	"path/filepath"
	"github.com/mattn/go-colorable"
)
type cmdfunc func(c *Console, args []string)
var exec map[string]cmdfunc = map[string]cmdfunc{
}

var (
	passwordRegexp = regexp.MustCompile(`personal.[nus]`)
	onlyWhitespace = regexp.MustCompile(`^\s*$`)
	exit           = regexp.MustCompile(`^\s*exit\s*;*\s*$`)

	cmdnames = []string{"account", "createspace", "recharge", "setattr", "ttx", "test", "setvrf"}
)

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
const HistoryFile = "history"

const DefaultPrompt = "> "

type Config struct {
	Name     string
	DataDir  string
	DocRoot  string
	Client   *rpc.Client
	Prompt   string
	Prompter UserPrompter
	Printer  io.Writer
	Preload  []string
}

func New(config Config) (*Console, error) {
	if config.Prompter == nil {
		config.Prompter = Stdin
	}
	if config.Prompt == "" {
		config.Prompt = DefaultPrompt
	}
	if config.Printer == nil {
		config.Printer = colorable.NewColorableStdout()
	}
	console := &Console{
		name:     config.Name,
		client:   config.Client,
		prompt:   config.Prompt,
		prompter: config.Prompter,
		printer:  config.Printer,
		histPath: filepath.Join(config.DataDir, HistoryFile),
		dataDir:  config.DataDir,
	}

	//console.prompter.SetWordCompleter(console.AutoCompleteInput)

	return console, nil
}
func (c *Console) Welcome() {
	fmt.Fprintf(c.printer, "Welcome to the unitcoin console!\n\n")
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

func countIndents(input string) int {
	var (
		indents     = 0
		inString    = false
		strOpenChar = ' '
		charEscaped = false
	)

	for _, c := range input {
		switch c {
		case '\\':
			if !charEscaped && inString {
				charEscaped = true
			}
		case '\'', '"':
			if inString && !charEscaped && strOpenChar == c {
				inString = false
			} else if !inString && !charEscaped {
				inString = true
				strOpenChar = c
			}
			charEscaped = false
		case '{', '(':
			if !inString {
				indents++
			}
			charEscaped = false
		case '}', ')':
			if !inString {
				indents--
			}
			charEscaped = false
		default:
			charEscaped = false
		}
	}

	return indents
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
