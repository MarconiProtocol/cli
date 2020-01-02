package credential_commands

import (
  "fmt"
  "github.com/MarconiProtocol/cli/console/execution/execution_flags"
  "github.com/MarconiProtocol/cli/console/modes"
  "github.com/MarconiProtocol/go-prompt"
  "io/ioutil"
  "strings"
)

// Commands
const (
  ACCOUNT = "account"
  KEY     = "key"
)

var COMMAND_MAP = map[string]func([]string){
  ACCOUNT: HandleAccountCommand,
  KEY:     HandleKeyCommand,
}

type jsonObject = map[string]interface{}

func getPasswordFromFlags(ef *execution_flags.ExecFlags) (string, bool, error) {

  // Priority: Check for Password File then Password
  if ef.CheckPasswordFileFlagSet() {
    // Try to open file and retrieve password from there
    data, err := ioutil.ReadFile(ef.GetPasswordFile())

    // Check if an error occured
    if err != nil {
      fmt.Println(err)
      return "", true, err
    }

    // File read was a success
    return strings.TrimSpace(string(data)), true, nil
  }

  // Checking if the password flag is set
  if ef.CheckPasswordFlagSet() {
    password := ef.GetPassword()
    // If the password is the default, then it is empty
    if password == "''" {
      password = ""
    }
    return password, true, nil
  }
  return "", false, nil
}

func getPassword(outputPrompt string, ef *execution_flags.ExecFlags) (string, bool, error) {

  // First see if password was provided in arguments
  password, found, err := getPasswordFromFlags(ef)

  if found == true {
    return password, false, err
  }

  // If not continue with usual prompt
  fmt.Print(outputPrompt, ": ")

  earlyExit := make(chan struct{}, 1)
  cancelled := false

  cancel := func(*prompt.Buffer) {
    cancelled = true
    earlyExit <- struct{}{}
  }

  password = prompt.Input("", modes.PasswordCompleter,
    prompt.OptionHiddenInput(),
    prompt.OptionAddKeyBind(prompt.KeyBind{Key: prompt.ControlC, Fn: cancel}),
    prompt.OptionSetEarlyExit(earlyExit))
  return password, cancelled, err
}
