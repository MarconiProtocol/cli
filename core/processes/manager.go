package processes

import (
  "fmt"
  "gitlab.neji.vm.tc/marconi/cli/core/configs"
  "io/ioutil"
  "os"
  "os/exec"
  "path/filepath"
  "strconv"
  "sync"
  "syscall"
  "time"
)

const (
  LOG_DIR = "var/log/marconi"
  PID_DIR = "var/pid/marconi"

  MIDDLEWARE_ID = "middleware"
)

/*
  Simple configurable manager that starts and stops processes
*/
type ProcessManager struct {
  processMap      map[string]*os.Process
  processesConfig configs.ProcessesConfig
  baseDir         string
}

var instance *ProcessManager
var once sync.Once

func Instance() *ProcessManager {
  once.Do(func() {
    instance = &ProcessManager{}
    instance.processMap = make(map[string]*os.Process)
  })
  return instance
}

func (pm *ProcessManager) InitProcessManager(baseDir string, processesConfig configs.ProcessesConfig) {
  pm.processesConfig = processesConfig
  pm.baseDir = baseDir
  pm.monitorAllProcesses()
}

func (pm *ProcessManager) GetSortedProcessConfigs() []configs.ProcessConfig {
  return buildDependencyGraph(pm.processesConfig.Processes).getOrderedProcessConfigs()
}

func (pm *ProcessManager) ContainsId(id string) bool {
  for _, config := range pm.processesConfig.Processes {
    if config.Id == id {
      return true
    }
  }
  return false
}

/*
  Start processes as defined by procConfigs
  New goroutines are spawned to start the processes
  TODO: we need to have way to signal between process_manager and process to orchestrate next process, ghetto sleep for now
*/
func (pm *ProcessManager) StartProcesses(process_names []string, background bool) {
  procConfigs := []configs.ProcessConfig{}
  for _, config := range pm.processesConfig.Processes {
    for _, name := range process_names {
      if config.Id == name {
        procConfigs = append(procConfigs, config)
        break
      }
    }
  }

  // Build dependency graph and calculate ordered execution
  sortedProcConfigs := buildDependencyGraph(procConfigs).getOrderedProcessConfigs()

  for _, config := range sortedProcConfigs {
    // Either run process command in a coroutine or in the same thread
    if config.WaitForCompletion {
      pm.startProcess(config, background)
    } else {
      go pm.startProcess(config, background)
      // This is retarded, but temporary
      if config.WaitTime > 0 {
        time.Sleep(time.Duration(config.WaitTime) * time.Second)
      }
    }
  }
}

