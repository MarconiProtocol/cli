package modes

import (
  "fmt"
  "gitlab.neji.vm.tc/marconi/go-prompt"
  "gitlab.neji.vm.tc/marconi/cli/console/context"
  "gitlab.neji.vm.tc/marconi/cli/console/util"
  "gitlab.neji.vm.tc/marconi/cli/console/mode_interface"
  "strings"
)

const (
  RETURN_TO_ROOT = "home"
  EXIT_CMD       = "exit"
  JUMP_CMD       = "j"
)

type Command struct {
  suggestion func([]string) []prompt.Suggest
  handler    func(args []string)
}

/*
  BaseMode handles breaking down the current input ( line ) into command and args
  Handles showing different suggestions depending on the command

  Used as a component in other structs to statisfy the cmd interface
*/
type BaseMode struct {
  baseSuggestions []prompt.Suggest
  commands        map[string]Command
  subcommands     map[string]map[string]Command
  contxt          *context.Context
}

/*
  Initializes the basemode
*/
func (bm *BaseMode) Init(contxt *context.Context) {
  bm.commands = make(map[string]Command)
  bm.subcommands = make(map[string]map[string]Command)
  bm.contxt = contxt

  bm.RegisterCommand(EXIT_CMD, bm.GetEmptySuggestions, bm.HandleExitCommand)
  bm.RegisterCommand(JUMP_CMD, bm.contxt.ModeSuggestionsFunc, bm.HandleJumpCommand)
}

func (bm *BaseMode) SetBaseSuggestions(suggestions []prompt.Suggest) {
  bm.baseSuggestions = suggestions
}

func (bm *BaseMode) RegisterCommand(command string, suggestionFunc func([]string) []prompt.Suggest, commandFunc func(args []string)) {
  c := Command{suggestionFunc, commandFunc}
  bm.commands[command] = c
}

func (bm *BaseMode) RegisterSubCommand(command string, subcommand string, suggestionFunc func([]string) []prompt.Suggest, commandFunc func(args []string)) {
  c := Command{suggestionFunc, commandFunc}
  if bm.subcommands[command] == nil {
    bm.subcommands[command] = make(map[string]Command)
  }
  bm.subcommands[command][subcommand] = c
}

/*
  Based on the current line, split into command and args,
  Pass this to the registered suggestion functions to get a specific suggestion for a command
*/
func (bm *BaseMode) Completer(d prompt.Document) []prompt.Suggest {
  // Try to split the current line by space to extract a command, if non exists, just use the base suggestions
  if line := util.SplitInput(d.CurrentLineBeforeCursor(), true); len(line) > 1 {
    // Return the correct command suggestion based on the command
    command := line[0]
    if c, exists := bm.commands[command]; exists {
      subcommand := line[1]
      if sc, exists := bm.subcommands[command][subcommand]; exists {
        return sc.suggestion(line[1:])
      }
      return c.suggestion(line)
    } else {
      return []prompt.Suggest{}
    }
  } else {
    return prompt.FilterHasPrefix(bm.baseSuggestions, d.GetWordBeforeCursor(), true)
  }
}

/*
  Once a selection has been made, pass args to command handler
*/
func (bm *BaseMode) HandleSelection(currCmd mode_interface.Mode, selection string, args []string) {
  if c, exists := bm.commands[selection]; exists {
    if len(args) > 0 {
      if sc, exists := bm.subcommands[selection][args[0]]; exists {
        sc.handler(args[1:])
        return
      }
    }
    c.handler(args)
  } else {
    fmt.Println("No command found named", selection)
  }
}

/*
  Just a simple prompt prefix
*/
func (bm *BaseMode) CliPrefix() (string, bool) {
  return "", false
}

/*
  For when there is no need for any suggestions
*/
func (bm *BaseMode) GetEmptySuggestions(line []string) []prompt.Suggest {
  return []prompt.Suggest{}
}

/*
  Triggers the mcli to exit by passing back a nil instead of a valid cmd
*/
func (bm *BaseMode) HandleExitCommand(args []string) {
  bm.contxt.Exit()
}

/*
  Updates the context's active mode based on the args provided
*/
func (bm *BaseMode) HandleJumpCommand(args []string) {
  if !ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", JUMP_CMD, "<mode_name>")
    return
  }

  modeStr := strings.TrimSpace(args[0])
  _, err := bm.contxt.SelectMode(modeStr)
  if err != nil {
    fmt.Println("Tried to jump to unregistered mode", modeStr)
  }
}

/*
  Jumps to the root mode
 */
func (bm *BaseMode) HandleReturnToRoot(args []string) {
  bm.HandleJumpCommand([]string{"home"})
}

/*
  Returns a suggestion to ask for the user's wallet password
*/
func PasswordCompleter(d prompt.Document) []prompt.Suggest {
  return []prompt.Suggest{{Text: "<PASSWORD>", Description: "Your wallet password"}}
}

/*
  A suggestion for the confirmation input
*/
func confirmationCompleter(d prompt.Document) []prompt.Suggest {
  return []prompt.Suggest{{Text: "", Description: "Please enter yes or no"}}
}

/*
  Loop to get user input for a confirmation
*/
func GetConfirmationInput() bool {
  for {
    fmt.Println("Please enter yes or no")
    input := prompt.Input("", confirmationCompleter)
    if input == "yes" || input == "y" || input == "Y" {
      return true
    } else if input == "no" || input == "n" || input == "N" {
      return false
    }
  }
}
