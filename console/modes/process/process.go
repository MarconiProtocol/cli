package process

import (
  "gitlab.neji.vm.tc/marconi/go-prompt"
  "gitlab.neji.vm.tc/marconi/cli/console/context"
  "gitlab.neji.vm.tc/marconi/cli/console/modes"
  "gitlab.neji.vm.tc/marconi/cli/console/util"
  "gitlab.neji.vm.tc/marconi/cli/core/processes"
  "fmt"
  "strconv"
)

const (
  START = "start"
  STOP = "stop"
  RESTART = "restart"
  LIST = "list"
)

// Mode suggestions
var PROCESS_SUGGESTIONS = []prompt.Suggest{
  {Text: START, Description: "Start a process and manage its lifetime."},
  {Text: STOP, Description: "Stop a managed process."},
  {Text: RESTART, Description: "Restart a managed process."},
  {Text: LIST, Description: "List running processes."},
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
  processMode.RegisterCommand(START, processMode.getSuggestions, processMode.handleStart)
  processMode.RegisterCommand(STOP, processMode.getSuggestions, processMode.handleStop)
  processMode.RegisterCommand(RESTART, processMode.getSuggestions, processMode.handleRestart)
  processMode.RegisterCommand(LIST, processMode.GetEmptySuggestions, processMode.handleList)

  processMode.RegisterCommand(modes.RETURN_TO_ROOT, processMode.GetEmptySuggestions, processMode.HandleReturnToRoot)
  processMode.RegisterCommand(modes.EXIT_CMD, processMode.GetEmptySuggestions, processMode.HandleExitCommand)

  return &processMode
}

func (pm *ProcessMode) CliPrefix() (string, bool) {
  return pm.Name(), false
}

func (pm *ProcessMode) Name() string {
  return "process"
}

func (pm *ProcessMode) getSuggestions(line []string) []prompt.Suggest {
  processConfigs := processes.Instance().GetSortedProcessConfigs()
  suggestions := make([]prompt.Suggest, len(processConfigs))

  for i, processConfig := range processConfigs {
    suggestions[i] = prompt.Suggest{Text: processConfig.Id, Description: ""}
  }

  return util.SimpleSubcommandCompleter(line, 1, suggestions)
}

func (pm *ProcessMode) parseArgs(args []string) (string, []string, bool) {
  if len(args) <= 0 {
    fmt.Println("No process argument provided")
    return "", nil, false
  }

  program := args[0]
  parsedArgs := args[1:]

  if !processes.Instance().ContainsId(program) {
    fmt.Println("Unrecognized process argument:", program)
    return "", nil, false
  }

  return program, parsedArgs, true
}

func (pm *ProcessMode) handleStart(args []string) {
  program, parsedArgs, exists := pm.parseArgs(args)
  if !exists {
    return
  }

  background := false
  if len(parsedArgs) == 1 {
    background, _ = strconv.ParseBool(parsedArgs[0])
  }
  processes.Instance().StartProcesses([]string{program}, background)
}

func (pm *ProcessMode) handleStop(args []string) {
  program, _, exists := pm.parseArgs(args)
  if !exists {
    return
  }

  processes.Instance().KillProcess(program)
}

func (pm *ProcessMode) handleRestart(args []string) {
  program, _, exists := pm.parseArgs(args)
  if !exists {
    return
  }

  pm.handleStop([]string{program})
  pm.handleStart(args)
}

func (pm *ProcessMode) handleList(args []string) {
  statuses := processes.Instance().GetProcessRunningMap()
  for _, processConfig := range processes.Instance().GetSortedProcessConfigs() {
    var runString string
    if statuses[processConfig.Id] {
      runString = "RUNNING"
    } else {
      runString = "STOPPED"
    }

    fmt.Printf("%-15s %s\n", processConfig.Id, runString)
  }
}
