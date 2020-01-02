package execution_flags

const (
  PATH                     = "--path"
  PASSWORD                 = "--password"
  PASSWORD_FILE            = "--password-file"
  NODE_KEY                 = "--node-key"
  SKIP_PROMPT_USE_DEFAULTS = "--skip-prompts"
)

var execFlagsMap = map[string]string{
  PATH:                     "''",
  PASSWORD:                 "''",
  PASSWORD_FILE:            "''",
  NODE_KEY:                 "0",
  SKIP_PROMPT_USE_DEFAULTS: "''",
}

type ExecFlags struct {
  path         string
  password     string
  passwordFile string
  nodeKey      string
  skipPrompts  bool
}

func NewExecFlags(args []string) *ExecFlags {
  execFlags := ExecFlags{}
  for index, value := range args {
    if index+1 < len(args) {
      execFlags.parseFlag(value, args[index+1])
    } else { //last value
      execFlags.parseFlag(value, value)
    }
  }
  return &execFlags
}

func (ef *ExecFlags) setFlagValue(flag string, value string) {
  switch flag {
  case PATH:
    ef.path = value
  case PASSWORD:
    ef.password = value
  case PASSWORD_FILE:
    ef.passwordFile = value
  case NODE_KEY:
    ef.nodeKey = value
  case SKIP_PROMPT_USE_DEFAULTS:
    ef.skipPrompts = true
  }
}

func (ef *ExecFlags) parseFlag(flag string, value string) {
  if defaultValue, flagPresent := execFlagsMap[flag]; flagPresent {
    // Checking if the second 'value' is in fact a flag
    if _, valuePresent := execFlagsMap[value]; valuePresent {
      ef.setFlagValue(flag, defaultValue)
    } else {
      ef.setFlagValue(flag, value)
    }
  }
}

func (ef *ExecFlags) CheckSkipPromptsFlagSet() bool {
  return ef.skipPrompts != false
}

func (ef *ExecFlags) CheckPathFlagSet() bool {
  return ef.path != ""
}

func (ef *ExecFlags) GetPath() string {
  return ef.path
}

func (ef *ExecFlags) CheckPasswordFlagSet() bool {
  return ef.password != ""
}

func (ef *ExecFlags) GetPassword() string {
  return ef.password
}

func (ef *ExecFlags) CheckPasswordFileFlagSet() bool {
  return ef.passwordFile != ""
}

func (ef *ExecFlags) GetPasswordFile() string {
  return ef.passwordFile
}

func (ef *ExecFlags) CheckNodeKeyFlagSet() bool {
  return ef.nodeKey != ""
}

func (ef *ExecFlags) GetNodeKey() string {
  return ef.nodeKey
}
