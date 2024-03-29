package configs

import (
  "encoding/json"
  "fmt"
  mlog "github.com/MarconiProtocol/log"
  "io/ioutil"
  "os"
)

// Loads a config file, with the assumption that it is a json format, and try to unmarshall it into a BaseConfig object

var (
  // note that the trailing %.s is not a typo; it's needed to suppress the
  // %!(EXTRA string=) that comes from having a Sprintf argument that's unused
  MAIN_CONFIG_FILE      = []string{"%s/configs/mcli.json", "./configs/mcli.json%.s"}
  PACKAGES_CONFIG_FILE  = []string{"%s/configs/packages_conf.json", "./configs/packages_conf.json%.s"}
  PROCESSES_CONFIG_FILE = []string{"%s/configs/processes_conf.json", "./configs/processes_conf.json%.s"}
)

func LoadBaseConf() *BaseConfig {
  fileBytes := loadFileBytes(MAIN_CONFIG_FILE)
  var bootstrapConfig BaseConfig
  err := json.Unmarshal(fileBytes, &bootstrapConfig)
  if err != nil {
    panic(fmt.Sprintf("Failed to parse bootstrap config: %v", err))
  }
  return &bootstrapConfig
}

func LoadPackagesConf() *PackagesConfig {
  fileBytes := loadFileBytes(PACKAGES_CONFIG_FILE)
  var packagesConfig PackagesConfig
  err := json.Unmarshal(fileBytes, &packagesConfig)
  if err != nil {
    panic(fmt.Sprintf("Failed to parse packages config: %v", err))
  }
  return &packagesConfig
}

func UpdatePackagesConf(newPackagesConfig *PackagesConfig) {
  bytes, err := json.MarshalIndent(newPackagesConfig, "", " ")
  if err != nil {
    panic(fmt.Sprintf("Failed to marshal packages config: %v", err))
  }
  writeFileBytes(bytes, PACKAGES_CONFIG_FILE)
}

func LoadProcessesConf(baseDir string) *ProcessesConfig {
  fileBytes := loadFileBytes(PROCESSES_CONFIG_FILE)
  var processesConfig ProcessesConfig
  err := json.Unmarshal(fileBytes, &processesConfig)
  if err != nil {
    panic(fmt.Sprintf("Failed to parse processes config: %v", err))
  }
  return &processesConfig
}

func loadFileBytes(filenames []string) []byte {
  formattedFilenames := make([]string, len(filenames))
  for i, filename := range filenames {
    formattedFilenames[i] = fmt.Sprintf(filename, GetBaseDir())
  }

  var fileBytes []byte
  var err error
  for _, filename := range formattedFilenames {
    // TODO: we need the logger so we can print stuff like this to the logs
    fileBytes, err = ioutil.ReadFile(filename)
    if err == nil && fileBytes != nil {
      break
    }
  }

  if fileBytes == nil || err != nil {
    mlog.GetLogger().Errorf("Failed to open configuration files at any locations: %v", formattedFilenames)
    panic(fmt.Sprintf("Failed to open configuration files at any locations: %v", formattedFilenames))
  }
  return fileBytes
}

func writeFileBytes(bytes []byte, filenames []string) {
  formattedFilenames := make([]string, len(filenames))
  for i, filename := range filenames {
    formattedFilenames[i] = fmt.Sprintf(filename, GetBaseDir())
  }

  var err error
  for _, filename := range formattedFilenames {
    _, err = os.Stat(filename)
    if err == nil {
      err = ioutil.WriteFile(filename, bytes, 0644)
      if err == nil {
        break
      }
    }
  }

  // TODO check the behavior of writeFileBytes and related functions
  if err != nil {
    mlog.GetLogger().Errorf("Failed to write to configuration files at any locations: %v", formattedFilenames)
    //panic(fmt.Sprintf("Failed to write to configuration files at any locations: %v", formattedFilenames))
  }
}
