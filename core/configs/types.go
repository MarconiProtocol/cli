package configs

type BaseConfig struct {
  Version         string
  MarconiNodeHost string
  MarconiNodePort string
  MarconidRPCPort string
}

// Config for packages to be downloaded
type PackagesConfig struct {
  Version           string
  AutoUpdateEnabled bool
  Packages          []PackageConfig
}

type PackageConfig struct {
  Id          string
  Dir         string
  VersionFile string
  IsEncrypted bool
  Manifest    string
}

type PackageManifest struct {
  Version  string
  Source   string
  Checksum string
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
