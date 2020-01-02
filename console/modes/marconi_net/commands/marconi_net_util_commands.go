package marconi_net_commands

import (
  "fmt"
  "github.com/MarconiProtocol/cli/api/middleware"
  "github.com/MarconiProtocol/cli/console/execution/execution_flags"
  "github.com/MarconiProtocol/cli/console/modes"
  "github.com/MarconiProtocol/cli/console/util"
  "github.com/MarconiProtocol/cli/core/configs"
  "github.com/MarconiProtocol/cli/core/mkey"
  "github.com/MarconiProtocol/marconid/util"
  "github.com/MarconiProtocol/go-prompt"
  "os"
  "path/filepath"
  "strings"
)

// Commands
const (
  GENERATE_32BITKEY = "generate_32bitkey"
  REGISTER          = "register"
  GET_MPIPE_PORT    = "get_mpipe_port"
  START_NETFLOW    = "start_netflow"

  DEFAULT_KEY_CHILD_PATH = "/etc/marconid/"
  NETFLOW_DIR = "/opt/marconi/var/log/netflow"
)

var UTIL_COMMAND_MAP = map[string]func([]string){
  GENERATE_32BITKEY: Generate32BitKey,
  REGISTER:          Register,
  GET_MPIPE_PORT:    GetMPipePort,
  START_NETFLOW:    StartNetflow,
}

func HandleUtilCommand(args []string) {
  if !modes.ArgsMinLenCheck(args, 1) {
    fmt.Println("USAGE: <command>")
    return
  }
  commandType := args[0]
  commandArgs := args[1:]
  if commandHandlerFunction, present := UTIL_COMMAND_MAP[commandType]; present {
    commandHandlerFunction(commandArgs)
  } else {
    fmt.Println("Invalid command " + commandType)
  }
}

func Generate32BitKey(args []string) {
  if !modes.ArgsLenCheckWithOptionalRange(args, 0, 1, 3) {
    fmt.Println("Usage:", GENERATE_32BITKEY, "[Optional:", execution_flags.PATH, "<PATH> |", execution_flags.SKIP_PROMPT_USE_DEFAULTS, "]")
    return
  }

  defaultPath := configs.GetFullPath(DEFAULT_KEY_CHILD_PATH)
  keyfileName := "l2.key"
  fmt.Println("Generating base64 32bit key.")

  executionFlags := execution_flags.NewExecFlags(args)
  path := ""

  // Checking if execution flag has a pathFlag value
  if executionFlags.CheckPathFlagSet() || executionFlags.CheckSkipPromptsFlagSet() {
    if executionFlags.GetPath() != "''" {
      path = executionFlags.GetPath()
    } else {
      fmt.Println("No path specified. Using the Default:", defaultPath)
      path = defaultPath
    }
  } else {
    fmt.Printf("Enter path in which to save the key (%s):\n", defaultPath)
    path = prompt.Input("", func(document prompt.Document) []prompt.Suggest {
      return []prompt.Suggest{}
    })
  }
  // use a default if nothing was entered
  path = strings.TrimSpace(path)
  if path == "" {
    path = defaultPath
  }

  // check if an existing key exists
  keyfilePath := filepath.Join(path, keyfileName)
  if _, err := os.Stat(keyfilePath); !os.IsNotExist(err) {
    if !executionFlags.CheckPathFlagSet() && !executionFlags.CheckSkipPromptsFlagSet() {
      fmt.Println(keyfilePath, "exists. Overwrite?")
      confirmed := modes.GetConfirmationInput()
      if !confirmed {
        fmt.Println("Skipping generation")
        return
      }
    } else {
      fmt.Println(keyfilePath, "exists, Overwriting.")
    }
  }
  _, err := mkey.Generate32BitKey(keyfilePath)
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Key successfully generated and written to", keyfilePath)
  }
}

func Register(args []string) {
  if !checkMiddlewareRunning() {
    return
  }
  if !modes.ArgsLenCheck(args, 2) {
    fmt.Println("Usage:", REGISTER, "<NODE_ID>, <MAC_HASH>")
    return
  }
  if !modes.ArgPubKeyHashCheck(args[0]) {
    return
  }

  result, err := middleware.GetClient().RegisterUser(args[0], args[1])
  if err != nil {
    fmt.Println("Error:", err)
  } else {
    fmt.Println("Registered a peer to the network:")
    fmt.Printf("%-24s %48s\n\n", "Registered Peer", mkey.AddPrefixPubKeyHash(result.PubKeyHash))
  }
}

func GetMPipePort(args []string) {
  if !modes.ArgsLenCheck(args, 2) {
    fmt.Println("Usage:", GET_MPIPE_PORT, "<PEER_NODE_ID>, <OTHER_PEER_NODE_ID>")
    return
  }
  if !modes.ArgPubKeyHashCheck(args[0]) ||
    !modes.ArgPubKeyHashCheck(args[1]) {
    return
  }
  mutualPort := mutil.GetMutualMPipePort(strings.TrimPrefix(args[0], util.PUB_KEY_PREFIX), strings.TrimPrefix(args[1], util.PUB_KEY_PREFIX))
  fmt.Printf("The MPipe port is: %d\n", mutualPort)
}

func StartNetflow(args []string) {
  if !checkMiddlewareRunning() {
    return
  }
  if !modes.ArgsLenCheckWithOptional(args, 3, 4) {
    fmt.Println("Usage:", START_NETFLOW, "<COLLECTOR_IP>, <COLLECTOR_PORT>, <INTERFACE_ID>, [LOGGING_DIR]")
    return
  }

  var directory string = NETFLOW_DIR
  if len(args) > 3 {
    directory = args[3]
    if args[0] != "127.0.0.1" && args[0] != "localhost" {
      fmt.Println("Cannot use logging directory if not collected on localhost.")
      return
    }
  }

  err := middleware.GetClient().StartNetflow(args[0], args[1], args[2], directory)
  if err != nil {
    fmt.Println("Error:", err)
  } else {
    fmt.Println("Started netflow monitoring")
  }
}
