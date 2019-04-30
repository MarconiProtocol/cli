package core

import (
  "fmt"
  "gitlab.neji.vm.tc/marconi/cli/core/configs"
  "gitlab.neji.vm.tc/marconi/cli/core/packages"
  "gitlab.neji.vm.tc/marconi/cli/core/processes"
)

func Bootstrap(baseDir string) {
  fmt.Println("Bootstrapping Marconi Client\n")

  // Load the configs
  packages_config := configs.LoadPackagesConf()

  // Check packages
  packages.Instance().UpdatePackages(baseDir, packages_config.Packages)
}

func StartProcessManager(baseDir string) {
  processes_config := configs.LoadProcessesConf(baseDir)
  processes.Instance().InitProcessManager(baseDir, *processes_config)
}

func Cleanup() {
  processes.Instance().StopProcesses()
}
