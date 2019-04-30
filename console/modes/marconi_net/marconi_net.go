package marconi_net

import (
  "fmt"
  "gitlab.neji.vm.tc/marconi/go-prompt"
  "gitlab.neji.vm.tc/marconi/cli/api/middleware"
  "gitlab.neji.vm.tc/marconi/cli/console/context"
  "gitlab.neji.vm.tc/marconi/cli/console/modes"
  "gitlab.neji.vm.tc/marconi/cli/console/util"
  "gitlab.neji.vm.tc/marconi/cli/core/configs"
  "gitlab.neji.vm.tc/marconi/cli/core/mkey"
  "gitlab.neji.vm.tc/marconi/cli/core/processes"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "strconv"
  "strings"
)

const (
  PEER             = "peer"
  UTIL             = "util"
  USE              = "use"
  CREATE_NETWORK   = "create"
  DELETE_NETWORK   = "delete"
  JOIN_NETWORK     = "join"
  GET_NETWORK_INFO = "info"
)

const (
  MARCONID_CONFIG_CHILD_PATH = "/etc/marconid/config.yml"
  EMPTY_CONTRACT_ADDRESS = "0x0000000000000000000000000000000000000000"

  NETWORK_INFO_FORMAT = "%-24s %48s\n"
)

// Mode suggestions
var MNET_SUGGESTIONS = []prompt.Suggest{
  {Text: PEER, Description: "Peer related commands"},
  {Text: UTIL, Description: "Utility commands"},
  {Text: USE, Description: "Set network to use with other commands"},
  {Text: CREATE_NETWORK, Description: "Create new network"},
  {Text: DELETE_NETWORK, Description: "Delete existing network"},
  {Text: JOIN_NETWORK, Description: "Join an existing network"},
  {Text: GET_NETWORK_INFO, Description: "Get network info"},
  {Text: modes.RETURN_TO_ROOT, Description: "Return to home menu"},
  {Text: modes.EXIT_CMD, Description: "Exit mcli"},
}

/*
  Menu/Mode for Marconi Net related operations
*/
type MarconiNetMode struct {
  modes.BaseMode
  middlewareClient *middleware.Client
  contractAddress  string
}

/*
  Create a new Marconi Net cmd struct,
  use this to create a MarconiNetMode otherwise suggestion and handlers wont be properly initialized
*/
func NewMarconiNetMode(c *context.Context) *MarconiNetMode {
  mnetMode := MarconiNetMode{}
  mnetMode.Init(c)
  mnetMode.SetBaseSuggestions(MNET_SUGGESTIONS)

  // suggestion and handler registrations
  mnetMode.RegisterCommand(USE, mnetMode.getUseNetworkSuggestions, mnetMode.handleUseNetwork)
  mnetMode.RegisterCommand(CREATE_NETWORK, mnetMode.getCreateNetworkSuggestions, mnetMode.handleCreateNetwork)
  mnetMode.RegisterCommand(DELETE_NETWORK, mnetMode.getDeleteNetworkSuggestions, mnetMode.handleDeleteNetwork)
  mnetMode.RegisterCommand(JOIN_NETWORK, mnetMode.getJoinNetworkSuggestions, mnetMode.handleJoinNetwork)
  mnetMode.RegisterCommand(GET_NETWORK_INFO, mnetMode.getGetNetworkInfoSuggestions, mnetMode.handleGetNetworkInfo)

  mnetMode.RegisterCommand(PEER, mnetMode.getPeerSuggestions, mnetMode.handlePeer)
  mnetMode.RegisterSubCommand(PEER, ADD_PEER, mnetMode.getAddPeerSuggestions, mnetMode.handleAddPeer)
  mnetMode.RegisterSubCommand(PEER, REMOVE_PEER, mnetMode.getRemovePeerSuggestions, mnetMode.handleRemovePeer)
  mnetMode.RegisterSubCommand(PEER, ADD_PEER_RELATION, mnetMode.getAddPeerRelationSuggestions, mnetMode.handleAddPeerRelation)
  mnetMode.RegisterSubCommand(PEER, REMOVE_PEER_RELATION, mnetMode.getRemovePeerRelationSuggestions, mnetMode.handleRemovePeerRelation)
  mnetMode.RegisterSubCommand(PEER, GET_PEER_RELATIONS, mnetMode.getGetPeerRelationsSuggestions, mnetMode.handleGetPeerRelations)
  mnetMode.RegisterSubCommand(PEER, GET_PEER_INFO, mnetMode.getGetPeerInfoSuggestions, mnetMode.handleGetPeerInfo)

  mnetMode.RegisterCommand(UTIL, mnetMode.getUtilSuggestions, mnetMode.handleUtil)
  mnetMode.RegisterSubCommand(UTIL, GENERATE_32BITKEY, mnetMode.getGenerate32BitKeySuggestions, mnetMode.handleGenerate32BitKey)
  mnetMode.RegisterSubCommand(UTIL, REGISTER, mnetMode.getRegisterUserSuggestions, mnetMode.handleRegisterUser)

  mnetMode.RegisterCommand(modes.RETURN_TO_ROOT, mnetMode.GetEmptySuggestions, mnetMode.HandleReturnToRoot)
  mnetMode.RegisterCommand(modes.EXIT_CMD, mnetMode.GetEmptySuggestions, mnetMode.HandleExitCommand)

  mnetMode.middlewareClient = middleware.GetClient()

  return &mnetMode
}

