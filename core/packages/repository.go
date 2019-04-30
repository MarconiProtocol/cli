package packages

import (
  "bufio"
  "fmt"
  "gitlab.neji.vm.tc/marconi/cli/core/configs"
  "os"
  "path/filepath"
  "strings"
)

/*
  Updates the package defined in the provided PackageConfig if necessary
*/
func updatePackage(baseDir string, config configs.PackageConfig) {
  if valid, err := checkConfigValid(baseDir, config); !valid && err != nil {
    handleErr(err)
  }
  if !checkPackageIsUpToDate(baseDir, config) {
    fmt.Println("Package", config.Id, "out of date, updating:\n")
    fetchPackage(baseDir, config)
  } else {
    fmt.Println("Package", config.Id, "is up to date.\n")
  }
}

func checkConfigValid(baseDir string, config configs.PackageConfig) (bool, error) {
  workingDir, err := os.Getwd()
  if err != nil {
    handleErr(err)
  }
  _, err = os.Stat(filepath.Join(baseDir, config.Dir))
  // if  config.Dir doesnt exist, if must not be the current directory
  if err != nil {
    return true, nil
  }
  _, err = os.Stat(workingDir)
  if err != nil {
    handleErr(err)
  }
  return true, nil
}

/*
  Check if the package defined in the provided PackageConfig exists and if it is up to date
*/
func checkPackageIsUpToDate(baseDir string, config configs.PackageConfig) bool {
  // Check for source directory and check version.txt to see if it matches correct version
  versionFilePath := ""
  if config.VersionFile != "" {
    versionFilePath = filepath.Join(baseDir, config.VersionFile)
  } else {
    versionFilePath = filepath.Join(baseDir, config.Dir, "version.txt")
  }
  file, err := os.Open(versionFilePath)
  if err != nil {
    // either package isnt downloaded or version.txt isnt there, either way, the package in our pov doesnt exist
    return false
  }
  defer file.Close()

  // check if the version is correct
  return checkVersionIsAccepted(file, config)
}

/*
  Check the version file and returns if it matches the version defined in config
*/
//TODO: Making a lot of assumptions here about the file format
func checkVersionIsAccepted(file *os.File, config configs.PackageConfig) bool {
  scanner := bufio.NewScanner(file)
  line := ""
  for scanner.Scan() {
    newline := strings.TrimSpace(scanner.Text())
    if len(newline) != 0 {
      if !strings.HasPrefix(newline, "#") {
        line = newline
        break
      }
    }
  }
  versionStrArr := strings.Split(line, "=")
  version := versionStrArr[1]

  // Check if the version defined in the file is accepted
  accepted, err := compareVersions(config.Version, version)
  if err != nil {
    handleErr(err)
  }
  return accepted
}

/*
  Fetch the package defined in the provided PackageConfig
*/
func fetchPackage(baseDir string, config configs.PackageConfig) {
  // Get the target file name based on the source
  sourceArr := strings.Split(config.Source, "/")
  filename := sourceArr[len(sourceArr)-1]
  var dirPath string = filepath.Join(baseDir, config.Dir)
  targetFile := filepath.Join(dirPath, filename)

  // Normally we would completely remove any existing directory
  // containing this package via a call like "removeDir(dirPath)",
  // that way we can start from a clean slate when re-downloading. But
  // because we switched from our old directory structure of
  // individually isolated components to a new structure of
  // overlapping paths (where components end up sharing paths such as
  // /bin and /etc), now there isn't a clean way to remove an old
  // package. Therefore we just re-download and overwrite, even though
  // this is error prone due to statefulness. A less error prone
  // option would be to delete everything for all packages and
  // re-download them all, but that will increase download time and
  // kind of defeats the purpose of having a package manager for
  // managing separate packages in the first place.

  // create the package dir
  createDir(dirPath)
  // download the file
  downloadFileWithHttp(targetFile, config.Source)

  // assume for now we are dealing with tarballs all the time
  getEulaAcknowledgementAndExtractTarball(targetFile, dirPath)

  // cleanup files
  removeFile(targetFile)
}
