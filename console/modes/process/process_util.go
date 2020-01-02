package process

import (
  "github.com/MarconiProtocol/cli/console/modes/process/commands"
  "github.com/MarconiProtocol/cli/console/util"
  "github.com/MarconiProtocol/go-prompt"
)

// Mode suggestions
var PROCESS_UTIL_SUGGESTIONS = []prompt.Suggest{
  {Text: process_commands.RESET, Description: "Back up user credentials and clean up all the files of gmeth, marconid and middleware."},
}

func (mnm *ProcessMode) getRestSuggestions(line []string) []prompt.Suggest {
  return []prompt.Suggest{}
}

/*
  Handle the reset command
*/
func (pm *ProcessMode) handleReset(args []string) {
  util.Logger.Info(process_commands.UTIL+" "+process_commands.RESET, util.ArgsToString(args))
  process_commands.Reset(args)
}