func (mnm *MarconiNetMode) CliPrefix() (string, bool) {
  var networkPrefix string
  if mnm.contractAddress != "" {
    networkPrefix = " <..." + mnm.contractAddress[len(mnm.contractAddress)-8:] + ">"
  }
  return mnm.Name() + networkPrefix, true
}

func (mnm *MarconiNetMode) Name() string {
  return "net"
}

/*
  Show prompt suggestions for entering the peer sub menu
*/
func (mnm *MarconiNetMode) getPeerSuggestions(line []string) []prompt.Suggest {
  return util.SimpleSubcommandCompleter(line, 1, MNET_PEER_SUGGESTIONS)
}

/*
  Show prompt suggestions for entering the util sub menu
*/
func (mnm *MarconiNetMode) getUtilSuggestions(line []string) []prompt.Suggest {
  return util.SimpleSubcommandCompleter(line, 1, MNET_UTIL_SUGGESTIONS)
}

/*
  Show prompt suggestions for use network command
*/
func (mnm *MarconiNetMode) getUseNetworkSuggestions(line []string) []prompt.Suggest {
  if len(line) == 2 {
    return []prompt.Suggest{{Text: "0x<NETWORK_CONTRACT_ADDRESS>", Description: "Network ID"}}
  }
  return []prompt.Suggest{}
}

/*
  Show prompt suggestions for create network command
*/
func (mnm *MarconiNetMode) getCreateNetworkSuggestions(line []string) []prompt.Suggest {
  return []prompt.Suggest{}
}

/*
  Show prompt suggestions for delete network command
*/
func (mnm *MarconiNetMode) getDeleteNetworkSuggestions(line []string) []prompt.Suggest {
  if len(line) == 2 {
    return []prompt.Suggest{{Text: "0x<NETWORK_CONTRACT_ADDRESS>", Description: "Network ID"}}
  }
  return []prompt.Suggest{}
}

/*
  Show prompt suggestions for join network command
*/
func (mnm *MarconiNetMode) getJoinNetworkSuggestions(line []string) []prompt.Suggest {
  if len(line) == 2 {
    return []prompt.Suggest{{Text: "0x<NETWORK_CONTRACT_ADDRESS>", Description: "Network ID"}}
  }
  return []prompt.Suggest{}
}

/*
  Show prompt suggestions for get network info command
*/
func (mnm *MarconiNetMode) getGetNetworkInfoSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<0xNETWORK_CONTRACT_ADDRESS>", Description: "The address of the smart contract for the network in which to inspect"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Handle the peers sub menu command
*/
func (mnm *MarconiNetMode) handlePeer(args []string) {
  util.HandleFurtherCommands(PEER, MNET_PEER_SUGGESTIONS)
}

/*
  Handle the util sub menu command
*/
func (mnm *MarconiNetMode) handleUtil(args []string) {
  util.HandleFurtherCommands(UTIL, MNET_UTIL_SUGGESTIONS)
}

/*
  Handle the use network command
*/
func (mnm *MarconiNetMode) handleUseNetwork(args []string) {
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", USE, "<NETWORK_CONTRACT_ADDRESS>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }
  mnm.contractAddress = util.GetEIP55Address(args[0])
  fmt.Println("Using network contract")
  fmt.Println(mnm.contractAddress)
}

func (mnm *MarconiNetMode) printNetworkInfo(action, networkId, admin, networkContract string) {
  fmt.Printf("%s:\n", action)

  // we will print this out later, but hide for now
  //if networkId != "" {
  //  fmt.Printf(NETWORK_INFO_FORMAT, "Network Id", networkId)
  //}
  if admin != "" {
    fmt.Printf(NETWORK_INFO_FORMAT, "Admin Account", util.GetEIP55Address(admin))
  }
  if networkContract != "" {
    fmt.Printf(NETWORK_INFO_FORMAT, "Network Contract Address", util.GetEIP55Address(networkContract))
  }
}