/*
  Start a single process, reference to os.Process object stored in process_manager's processMap
  Pipes process output to a log
*/
func (pm *ProcessManager) startProcess(cfg configs.ProcessConfig, background bool) {
  fmt.Println("STARTING PROCESS: ", cfg.Id)

  // Create directories if they don't already exist
  err := os.MkdirAll(filepath.Join(pm.baseDir, LOG_DIR), 0700)
  if err != nil {
    fmt.Println(err)
  }
  err = os.MkdirAll(filepath.Join(pm.baseDir, PID_DIR), 0700)
  if err != nil {
    fmt.Println(err)
  }

  // if pid file already exist (process already running), don't start the process
  if pm.checkPidFileExists(cfg.PidFilename) {
    pid, err := pm.getPidFromPidFile(cfg.PidFilename)
    if err != nil {
      fmt.Println(err)
    }
    fmt.Println("ProcessManager did not startProcess as an instance with pid=", pid, " is already running")
    return
  }

  // Open the logfile
  logFile, err := os.OpenFile(filepath.Join(pm.baseDir, LOG_DIR, cfg.LogFilename), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
  if err != nil {
    fmt.Println(err)
  }
  defer logFile.Close()

  // Create the command
  fmt.Println(fmt.Sprintf("Going to excute %v, %v", cfg.Command, cfg.Arguments))
  cmd := exec.Command(cfg.Command, cfg.Arguments...)
  cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

  cmd.Stdout = logFile
  cmd.Stderr = logFile
  cmd.Dir = filepath.Join(pm.baseDir, cfg.Dir)

  // Start command
  if err := cmd.Start(); err != nil {
    fmt.Println("ProcessManager:: Failed starting configured process:")
    fmt.Printf("  Tried to run command: %s with arguments: %s \n", cfg.Command, cfg.Arguments)
    fmt.Printf("  ERROR: %s\n\n", err)
    return
  }

  if !background {
    pm.processMap[cfg.Id] = cmd.Process
  }
  // write the process id to file
  content := []byte(strconv.Itoa(cmd.Process.Pid))
  pidFilePath := pm.getPidFilePath(cfg.PidFilename)
  err = ioutil.WriteFile(pidFilePath, content, 0644)
  if err != nil {
    fmt.Printf("ProcessManager Failed writing pid to file", pidFilePath)
  }

  // Blocks until command is done execution
  cmd.Wait()
}

/*
  Stop all processes as defined in the process map by signalling SIGTERM
*/
func (pm *ProcessManager) StopProcesses() {
  for id, process := range pm.processMap {
    fmt.Println(fmt.Sprintf("Stopping process %v with pid: %v ...", id, process.Pid))
    syscall.Kill(-process.Pid, syscall.SIGTERM)

    // clean up the pid files
    for _, config := range pm.processesConfig.Processes {
      pm.removePidFile(config.PidFilename)
    }
  }
}

/*
  Kill a process by name
*/
func (pm *ProcessManager) KillProcess(processName string) {
  fmt.Println("Stopping process", processName)
  _, err := exec.Command("pkill", processName).Output()
  if err != nil {
    fmt.Println(err)
  }

  // clean up the pid file for this process
  for _, config := range pm.processesConfig.Processes {
    if config.Id == processName {
      pm.removePidFile(config.PidFilename)
      break
    }
  }
}

func (pm *ProcessManager) getPidFilePath(filename string) string {
  return filepath.Join(pm.baseDir, PID_DIR, filename)
}

func (pm *ProcessManager) removePidFile(filename string) {
  err := os.RemoveAll(pm.getPidFilePath(filename))
  if err != nil {
    fmt.Printf("ProcessManager failed to remove pid file", filename, err)
  }
}

func (pm *ProcessManager) checkPidFileExists(filename string) bool {
  _, err := os.Stat(pm.getPidFilePath(filename))

  // assume that any error means file doesn't exist
  return err == nil
}

func (pm *ProcessManager) getPidFromPidFile(filename string) (int, error) {
  content, err := ioutil.ReadFile(pm.getPidFilePath(filename))
  if err != nil {
    return 0, err
  }

  pid, err := strconv.Atoi(string(content))
  if err != nil {
    return 0, err
  }

  return pid, nil
}

func (pm *ProcessManager) GetProcessRunningMap() map[string]bool {
  statuses := make(map[string]bool, len(pm.processesConfig.Processes))

  for _, config := range pm.processesConfig.Processes {
    statuses[config.Id] = pm.checkPidFileExists(config.PidFilename)
  }

  return statuses
}

func (pm *ProcessManager) monitorAllProcesses() {
  for _, process := range pm.processesConfig.Processes {
    go pm.monitorProcess(process.PidFilename)
  }
}

func (pm *ProcessManager) monitorProcess(filename string) {
  const SLEEP_DURATION = time.Duration(1) * time.Second
  for {
    if pm.checkPidFileExists(filename) { // pidfile exists

      pid, err := pm.getPidFromPidFile(filename)
      if err == nil { // successfully parsed PID

        // note that sending the null signal is essentially a "dry-run"; no
        // signals are actually sent see kill(2) for more details
        killErr := syscall.Kill(pid, syscall.Signal(0))
        if killErr != nil && killErr == syscall.ESRCH { // process doesn't exist
          pm.removePidFile(filename)
        }
      }
    }

    time.Sleep(SLEEP_DURATION)
  }
}
