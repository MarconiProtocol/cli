package packages

import (
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
func (pm *PackageManager) UpdatePackages(baseDir string, packagesConfig *configs.PackagesConfig) {
  autoUpdateEnabled := packagesConfig.AutoUpdateEnabled
  for _, config := range packagesConfig.Packages {
    // Check if the package's version file exist as a proxy for whether the package itself exists
    packageExists := doesPackageVersionFileExist(baseDir, config)
    // If the package does not exist, just download
    if !packageExists {
      downloadPackage(baseDir, config)
    } else {
      updatePackage(baseDir, config, autoUpdateEnabled)
    }
  }
}
