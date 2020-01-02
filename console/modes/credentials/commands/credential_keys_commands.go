package credential_commands

import (
  "fmt"
  "github.com/MarconiProtocol/cli/console/execution/execution_flags"
  "github.com/MarconiProtocol/cli/console/modes"
  "github.com/MarconiProtocol/cli/core/mkey"
  "github.com/MarconiProtocol/go-prompt"
  "strconv"
)

// Commands
const (
  USE_MPKEY         = "use"
  GENERATE_MP_KEY   = "generate"
  EXPORT_MP_KEY     = "export"
  LIST_MPKEY_HASHES = "list"
  CANCEL            = "cancel"
)

var KEY_COMMAND_MAP = map[string]func([]string){
  USE_MPKEY:         UseMPKey,
  GENERATE_MP_KEY:   GenerateMPKey,
  EXPORT_MP_KEY:     ExportMPKey,
  LIST_MPKEY_HASHES: ListMPKeyHashes,
}

func HandleKeyCommand(args []string) {
  if !modes.ArgsMinLenCheck(args, 1) {
    fmt.Println("USAGE: <command>")
    return
  }
  commandType := args[0]
  commandArgs := args[1:]
  if commandHandlerFunction, present := KEY_COMMAND_MAP[commandType]; present {
    commandHandlerFunction(commandArgs)
  } else {
    fmt.Println("Invalid command " + commandType)
  }
}

func GenerateMPKey(args []string) {
  if !modes.ArgsLenCheckWithOptional(args, 1, 2) {
    fmt.Println("Usage:", GENERATE_MP_KEY, "<0xACCOUNT_ADDRESS> [Optional:", execution_flags.PASSWORD, "<password> ", execution_flags.PASSWORD_FILE, "<passwordfile>]")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }

  executionFlags := execution_flags.NewExecFlags(args)

  password, cancelled, err := getPassword("Please enter your account password", executionFlags)
  if cancelled || err != nil {
    return
  }

  keystore, err := mkey.GetAccountForAddress(args[0])
  if err != nil {
    fmt.Println(err)
    return
  }

  // Do hacky validation of password by calling GetGoMarconiKey which
  // tries to decrypt it.
  _, err = keystore.GetGoMarconiKey(password)
  if err != nil {
    fmt.Println("Failed to validate password:", err)
    return
  }

  marconiKey, err := keystore.GenerateMarconiKey(password)
  if err != nil {
    fmt.Println("Failed to generate nodekey", err)
    return
  }
  fmt.Println("nodeID:")
  fmt.Println(mkey.AddPrefixPubKeyHash(marconiKey.PublicKeyHash))
  err = keystore.ExportMarconiKeys(password)
  if err != nil {
    fmt.Println("Failed to export nodekeys", err)
    return
  }

  // use the newly generated key
  index := len(keystore.MarconiKeys) - 1
  keyName := mkey.MARCONI_PRIVATE_KEY_FILENAME + strconv.Itoa(index)
  err = keystore.UseMarconiKey(keyName, password)
  if err != nil {
    fmt.Println("Failed to use nodekey", keyName, err)
  }
}

