package core

import (
  "fmt"
  "github.com/MarconiProtocol/cli/core/configs"
  "github.com/MarconiProtocol/cli/core/packages"
  "github.com/MarconiProtocol/cli/core/processes"
)

func Bootstrap(baseDir string) {
  fmt.Println("Bootstrapping Marconi Client...")

  // Load the configs
  packages_config := configs.LoadPackagesConf()

  // Check packages
  packages.Instance().UpdatePackages(baseDir, packages_config)
}

func StartProcessManager(baseDir string) {
  processes_config := configs.LoadProcessesConf(baseDir)
  processes.Instance().InitProcessManager(baseDir, *processes_config)
}

func Cleanup() {
  processes.Instance().StopProcesses()
}
