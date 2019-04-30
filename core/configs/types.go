package configs

// Main config - more of a placeholder for now, maybe include modes settings?
type BaseConfig struct {
  Version         string
  MarconiNodeHost string
  MarconiNodePort string
}

// Config for packages to be downloaded
type PackagesConfig struct {
  Version  string
  Packages []PackageConfig
}

type PackageConfig struct {
  Id          string
  Dir         string
  Source      string
  Version     string
  VersionFile string
  IsEncrypted bool
}

// Config for an individual process to be started with the process manager
type ProcessesConfig struct {
  Version   string
  Processes []ProcessConfig
}

type ProcessConfig struct {
  Id                string
  Dependencies      []string
  Dir               string
  Command           string
  Arguments         []string
  LogFilename       string
  WaitForCompletion bool
  WaitTime          int
  PidFilename       string
}
