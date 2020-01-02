package execution

import (
  "fmt"
  "github.com/MarconiProtocol/cli/console/context"
  "github.com/MarconiProtocol/cli/console/modes"
  "github.com/MarconiProtocol/cli/console/util"
)

type ExecMode struct {
  context *context.Context
}

func NewExecMode(c *context.Context) *ExecMode {
  execMode := ExecMode{}
  execMode.context = c
  return &execMode
}

/*
  Handle commands when run in exec mode.
  Expected format is <mode> <command>
  ie. credential account list
*/
func (em *ExecMode) HandleCommand(args []string) {
  util.Logger.Info("EXEC MODE:", util.ArgsToString(args))

  if !modes.ArgsMinLenCheck(args, 2) {
    fmt.Println("USAGE: <mode> <command>")
    return
  }

  modeType := args[0]
  command := args[1:]

  mode, _ := em.context.SelectMode(modeType)

  if mode != nil {
    mode.HandleCommand(command)
  } else {
    fmt.Println("Invalid Mode Specified")
  }
}