func UseMPKey(args []string) {
  if !modes.ArgsLenCheckWithOptionalRange(args, 1, 2, 5) {
    fmt.Println("Usage:", USE_MPKEY, "<0xACCOUNT_ADDRESS> [Optional:", execution_flags.PASSWORD, "<password> |", execution_flags.PASSWORD_FILE, "<password file> |", execution_flags.NODE_KEY, "<node key> (default 0) |", execution_flags.SKIP_PROMPT_USE_DEFAULTS, "]")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }

  executionFlags := execution_flags.NewExecFlags(args)

  accountAddress := args[0]
  mpkeys, err := mkey.GetMPKeyHashesForAddress(accountAddress)
  if err != nil {
    fmt.Println("Failed to retrieve NodeKey for the account", accountAddress, err)
    return
  } else if len(mpkeys) == 0 {
    fmt.Println("There are no nodekeys for the account", accountAddress)
    return
  }

  fmt.Println("nodekeys associated with this account:")
  err = internalListMPKeyHashes(mpkeys)
  if err != nil {
    return
  }

  nodeKey := ""

  // Check if value is provided by a flag, if not ask user for it
  if executionFlags.CheckNodeKeyFlagSet() || executionFlags.CheckSkipPromptsFlagSet() {

    // If skip prompts is set but the flag isn't set it to the default 0
    if !executionFlags.CheckNodeKeyFlagSet() {
      nodeKey = "0"
    } else {
      nodeKey = executionFlags.GetNodeKey()
    }

    // Checking if it is valid
    if !containsMPKey(mpkeys, nodeKey) {
      fmt.Println("Invalid node key")
      return
    }

    fmt.Println("Using node key:", nodeKey)

  } else { // Value is not provided by a flag and skip prompts is not set
    fmt.Printf("\nPlease 'Enter' to use nodekey 0 or specify a different one (e.g. '1', '2'):")
    keyInput := ""
    for {
      keyInput = prompt.Input("", func(document prompt.Document) []prompt.Suggest {
        return []prompt.Suggest{}
      })
      if keyInput != "" && keyInput != CANCEL && !containsMPKey(mpkeys, keyInput) {
        fmt.Println("Please enter a nodekey or", CANCEL, "to cancel")
      } else {
        break
      }
    }

    if keyInput == CANCEL {
      fmt.Println("Cancelled")
      return
    } else {
      if keyInput == "" {
        keyInput = "0"
      }
    }

    nodeKey = keyInput
  }

  keystore, err := mkey.GetAccountForAddress(accountAddress)
  if err != nil {
    fmt.Println(err)
    return
  }
  password, cancelled, err := getPassword("Please enter your account password", executionFlags)
  if cancelled || err != nil {
    return
  }

  keyName := mkey.MARCONI_PRIVATE_KEY_FILENAME + nodeKey
  err = keystore.UseMarconiKey(keyName, password)
  if err != nil {
    fmt.Println("Failed to use nodekey", err)
    return
  }
}

func ExportMPKey(args []string) {
  if !modes.ArgsLenCheckWithOptional(args, 1, 2) {
    fmt.Println("Usage:", EXPORT_MP_KEY, "<0xACCOUNT_ADDRESS> [Optional:", execution_flags.PASSWORD, "<password> |", execution_flags.PASSWORD_FILE, "<password file> ]")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }

  executionFlags := execution_flags.NewExecFlags(args)

  password, cancelled, err := getPassword("Please enter your account password", executionFlags)
  if cancelled || err != nil {
    return
  }

  keystore, err := mkey.GetAccountForAddress(args[0])
  if err != nil {
    fmt.Println(err)
    return
  }
  err = keystore.ExportMarconiKeys(password)
  if err != nil {
    fmt.Println("Failed to export nodekey", err)
    return
  }

  return
}

func containsMPKey(mpkeys []mkey.MpkeyListItem, keyName string) bool {
  for _, mpkey := range mpkeys {
    if strconv.Itoa(mpkey.Idx) == keyName {
      return true
    }
  }
  return false
}

func internalListMPKeyHashes(mpkeys []mkey.MpkeyListItem) error {
  fmt.Println("nodeIDs:")
  for _, mpkey := range mpkeys {
    fmt.Printf("%d %48s\n", mpkey.Idx, mkey.AddPrefixPubKeyHash(mpkey.Pubkeyhash))
  }
  return nil
}

func ListMPKeyHashes(args []string) {
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", LIST_MPKEY_HASHES, "<0xACCOUNT_ADDRESS>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }

  mpkeys, err := mkey.GetMPKeyHashesForAddress(args[0])
  if err != nil {
    fmt.Println("Failed to retrieve nodekeys for the account", args[0], err)
    return
  }
  internalListMPKeyHashes(mpkeys)
}