/*
  Handle the create network command
*/
func (mnm *MarconiNetMode) handleCreateNetwork(args []string) {
  if !mnm.checkMiddlewareRunning() {
    return
  }
  if !modes.ArgsLenCheckWithOptional(args, 0, 1) {
    fmt.Println("Usage:", CREATE_NETWORK)
    return
  }

  fmt.Println("This may take up to a minute...")
  result, err := mnm.middlewareClient.CreateNetwork()
  if err != nil {
    fmt.Println("Failed to create network:", err)
  } else {
    mnm.printNetworkInfo("Created a new network", result.NetworkId, result.Admin, result.NetworkContract)

    // use the network is the optional argument is provided
    if len(args) == 1 && args[0] == "use" {
      mnm.handleUseNetwork([]string{result.NetworkContract})
    }
  }
}

/*
  Handle the delete network command
*/
func (mnm *MarconiNetMode) handleDeleteNetwork(args []string) {
  if !mnm.checkMiddlewareRunning() {
    return
  }
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", DELETE_NETWORK, "<0xNETWORK_CONTRACT_ADDRESS>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }
  result, err := mnm.middlewareClient.DeleteNetwork(args[0])
  if err != nil {
    fmt.Println("Error:", err)
  } else {
    mnm.printNetworkInfo("Deleted network", result.NetworkId, result.Admin, "")
  }
}

/*
  Handle the delete network command
*/
func (mnm *MarconiNetMode) handleJoinNetwork(args []string) {
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", JOIN_NETWORK, "<0xNETWORK_CONTRACT_ADDRESS>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }
  // update config.yml with the Network Address
  confPath := configs.GetFullPath(MARCONID_CONFIG_CHILD_PATH)
  data, err := ioutil.ReadFile(confPath)

  configMap := make(map[interface{}]interface{})
  err = yaml.Unmarshal(data, &configMap)
  if err != nil {
    fmt.Println("Failed to parse", confPath, err)
    return
  }
  blockchain := configMap["blockchain"].(map[interface{}]interface{})
  networkAddress := blockchain["networkContractAddress"]
  if err == nil {
    if networkAddress == EMPTY_CONTRACT_ADDRESS {
      blockchain["networkContractAddress"] = args[0]
      newContents, err := yaml.Marshal(configMap)
      if err == nil {
        err = ioutil.WriteFile(confPath, []byte(newContents), 0644)
        if err == nil {
          fmt.Println("Joined network")
        } else {
          fmt.Println("Failed to join network", err)
        }
      } else {
        fmt.Println("Failed to marshal config object", err)
      }
    } else {
      fmt.Println("Join network skipped, user is already part of network", networkAddress)
    }
  }
}

/*
  Handle the get network info command
*/
func (mnm *MarconiNetMode) handleGetNetworkInfo(args []string) {
  if !mnm.checkMiddlewareRunning() {
    return
  }
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", GET_NETWORK_INFO, "<0xNETWORK_CONTRACT_ADDRESS>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }

  networkAddress := args[0]
  result, err := mnm.middlewareClient.GetNetworkInfo(networkAddress)
  if err != nil {
    fmt.Println("Error:", err)
  } else {
    mnm.printNetworkInfo("Network Info", result.NetworkId, result.NetworkAdmin, networkAddress)
    fmt.Println("Peers:")
    if result.Peers != "" {
      peers := strings.Split(result.Peers, ",")
      for _, peer := range peers {
        // print the peer only if it exist and also active
        // GetPeerInfo return err if the peer doesn't exist, and result.Active = false if it exist but is in-active
        result, err := mnm.middlewareClient.GetPeerInfo(networkAddress, peer)
        if err == nil && result != nil && result.Active {
          fmt.Printf(NETWORK_INFO_FORMAT, result.IP, mkey.AddPrefixPubKeyHash(peer))
        }
      }
    }
  }
}

/*
  Check if middleware is currently running, message if it is not
*/
func (mnm *MarconiNetMode) checkMiddlewareRunning() bool {
  statuses := processes.Instance().GetProcessRunningMap()
  if !statuses[processes.MIDDLEWARE_ID] {
    fmt.Println("Middleware is currently not running")
    return false
  }
  return true
}

/*
  Check if the network contract address has been set, message if it has not been set
*/
func (mnm *MarconiNetMode) checkNetworkSet() bool {
  networkSet := mnm.contractAddress != ""
  if !networkSet {
    fmt.Println("Network has not been set, please set with 'use' command")
  }
  return networkSet
}

/*
  Parse and check whether the waitForReceipt argument is true/false
*/
func parseWaitForReceipt(arg string) bool {
  waitForReceipt := false
  var err error
  waitForReceipt, err = strconv.ParseBool(arg)
  if err != nil {
    return false
  }
  return waitForReceipt
}

/*
  Helper function for printing out response containing the transaction hash
*/
func printTransactionHashResponse(result interface{}, err error) {
  if err != nil {
    fmt.Println("Failed to Submit Transaction. Error:", err)
  } else {
    fmt.Println("Hash for submitted transaction:")
    fmt.Printf("%s\n\n", result.(middleware.TransactionHashResult).TransactionHash)
  }
}
