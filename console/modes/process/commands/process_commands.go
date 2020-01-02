package process_commands

import (
  "fmt"
  "github.com/MarconiProtocol/cli/console/util"
  "github.com/MarconiProtocol/cli/core/configs"
  "github.com/MarconiProtocol/cli/core/processes"
  "io/ioutil"
  "path/filepath"
  "strconv"
  "strings"
  "time"
)

// Commands
const (
  START   = "start"
  STOP    = "stop"
  RESTART = "restart"
  LIST    = "list"
  VERSION = "version"
  UTIL    = "util"
)

var COMMAND_MAP = map[string]func([]string){
  START:   StartProcesses,
  STOP:    StopProcess,
  RESTART: RestartProcess,
  LIST:    ListProcesses,
  VERSION: ListProcessVersions,
  UTIL:    handleUtilCommand,
}

func parseArgs(args []string) (string, []string, bool) {
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

func StartProcesses(args []string) {
  program, parsedArgs, exists := parseArgs(args)
  if !exists {
    return
  }

  background := false
  if len(parsedArgs) == 1 {
    background, _ = strconv.ParseBool(parsedArgs[0])
  }

  processes.Instance().StartProcesses([]string{program}, background)
}

func StopProcess(args []string) {
  program, _, exists := parseArgs(args)
  if !exists {
    return
  }

  processes.Instance().KillProcess(program)
}

func RestartProcess(args []string) {
  program, _, exists := parseArgs(args)
  if !exists {
    return
  }

  StopProcess(args)
  // use ps aux to check if this program is terminated
  for attempts := 0; attempts <= 3; attempts++ {
    if attempts == 3 {
      fmt.Println("restart " + program + " failed")
      util.Logger.Error("Error: restart " + program + " failed")
      return
    }
    time.Sleep(3 * time.Second)
    if processes.Instance().CheckProcessExistence(program) {
      fmt.Println(program + " is still running, retrying...")
    } else {
      break
    }
  }
  StartProcesses(args)
}

func ListProcesses(args []string) {
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

func ListProcessVersions(args []string) {
  packagesConf := configs.LoadPackagesConf()
  for _, config := range packagesConf.Packages {
    versionFilePath := ""
    if config.VersionFile != "" {
      versionFilePath = filepath.Join(configs.GetBaseDir(), config.VersionFile)
    } else {
      versionFilePath = filepath.Join(configs.GetBaseDir(), config.Dir, "version.txt")
    }
    data, err := ioutil.ReadFile(versionFilePath)
    if err != nil {
      panic(err)
    }
    version := strings.Trim(strings.Split(string(data), "=")[1], "\n")
    fmt.Println(config.Id + ": " + version)
  }
  fmt.Println("mcli:", packagesConf.Version)
}
