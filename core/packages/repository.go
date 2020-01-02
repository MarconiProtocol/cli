package packages

import (
  "bufio"
  "errors"
  "fmt"
  "github.com/MarconiProtocol/cli/core/configs"
  "os"
  "path/filepath"
  "strings"
)

/*
  Download the package
*/
func downloadPackage(baseDir string, config configs.PackageConfig) {
  if valid, err := checkConfigValid(baseDir, config); !valid && err != nil {
    handleErr(err)
  }
  packageManifest, err := getPackageManifest(baseDir, config)
  if err != nil {
    handleErr(err)
  }
  fetchPackage(baseDir, config, packageManifest.Source)
}

/*
  Updates the package defined in the provided PackageConfig if necessary
*/
func updatePackage(baseDir string, config configs.PackageConfig, autoUpdateEnabled bool) {
  if valid, err := checkConfigValid(baseDir, config); !valid && err != nil {
    handleErr(err)
  }
  if isUpToDate, packageLocation := checkPackageIsUpToDate(baseDir, config, autoUpdateEnabled); !isUpToDate {
    fmt.Printf("Package %s out of date, updating:\n", config.Id)
    fetchPackage(baseDir, config, packageLocation)
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
  Get the manifest data for a specific package
*/
func getPackageManifest(baseDir string, config configs.PackageConfig) (*configs.PackageManifest, error) {
  // Attempt to grab the manifest file
  manifest, err := downloadManifestWithHttp(config.Manifest)
  if err != nil {
    return nil, err
  }
  return manifest, nil
}

func getVersionFilePath(baseDir string, config configs.PackageConfig) string {
  versionFilePath := ""
  if config.VersionFile != "" {
    versionFilePath = filepath.Join(baseDir, config.VersionFile)
  } else {
    versionFilePath = filepath.Join(baseDir, config.Dir, "version.txt")
  }
  return versionFilePath
}

/*
  Return the version number of a package
*/
func getPackageVersion(baseDir string, config configs.PackageConfig) (string, error) {
  // Check for source directory and check version.txt to see if it matches correct version
  versionFilePath := getVersionFilePath(baseDir, config)
  file, err := os.Open(versionFilePath)
  if err != nil {
    // either package isnt downloaded or version.txt isnt there, either way, the package in our pov doesnt exist
    return "", errors.New(fmt.Sprintf("Failed to find version file for package: %s at %s", config.Id, config.VersionFile))
  }
  defer file.Close()

  // parse version file
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

  return version, nil
}

/*
  Check if the package's version file exists
*/
func doesPackageVersionFileExist(baseDir string, config configs.PackageConfig) bool {
  versionFilePath := getVersionFilePath(baseDir, config)
  _, err := os.Stat(versionFilePath)
  if err != nil && os.IsNotExist(err) {
    return false
  }
  return true
}

/*
  Check if the package defined in the provided PackageConfig exists and if it is up to date
*/
func checkPackageIsUpToDate(baseDir string, config configs.PackageConfig, autoUpdateEnabled bool) (bool, string) {
  manifest, err := getPackageManifest(baseDir, config)
  if err != nil {
    handleErr(err)
  }
  version, err := getPackageVersion(baseDir, config)
  if err != nil {
    handleErr(err)
  }

  // check if the version is correct
  isUpToDate := checkVersionIsAccepted(version, manifest.Version)
  if isUpToDate {
    return true, ""
  }

  // if autoUpdate is not enabled, prompt the user for permission to update
  if !autoUpdateEnabled {
    update, optIn := askForPackageUpdateAcknowledgement(config.Id, version, manifest.Version)
    // If we don't want to update, the package is considered up to date
    if !update {
      return true, ""
    }
    // If the user has decided to opt in, update the packages config
    if optIn {
      packages_config := configs.LoadPackagesConf()
      packages_config.AutoUpdateEnabled = true
      configs.UpdatePackagesConf(packages_config)
    }
  }

  return false, manifest.Source
}

/*
  Check the version file and returns if it matches the version defined in config
*/
func checkVersionIsAccepted(version string, manifestVersion string) bool {
  // Check if the version defined in the file is accepted
  accepted, err := compareVersions(manifestVersion, version)
  if err != nil {
    handleErr(err)
  }
  return accepted
}

/*
  Fetch the package defined in the provided PackageConfig
*/
func fetchPackage(baseDir string, config configs.PackageConfig, packageLocation string) {
  // Get the target file name based on the source
  sourceArr := strings.Split(packageLocation, "/")
  filename := sourceArr[len(sourceArr)-1]
  dirPath := filepath.Join(baseDir, config.Dir)
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
  downloadFileWithHttp(targetFile, packageLocation, true)

  // TODO: temp decrypt the file
  var finalFilename string
  if config.IsEncrypted {
    finalFilename = decryptFile(targetFile)
  } else {
    finalFilename = targetFile
  }

  // assume for now we are dealing with tarballs all the time
  getEulaAcknowledgementAndExtractTarball(finalFilename, dirPath)

  // cleanup files
  removeFile(targetFile)
  if config.IsEncrypted {
    removeFile(finalFilename)
  }
}
