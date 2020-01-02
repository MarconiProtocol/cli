package process

import (
  "fmt"
  "github.com/MarconiProtocol/cli/console/context"
  "github.com/MarconiProtocol/cli/console/modes"
  "github.com/MarconiProtocol/cli/console/modes/process/commands"
  "github.com/MarconiProtocol/cli/console/util"
  "github.com/MarconiProtocol/cli/core/processes"
  "github.com/MarconiProtocol/go-prompt"
)

// Name of Mode
const NAME = "process"

// Mode suggestions
var PROCESS_SUGGESTIONS = []prompt.Suggest{
  {Text: process_commands.START, Description: "Start a process and manage its lifetime."},
  {Text: process_commands.STOP, Description: "Stop a managed process."},
  {Text: process_commands.RESTART, Description: "Restart a managed process."},
  {Text: process_commands.LIST, Description: "List running processes."},
  {Text: process_commands.VERSION, Description: "Show version of each component."},
  {Text: process_commands.UTIL, Description: "Utility commands"},
  {Text: modes.RETURN_TO_ROOT, Description: "Return to home menu"},
  {Text: modes.EXIT_CMD, Description: "Exit mcli"},
}

/*
  Menu/Mode for process operations
*/
type ProcessMode struct {
  modes.BaseMode
}

/*
  Create a new process cmd struct,
  use this to create a ProcessMode otherwise suggestion and handlers wont be properly initialized
*/
func NewProcessMode(c *context.Context) *ProcessMode {
  processMode := ProcessMode{}
  processMode.Init(c)
  processMode.SetBaseSuggestions(PROCESS_SUGGESTIONS)

  // suggestion and handler registrations
  processMode.RegisterCommand(process_commands.START, processMode.getSuggestions, processMode.handleStart)
  processMode.RegisterCommand(process_commands.STOP, processMode.getSuggestions, processMode.handleStop)
  processMode.RegisterCommand(process_commands.RESTART, processMode.getSuggestions, processMode.handleRestart)
  processMode.RegisterCommand(process_commands.LIST, processMode.GetEmptySuggestions, processMode.handleList)
  processMode.RegisterCommand(process_commands.VERSION, processMode.GetEmptySuggestions, processMode.handleVersion)
  processMode.RegisterCommand(process_commands.UTIL, processMode.getUtilSuggestions, processMode.handleUtil)
  processMode.RegisterSubCommand(process_commands.UTIL, process_commands.RESET, processMode.getRestSuggestions, processMode.handleReset)

  processMode.RegisterCommand(modes.RETURN_TO_ROOT, processMode.GetEmptySuggestions, processMode.HandleReturnToRoot)
  processMode.RegisterCommand(modes.EXIT_CMD, processMode.GetEmptySuggestions, processMode.HandleExitCommand)

  return &processMode
}

func (pm *ProcessMode) CliPrefix() (string, bool) {
  return pm.Name(), false
}

func (pm *ProcessMode) Name() string {
  return NAME
}

func (pm *ProcessMode) HandleCommand(args []string) {
  if !modes.ArgsMinLenCheck(args, 1) {
    fmt.Println("USAGE:  <command>")
    return
  }
  commandType := args[0]
  commandArgs := args[1:]
  if commandHandlerFunction, present := process_commands.COMMAND_MAP[commandType]; present {
    commandHandlerFunction(commandArgs)
  } else {
    fmt.Println("Invalid command " + commandType)
  }
}

func (pm *ProcessMode) getSuggestions(line []string) []prompt.Suggest {
  processConfigs := processes.Instance().GetSortedProcessConfigs()
  suggestions := make([]prompt.Suggest, len(processConfigs))

  for i, processConfig := range processConfigs {
    suggestions[i] = prompt.Suggest{Text: processConfig.Id, Description: ""}
  }

  return util.SimpleSubcommandCompleter(line, 1, suggestions)
}

/*
  Show prompt suggestions for entering the util sub menu
*/
func (pm *ProcessMode) getUtilSuggestions(line []string) []prompt.Suggest {
  return util.SimpleSubcommandCompleter(line, 1, PROCESS_UTIL_SUGGESTIONS)
}

func (pm *ProcessMode) handleStart(args []string) {
  util.Logger.Info(process_commands.START, util.ArgsToString(args))
  process_commands.StartProcesses(args)
}

func (pm *ProcessMode) handleStop(args []string) {
  util.Logger.Info(process_commands.STOP, util.ArgsToString(args))
  process_commands.StopProcess(args)
}

func (pm *ProcessMode) handleRestart(args []string) {
  util.Logger.Info(process_commands.RESTART, util.ArgsToString(args))
  process_commands.RestartProcess(args)
}

func (pm *ProcessMode) handleList(args []string) {
  util.Logger.Info(process_commands.LIST, util.ArgsToString(args))
  process_commands.ListProcesses(args)
}

func (pm *ProcessMode) handleVersion(args []string) {
  util.Logger.Info(process_commands.VERSION, util.ArgsToString(args))
  process_commands.ListProcessVersions(args)
}

func (pm *ProcessMode) handleUtil(args []string) {
  util.HandleFurtherCommands(process_commands.UTIL, PROCESS_UTIL_SUGGESTIONS)
}
