package main

import (
  "flag"
  "fmt"
  "github.com/MarconiProtocol/cli/console"
  "github.com/MarconiProtocol/cli/console/context"
  "github.com/MarconiProtocol/cli/console/execution"
  "github.com/MarconiProtocol/cli/console/modes/credentials"
  "github.com/MarconiProtocol/cli/console/modes/marconi_net"
  "github.com/MarconiProtocol/cli/console/modes/process/commands"
  "github.com/MarconiProtocol/cli/console/util"
  "github.com/MarconiProtocol/cli/core"
  "github.com/MarconiProtocol/cli/core/configs"
  "github.com/MarconiProtocol/cli/core/packages"
  mlog "github.com/MarconiProtocol/log"
  "os"
  "os/signal"
  "strings"
  "syscall"
)

const (
  MODE_EXEC      = "exec"
  MODE_NODE      = "node"
  MODE_UPGRADE   = "upgrade"
  MODE_DAEMON    = "daemon"
  LOG_CHILD_PATH = "/var/log/marconi"
)

func main() {
  mode := flag.String("mode", "node", "The mode Marconi Client will launch in")
  baseDir := flag.String("basedir", "/opt/marconi", "Base of directory tree that will "+
    "be used to store all configs and data files related to Marconi components, including "+
    "middleware, marconid, and meth")
  execCommand := flag.String("command", "", "The command to be executed")
  processStart := flag.String("start", "", "Start a process in background mode")
  processStop := flag.String("stop", "", "Stop a running process")
  processRestart := flag.String("restart", "", "Restart a running process")
  processList := flag.Bool("list", false, "List running processes")
  readCommandsFromStdin := flag.Bool("read-commands-from-stdin", false, "Whether to read commands from stdin")

  flag.Parse()
  configs.SetBaseDir(*baseDir)
  console.Init()
  mlog.Init(configs.GetFullPath(LOG_CHILD_PATH), "info")
  util.Logger, _ = mlog.GetLogInstance("mcli")
  defer core.Cleanup()

  // Catch sigterm and make sure we do core.Cleanup
  osIntChan := make(chan os.Signal, 1)
  signal.Notify(osIntChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
  go func() {
    <-osIntChan
    core.Cleanup()
    os.Exit(1)
  }()

  if !packages.CheckOrAskForMcliEulaAcknowledgement(*baseDir) {
    os.Exit(1)
  }

  signal := make(chan console.Exit)
  switch *mode {
  case MODE_NODE:
    // Bootstrap the mcli client
    core.Bootstrap(*baseDir)

    // Load the managed processes config
    core.StartProcessManager(*baseDir)

    go console.LaunchREPL(signal, true, *readCommandsFromStdin)

    // wait for signal
    <-signal

  case MODE_DAEMON:
    core.StartProcessManager(*baseDir)

    var processName string
    var mode string
    var argumentCount = 0

    if *processList {
      argumentCount++
    }
    if *processStart != "" {
      processName = *processStart
      mode = process_commands.START
      argumentCount++
    }
    if *processStop != "" {
      processName = *processStop
      mode = process_commands.STOP
      argumentCount++
    }
    if *processRestart != "" {
      processName = *processRestart
      mode = process_commands.RESTART
      argumentCount++
    }

    if argumentCount != 1 {
      fmt.Println("Only one of -start, -stop, -restart, -list can be specified at once")
      return
    }

    if *processList {
      console.ListProcesses()
    } else {
      console.LaunchProcess(mode, processName)
    }

  case MODE_UPGRADE:
    // NO-OP, nothing more to do than bootstrap
    core.Bootstrap(*baseDir)

  case MODE_EXEC:
    core.StartProcessManager(*baseDir)

    context := context.NewContext()

    // Register Modes onto the context
    context.RegisterMode(credentials.NewCredsMode(context), "Credential Mode")
    context.RegisterMode(marconi_net.NewMarconiNetMode(context), "Marconi Net Mode")

    execMode := execution.NewExecMode(context)

    if *execCommand == "" {
      fmt.Println("No commands specified Command Format: <mode> <command> with each command separated by ';'")
    } else {
      commands := strings.Split(*execCommand, ";")
      for _, command := range commands {
        args := strings.Split(strings.TrimSpace(command), " ")
        execMode.HandleCommand(args)
      }
    }

  default:
    fmt.Fprintln(os.Stderr, "Invalid mode.")
    flag.PrintDefaults()
    os.Exit(1)
  }
}
