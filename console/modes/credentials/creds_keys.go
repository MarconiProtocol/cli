package credentials

import (
  "fmt"
  "github.com/MarconiProtocol/go-prompt"
  "github.com/MarconiProtocol/cli/console/modes"
  "github.com/MarconiProtocol/cli/core/mkey"
  "strconv"
)

const (
  USE_MPKEY         = "use"
  GENERATE_MP_KEY   = "generate"
  EXPORT_MP_KEY     = "export"
  LIST_MPKEY_HASHES = "list"
  CANCEL            = "cancel"
)

// Mode suggestions
var METH_KEY_SUGGESTIONS = []prompt.Suggest{
  {Text: GENERATE_MP_KEY, Description: "Generate nodekey"},
  {Text: USE_MPKEY, Description: "Set nodekey to use with other commands"},
  {Text: EXPORT_MP_KEY, Description: "Export nodekey"},
  {Text: LIST_MPKEY_HASHES, Description: "List nodekeys"},
}

/*
  Show prompt suggestions for generate Marconi Node Private key command
*/
func (mm *CredsMode) getGenerateMPKeySuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<0xACCOUNT_ADDRESS>", Description: "The GoMarconi account you wish to generate an nodekey for"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Show prompt suggestions for export Marconi Node Private key command
*/
func (mm *CredsMode) getExportMPKeySuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<0xACCOUNT_ADDRESS>", Description: "The GoMarconi account you wish to export nodekey for"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Show prompt suggestions for export marconi keystore command
*/
func (mm *CredsMode) getExportGMrcKeySuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<0xACCOUNT_ADDRESS>", Description: "The GoMarconi account you wish to export the go Marconi keystore for"}}
  case len(line) == 3:
    return []prompt.Suggest{{Text: "<GO-MARCONI_DATA_DIR_PATH>", Description: "The GoMarconi data directory path"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Show prompt suggestions for generating a 32 bit key
*/
func (mm *CredsMode) getGenerate32BitKeySuggestions(line []string) []prompt.Suggest {
  return []prompt.Suggest{}
}

/*
  Show prompt suggestions for showing mpkeys for an account
*/
func (mm *CredsMode) getListMpkKeyHashesSuggestions(line []string) []prompt.Suggest {
  return []prompt.Suggest{{Text: "<0xACCOUNT_ADDRESS>", Description: "The GoMarconi account whose nodekey you wish to view"}}
}

/*
  Show prompt suggestions for copying mpkeys for use with marconid
*/
func (mm *CredsMode) getUseMpkKeySuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<0xACCOUNT_ADDRESS>", Description: "The GoMarconi account that has the nodekeys"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Handle the  generate Marconi Node Private key command
*/
func (mm *CredsMode) handleGenerateMPKey(args []string) {
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", GENERATE_MP_KEY, "<0xACCOUNT_ADDRESS>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }

  password, cancelled := getPassword("Please enter your account password")
  if cancelled {
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

/*
  Handle the export Marconi Node Private keys command
*/
func (mm *CredsMode) handleExportMPKey(args []string) {
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", EXPORT_MP_KEY, "<0xACCOUNT_ADDRESS>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }

  password, cancelled := getPassword("Please enter your account password")
  if cancelled {
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

func internalListMPKeyHashes(mpkeys []mkey.MpkeyListItem) error {
  fmt.Println("nodeIDs:")
  for _, mpkey := range mpkeys {
    fmt.Printf("%d %48s\n", mpkey.Idx, mkey.AddPrefixPubKeyHash(mpkey.Pubkeyhash))
  }
  return nil
}

func (mm *CredsMode) handleListMPKeyHashes(args []string) {
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

func containsMPKey(mpkeys []mkey.MpkeyListItem, keyName string) bool {
  for _, mpkey := range mpkeys {
    if strconv.Itoa(mpkey.Idx) == keyName {
      return true
    }
  }
  return false
}

func (mm *CredsMode) handleUseMPKey(args []string) {
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", USE_MPKEY, "<0xACCOUNT_ADDRESS>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }

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

  fmt.Printf("\nPlease 'Enter' to use nodekey 0 or specify a different one (e.g. '1', '2'):")
  keyInput := ""
  for {
    keyInput = prompt.Input("", func(document prompt.Document) []prompt.Suggest {
      return []prompt.Suggest{}
    })
    if keyInput != "" && keyInput != CANCEL && !containsMPKey(mpkeys, keyInput) {
      fmt.Println("Please enter a nodekey or", CANCEL ,"to cancel")
    } else {
      break
    }
  }
  if keyInput == CANCEL {
    fmt.Println("Cancelled")
  } else {
    if keyInput == "" {
      keyInput = "0"
    }
    keystore, err := mkey.GetAccountForAddress(accountAddress)
    if err != nil {
      fmt.Println(err)
      return
    }
    password, cancelled := getPassword("Please enter your account password")
    if cancelled {
      return
    }
    keyName := mkey.MARCONI_PRIVATE_KEY_FILENAME + keyInput
    err = keystore.UseMarconiKey(keyName, password)
    if err != nil {
      fmt.Println("Failed to use nodekey", err)
      return
    }
  }
}
