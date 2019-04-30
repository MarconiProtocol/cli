package main

import (
  "flag"
  "fmt"
  "os"
  "os/signal"
  "syscall"

  "gitlab.neji.vm.tc/marconi/cli/console"
  "gitlab.neji.vm.tc/marconi/cli/console/modes/process"
  "gitlab.neji.vm.tc/marconi/cli/core"
  "gitlab.neji.vm.tc/marconi/cli/core/configs"
  "gitlab.neji.vm.tc/marconi/cli/core/packages"
)

const (
  MODE_NODE     = "node"
  MODE_DOWNLOAD = "download"
  MODE_DAEMON   = "daemon"
)

func main() {
  mode := flag.String("mode", "node", "The mode Marconi Client will launch in")
  baseDir := flag.String("basedir", "/opt/marconi", "Base of directory tree that will "+
    "be used to store all configs and data files related to Marconi components, including "+
    "middleware, marconid, and meth")
  runConsole := flag.Bool("console", false, "Whether the console will be started")
  processStart := flag.String("start", "", "Start a process in background mode")
  processStop := flag.String("stop", "", "Stop a running process")
  processRestart := flag.String("restart", "", "Restart a running process")
  processList := flag.Bool("list", false, "List running processes")
  readCommandsFromStdin := flag.Bool("read-commands-from-stdin", false, "Whether to read commands from stdin")

  flag.Parse()

  configs.SetBaseDir(*baseDir)

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
    // Load the managed processes config
    core.StartProcessManager(*baseDir)

    go console.LaunchREPL(&signal, *runConsole, *readCommandsFromStdin)

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
      mode = process.START
      argumentCount++
    }
    if *processStop != "" {
      processName = *processStop
      mode = process.STOP
      argumentCount++
    }
    if *processRestart != "" {
      processName = *processRestart
      mode = process.RESTART
      argumentCount++
    }

    if (argumentCount != 1) {
      fmt.Println("Only one of -start, -stop, -restart, -list can be specified at once")
      return
    }

    if *processList {
      console.ListProcesses()
    } else {
      console.LaunchProcess(mode, processName)
    }

  case MODE_DOWNLOAD:
    // NO-OP, nothing more to do than bootstrap
    core.Bootstrap(*baseDir)

  default:
    fmt.Fprintln(os.Stderr, "Invalid mode.")
    flag.PrintDefaults()
    os.Exit(1)
  }
}
