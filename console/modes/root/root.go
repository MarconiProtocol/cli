package root

import (
  "gitlab.neji.vm.tc/marconi/go-prompt"
  "gitlab.neji.vm.tc/marconi/cli/console/context"
  "gitlab.neji.vm.tc/marconi/cli/console/modes"
)

const (
  CREDS_CMD   = "credential"
  MNET_CMD    = "net"
  PROCESS_CMD = "process"
)

var ROOT_SUGGESTIONS = []prompt.Suggest{
  {Text: modes.JUMP_CMD, Description: "Mode jumping menu"},
  {Text: modes.EXIT_CMD, Description: "Exit mcli"},
}

/*
  Root cmd is meant to be the first menu/cmd in the marconi cli
*/
type RootMode struct {
  modes.BaseMode
}

/*
  Create a new root cmd struct,
  use this to create a RootMode otherwise suggestion and handlers wont be properly initialized
*/
func NewRootMode(c *context.Context) *RootMode {
  rootMode := RootMode{}
  rootMode.Init(c)
  rootMode.SetBaseSuggestions(ROOT_SUGGESTIONS)

  // handler and suggestions registration
  rootMode.RegisterCommand(modes.EXIT_CMD, rootMode.GetEmptySuggestions, rootMode.HandleExitCommand)

  return &rootMode
}

func (rm *RootMode) Name() string {
  return "home"
}

/*
  Handle the create meth cmd
*/
func (rm *RootMode) handleCredsCommand(args []string) {
  rm.HandleJumpCommand([]string{CREDS_CMD})
}

/*
  Handle the create mnet cmd
*/
func (rm *RootMode) handleMNetCommand(args []string) {
  rm.HandleJumpCommand([]string{MNET_CMD})
}

/*
  Handle the process cmd
*/
func (rm *RootMode) handleProcessCommand(args []string) {
  rm.HandleJumpCommand([]string{PROCESS_CMD})
}
