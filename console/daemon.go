package console

import (
  "github.com/MarconiProtocol/cli/console/mode_interface"
  "github.com/MarconiProtocol/cli/console/modes/process/commands"
)

/*
  Send a "list" command to the process mode
*/
func ListProcesses() {
  processMode := getProcessMode()
  processMode.HandleSelection(processMode, "list", []string{})
}

/*
  Send a command to the process mode
*/
func LaunchProcess(mode string, processName string) {
  processMode := getProcessMode()

  args := []string{processName}
  // When the mode is START or RESTART we want to run the process in the background, passing true as the 2nd arg will do that
  if mode == process_commands.START || mode == process_commands.RESTART {
    args = append(args, "true")
  }

  processMode.HandleSelection(processMode, mode, args)
}

/*
  Returns the registered processMode from the contxt
*/
func getProcessMode() mode_interface.Mode {
  processMode, _ := contxt.SelectMode("process")
  return processMode
}
