package console

import (
  "bufio"
  "fmt"
  "gitlab.neji.vm.tc/marconi/go-prompt"
  "gitlab.neji.vm.tc/marconi/cli/console/context"
  "gitlab.neji.vm.tc/marconi/cli/console/util"
  "gitlab.neji.vm.tc/marconi/cli/console/modes/credentials"
  "gitlab.neji.vm.tc/marconi/cli/console/modes/marconi_net"
  "gitlab.neji.vm.tc/marconi/cli/console/modes/process"
  "gitlab.neji.vm.tc/marconi/cli/console/modes/root"
  "io"
  "os"
  "strings"
  "syscall"
)

const (
  PROMPT_SUFFIX = "> "
)

type Exit struct {}

var contxt *context.Context

func dummyExecutor(in string) { return }

// Initialization function for this package called by golang runtime
func init() {
  initializeContext()
}

/*
  Simple REPL that loops until the user exits
*/
func LaunchREPL(sig *chan Exit, runConsole bool, readCommandsFromStdin bool) {

  // Grab input commands from stdin if required
  var inputs_from_stdin []string
  if readCommandsFromStdin {
    inputs_from_stdin = readInputFromStdin()
  }

  // Set the current mode to home
  currentMode, _ := contxt.SelectMode("home")
  var prompt *prompt.Prompt = nil
  if runConsole {
    prompt = getPrompt(currentMode.Completer, currentMode.CliPrefix)
  }

  // Loop until console is exited
  for {
    currentMode, _ := contxt.GetCurrentMode()

    input := ""
    shouldExit := false

    // Get input for current mode based on conditions
    if len(inputs_from_stdin) > 0 {
      input, inputs_from_stdin = inputs_from_stdin[0], inputs_from_stdin[1:]
    } else if runConsole {
      input, shouldExit = prompt.InputReturnShouldExit()
    } else {
      shouldExit = true
    }

    // If the cli input should not cause an exit, continue to input parsing
    if !shouldExit {
      inputSlice := util.SplitInput(input, false)
      if len(inputSlice) > 0 {
        // command_selection and arguments
        currentMode.HandleSelection(currentMode, inputSlice[0], inputSlice[1:])
      }
    } else {
      // Otherwise trigger a context exit which will be caught below
      contxt.Exit()
    }

    // Check for the current command after handling the input
    newMode, _ := contxt.GetCurrentMode()
    if newMode == nil {
      *sig <- Exit{}
      break
    } else if newMode != currentMode {
      currentMode = newMode
      if runConsole {
        prompt = getPrompt(currentMode.Completer, currentMode.CliPrefix)
      }
    }
  }
}

/*
  Initialize contxt and register modes to it
*/
func initializeContext() {
  // TODO: make root part of mcli framework, support loading plugins
  contxt = context.NewContext()

  // future - register modes defined by loaded plugins
  contxt.RegisterMode(credentials.NewCredsMode(contxt), "Credential Mode")
  contxt.RegisterMode(marconi_net.NewMarconiNetMode(contxt), "Marconi Net Mode")
  contxt.RegisterMode(process.NewProcessMode(contxt), "Process Mode")

  contxt.RegisterMode(root.NewRootMode(contxt), "Home")
}

/*
  Parse stdin as command inputs and return a slice of them
*/
func readInputFromStdin() []string {
  var inputs_from_stdin []string
  // Read everything, then later feed it to the CLI.
  reader := bufio.NewReader(os.Stdin)
  for {
    line, err := reader.ReadString('\n')
    if len(line) > 0 {
      line = strings.TrimSuffix(line, "\n")
      inputs_from_stdin = append(inputs_from_stdin, line)
    }
    if err == io.EOF {
      break
    } else if err != nil {
      fmt.Println("Error reading stdin:", err)
      break
    } else if len(inputs_from_stdin) > 10000 {
      fmt.Println("Too many inputs on stdin, ignoring the rest.")
      break
    }
  }
  return inputs_from_stdin
}

/*
  Get the prompt to use with the console
*/
func getPrompt(completer func(d prompt.Document) []prompt.Suggest, cliPrefixFunc func() (prefix string, useLivePrefix bool)) *prompt.Prompt {

  cliPrefixFuncAppended := func() (prefix string, useLivePrefix bool) {
    prefix, useLivePrefix = cliPrefixFunc()
    prefix += PROMPT_SUFFIX
    return
  }

  prefix, _ := cliPrefixFuncAppended()

  return prompt.New(dummyExecutor, completer,
    prompt.OptionTitle("Marconi Client"),
    prompt.OptionAddKeyBind(
      prompt.KeyBind{Key: prompt.ControlZ, Fn: sigstopSelf}),
    prompt.OptionLivePrefix(cliPrefixFuncAppended),
    prompt.OptionPrefix(prefix),
    prompt.OptionPrefixTextColor(prompt.Blue),
    prompt.OptionInputTextColor(prompt.White),
    prompt.OptionPreviewSuggestionTextColor(prompt.Black),
    prompt.OptionPreviewSuggestionBGColor(prompt.White),
    prompt.OptionSuggestionTextColor(prompt.White),
    prompt.OptionSuggestionBGColor(prompt.Blue),
    prompt.OptionDescriptionTextColor(prompt.Black),
    prompt.OptionDescriptionBGColor(prompt.White),
    prompt.OptionSelectedSuggestionTextColor(prompt.White),
    prompt.OptionSelectedSuggestionBGColor(prompt.DarkBlue),
    prompt.OptionSelectedDescriptionTextColor(prompt.White),
    prompt.OptionSelectedDescriptionBGColor(prompt.Blue),
    prompt.OptionScrollbarThumbColor(prompt.DarkBlue),
    prompt.OptionScrollbarBGColor(prompt.Black),
  )
}

/*
  Kills this process
*/
func sigstopSelf(buffer *prompt.Buffer) {
  pid := os.Getpid()
  syscall.Kill(pid, syscall.SIGSTOP)
}
