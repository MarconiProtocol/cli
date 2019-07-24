package packages

import (
  "fmt"
  "github.com/MarconiProtocol/cli/core/configs"
  "sync"
)

type PackageManager struct {
}

var instance *PackageManager
var once sync.Once

func Instance() *PackageManager {
  once.Do(func() {
    instance = &PackageManager{}
  })
  return instance
}

/*
  Downloads or updates the packages as defined in config
*/
func (pm *PackageManager) UpdatePackages(baseDir string, packageConfigs []configs.PackageConfig) {
  for _, config := range packageConfigs {
    fmt.Println("\nReading package config for", config.Id)
    updatePackage(baseDir, config)
  }
}
