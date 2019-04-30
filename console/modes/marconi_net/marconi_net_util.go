package marconi_net

import (
  "fmt"
  "gitlab.neji.vm.tc/marconi/go-prompt"
  "gitlab.neji.vm.tc/marconi/cli/console/modes"
  "gitlab.neji.vm.tc/marconi/cli/core/configs"
  "gitlab.neji.vm.tc/marconi/cli/core/mkey"
  "os"
  "path/filepath"
  "strings"
)

const (
  GENERATE_32BITKEY = "generate_32bitkey"
  REGISTER          = "register"

  DEFAULT_KEY_CHILD_PATH = "/etc/marconid/"
)

// Mode suggestions
var MNET_UTIL_SUGGESTIONS = []prompt.Suggest{
  {Text: GENERATE_32BITKEY, Description: "Generate a 32 bit key"},
  {Text: REGISTER, Description: "Register a nodeID"},
}

func (mnm *MarconiNetMode) getGenerate32BitKeySuggestions(line []string) []prompt.Suggest {
  return []prompt.Suggest{}
}

/*
  Show prompt suggestions for register command
*/
func (mnm *MarconiNetMode) getRegisterUserSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<NODE_ID>", Description: "Hash of nodekey"}}
  case len(line) == 3:
    return []prompt.Suggest{{Text: "<MAC_HASH>", Description: "Hash of MAC address"}}
  }
  return []prompt.Suggest{}
}

func (mnm *MarconiNetMode) handleGenerate32BitKey(args []string) {
  defaultPath := configs.GetFullPath(DEFAULT_KEY_CHILD_PATH)
  keyfileName := "l2.key"
  fmt.Println("Generating base64 32bit key.")
  fmt.Printf("Enter path in which to save the key (%s):\n", defaultPath)
  path := prompt.Input("", func(document prompt.Document) []prompt.Suggest {
    return []prompt.Suggest{}
  })
  // use a default if nothing was entered
  path = strings.TrimSpace(path)
  if path == "" {
    path = defaultPath
  }

  // check if an existing key exists
  keyfilePath := filepath.Join(path, keyfileName)
  if _, err := os.Stat(keyfilePath); !os.IsNotExist(err) {
    fmt.Println(keyfilePath, "exists. Overwrite?")
    confirmed := modes.GetConfirmationInput()
    if !confirmed {
      fmt.Println("Skipping generation")
      return
    }
  }
  _, err := mkey.Generate32BitKey(keyfilePath)
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Key successfully generated")
  }
}

/*
  Handle the register command
*/
func (mnm *MarconiNetMode) handleRegisterUser(args []string) {
  if !mnm.checkMiddlewareRunning() {
    return
  }
  if !modes.ArgsLenCheck(args, 2) {
    fmt.Println("Usage:", REGISTER, "<NODE_ID>, <MAC_HASH>")
    return
  }
  if !modes.ArgPubKeyHashCheck(args[0]) {
    return
  }

  result, err := mnm.middlewareClient.RegisterUser(args[0], args[1])
  if err != nil {
    fmt.Println("Error:", err)
  } else {
    fmt.Println("Registered a peer to the network:")
    fmt.Printf("%-24s %48s\n\n", "Registered Peer", mkey.AddPrefixPubKeyHash(result.PubKeyHash))
  }
}
