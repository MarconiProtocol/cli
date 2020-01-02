package process_commands

import (
  "fmt"
  "github.com/MarconiProtocol/cli/console/modes"
  "github.com/MarconiProtocol/cli/core/configs"
  "io/ioutil"
  "os"
  "path"
  "time"
)

// Commands
const (
  RESET = "reset"
)

var UTIL_COMMAND_MAP = map[string]func([]string){
  RESET: Reset,
}

func handleUtilCommand(args []string) {
  if !modes.ArgsMinLenCheck(args, 1) {
    fmt.Println("USAGE: <command>")
    return
  }
  commandType := args[0]
  commandArgs := args[1:]
  if commandHandlerFunction, present := UTIL_COMMAND_MAP[commandType]; present {
    commandHandlerFunction(commandArgs)
  } else {
    fmt.Println("Invalid command " + commandType)
  }
}

func Reset(args []string) {
  // stop all marconi processes
  StopProcess([]string{"gmeth"})
  StopProcess([]string{"marconid"})
  StopProcess([]string{"middleware"})

  // create the backup folder in baseDir if it does not exist
  baseDir := configs.GetBaseDir()
  backUpDir := path.Join(baseDir, "backup")
  if _, err := os.Stat(backUpDir); os.IsNotExist(err) {
    if err1 := os.Mkdir(backUpDir, os.ModePerm); err1 != nil {
      fmt.Printf("Reset failed: could not create folder %s, err: %s", backUpDir, err1.Error())
      return
    }
  }

  // back up the user's credentials
  accountsDir := path.Join(baseDir, "accounts")
  if _, err := os.Stat(accountsDir); !os.IsNotExist(err) {
    // create a folder in the backup folder using the current timestamp
    t := time.Now()
    formattedTimestamp := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    backUpTSDir := path.Join(backUpDir, formattedTimestamp)

    if err1 := os.Mkdir(backUpTSDir, os.ModePerm); err1 != nil {
      fmt.Printf("Reset failed: could not create folder %s, err: %s", backUpTSDir, err1.Error())
      return
    }

    // back up the credentials
    if err1 := os.Rename(accountsDir, path.Join(backUpTSDir, "accounts")); err1 != nil {
      fmt.Println("Reset failed: err:", err1)
      return
    } else {
      fmt.Println("User's credentials are backed up to", backUpTSDir)
    }
  }

  // remove all the downloaded binaries and other data, get everything back to when MCLI was just extracted from the tarball
  files, err := ioutil.ReadDir(baseDir)
  if err != nil {
    fmt.Printf("Reset failed: could not read from path %s, err: %s", baseDir, err.Error())
    return
  }

  for _, file := range files {
    switch file.Name() {
    case "bin":
      binDir := path.Join(baseDir, "bin")
      bins, err1 := ioutil.ReadDir(binDir)
      if err1 != nil {
        fmt.Printf("Reset failed: could not read from path %s, err: %s", binDir, err1.Error())
        return
      }

      for _, bin := range bins {
        if bin.Name() != "mcli" {
          if err1 := os.RemoveAll(path.Join(baseDir, "bin", bin.Name())); err1 != nil {
            fmt.Println("Reset failed: err:", err1)
            return
          }
        }
      }
    case "configs",
      "EULA.txt",
      "backup":
    default:
      if err1 := os.RemoveAll(path.Join(baseDir, file.Name())); err1 != nil {
        fmt.Println("Reset failed: err:", err1)
        return
      }
    }
  }
  fmt.Println("Reset successfully. Please exit and restart MCLI.")
}
