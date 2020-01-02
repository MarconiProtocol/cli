package marconi_net

import (
  "fmt"
  "github.com/MarconiProtocol/cli/api/middleware"
  "github.com/MarconiProtocol/cli/console/context"
  "github.com/MarconiProtocol/cli/console/modes"
  "github.com/MarconiProtocol/cli/console/modes/marconi_net/commands"
  "github.com/MarconiProtocol/cli/console/util"
  "github.com/MarconiProtocol/go-prompt"
)

// Name of Mode
const NAME = "net"

const NETWORK_INFO_FORMAT = "%-24s %48s\n"

// Mode suggestions
var MNET_SUGGESTIONS = []prompt.Suggest{
  {Text: marconi_net_commands.PEER, Description: "Peer related commands"},
  {Text: marconi_net_commands.UTIL, Description: "Utility commands"},
  {Text: marconi_net_commands.TRAFFIC_CONTROL, Description: "Traffic control commands"},
  {Text: marconi_net_commands.USE, Description: "Set network to use with other commands"},
  {Text: marconi_net_commands.CREATE_NETWORK, Description: "Create new network"},
  {Text: marconi_net_commands.DELETE_NETWORK, Description: "Delete existing network"},
  {Text: marconi_net_commands.JOIN_NETWORK, Description: "Join an existing network"},
  {Text: marconi_net_commands.GET_NETWORK_INFO, Description: "Get network info"},
  {Text: modes.RETURN_TO_ROOT, Description: "Return to home menu"},
  {Text: modes.EXIT_CMD, Description: "Exit mcli"},
}

/*
  Menu/Mode for Marconi Net related operations
*/
type MarconiNetMode struct {
  modes.BaseMode
  middlewareClient *middleware.Client
  contractAddress  *string
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
  mnetMode.RegisterCommand(marconi_net_commands.USE, mnetMode.getUseNetworkSuggestions, mnetMode.handleUseNetwork)
  mnetMode.RegisterCommand(marconi_net_commands.CREATE_NETWORK, mnetMode.getCreateNetworkSuggestions, mnetMode.handleCreateNetwork)
  mnetMode.RegisterCommand(marconi_net_commands.DELETE_NETWORK, mnetMode.getDeleteNetworkSuggestions, mnetMode.handleDeleteNetwork)
  mnetMode.RegisterCommand(marconi_net_commands.JOIN_NETWORK, mnetMode.getJoinNetworkSuggestions, mnetMode.handleJoinNetwork)
  mnetMode.RegisterCommand(marconi_net_commands.GET_NETWORK_INFO, mnetMode.getGetNetworkInfoSuggestions, mnetMode.handleGetNetworkInfo)

  mnetMode.RegisterCommand(marconi_net_commands.PEER, mnetMode.getPeerSuggestions, mnetMode.handlePeer)
  mnetMode.RegisterSubCommand(marconi_net_commands.PEER, marconi_net_commands.ADD_PEER, mnetMode.getAddPeerSuggestions, mnetMode.handleAddPeer)
  mnetMode.RegisterSubCommand(marconi_net_commands.PEER, marconi_net_commands.REMOVE_PEER, mnetMode.getRemovePeerSuggestions, mnetMode.handleRemovePeer)
  mnetMode.RegisterSubCommand(marconi_net_commands.PEER, marconi_net_commands.ADD_PEER_RELATION, mnetMode.getAddPeerRelationSuggestions, mnetMode.handleAddPeerRelation)
  mnetMode.RegisterSubCommand(marconi_net_commands.PEER, marconi_net_commands.REMOVE_PEER_RELATION, mnetMode.getRemovePeerRelationSuggestions, mnetMode.handleRemovePeerRelation)
  mnetMode.RegisterSubCommand(marconi_net_commands.PEER, marconi_net_commands.GET_PEER_RELATIONS, mnetMode.getGetPeerRelationsSuggestions, mnetMode.handleGetPeerRelations)
  mnetMode.RegisterSubCommand(marconi_net_commands.PEER, marconi_net_commands.GET_PEER_INFO, mnetMode.getGetPeerInfoSuggestions, mnetMode.handleGetPeerInfo)

  mnetMode.RegisterCommand(marconi_net_commands.UTIL, mnetMode.getUtilSuggestions, mnetMode.handleUtil)
  mnetMode.RegisterSubCommand(marconi_net_commands.UTIL, marconi_net_commands.GENERATE_32BITKEY, mnetMode.getGenerate32BitKeySuggestions, mnetMode.handleGenerate32BitKey)
  mnetMode.RegisterSubCommand(marconi_net_commands.UTIL, marconi_net_commands.REGISTER, mnetMode.getRegisterUserSuggestions, mnetMode.handleRegisterUser)
  mnetMode.RegisterSubCommand(marconi_net_commands.UTIL, marconi_net_commands.GET_MPIPE_PORT, mnetMode.getGetMpipePortSuggestions, mnetMode.handleGetMPipePort)
  mnetMode.RegisterSubCommand(marconi_net_commands.UTIL, marconi_net_commands.START_NETFLOW, mnetMode.getStartNetflowSuggestions, mnetMode.handleStartNetflow)

  mnetMode.RegisterCommand(marconi_net_commands.TRAFFIC_CONTROL, mnetMode.getTCSuggestions, mnetMode.handleTC)
  mnetMode.RegisterSubCommand(marconi_net_commands.TRAFFIC_CONTROL, marconi_net_commands.INFO, mnetMode.getInfoSuggestions, mnetMode.handleInfo)
  mnetMode.RegisterSubCommand(marconi_net_commands.TRAFFIC_CONTROL, marconi_net_commands.SET_BANDWIDTH, mnetMode.getSetBandwidthSuggestions, mnetMode.handleSetBandwidth)
  mnetMode.RegisterSubCommand(marconi_net_commands.TRAFFIC_CONTROL, marconi_net_commands.SET_CORRUPT, mnetMode.getSetCorruptSuggestions, mnetMode.handleSetCorrupt)
  mnetMode.RegisterSubCommand(marconi_net_commands.TRAFFIC_CONTROL, marconi_net_commands.SET_DELAY, mnetMode.getSetDelaySuggestions, mnetMode.handleSetDelay)
  mnetMode.RegisterSubCommand(marconi_net_commands.TRAFFIC_CONTROL, marconi_net_commands.SET_DUPLICATE, mnetMode.getSetDuplicateSuggestions, mnetMode.handleSetDuplicate)
  mnetMode.RegisterSubCommand(marconi_net_commands.TRAFFIC_CONTROL, marconi_net_commands.SET_LOSS, mnetMode.getSetLossSuggestions, mnetMode.handleSetLoss)
  mnetMode.RegisterSubCommand(marconi_net_commands.TRAFFIC_CONTROL, marconi_net_commands.SET_REORDER, mnetMode.getSetReorderSuggestions, mnetMode.handleSetReorder)
  mnetMode.RegisterSubCommand(marconi_net_commands.TRAFFIC_CONTROL, marconi_net_commands.RESET, mnetMode.getResetSuggestions, mnetMode.handleReset)

  mnetMode.RegisterCommand(modes.RETURN_TO_ROOT, mnetMode.GetEmptySuggestions, mnetMode.HandleReturnToRoot)
  mnetMode.RegisterCommand(modes.EXIT_CMD, mnetMode.GetEmptySuggestions, mnetMode.HandleExitCommand)

  mnetMode.middlewareClient = middleware.GetClient()
  mnetMode.contractAddress = &marconi_net_commands.ContractAddress
  return &mnetMode
}

func (mnm *MarconiNetMode) CliPrefix() (string, bool) {
  var networkPrefix string
  if *mnm.contractAddress != "" {
    contractAddress := *mnm.contractAddress
    networkPrefix = " <..." + contractAddress[len(contractAddress)-8:] + ">"
  }
  return mnm.Name() + networkPrefix, true
}

func (mnm *MarconiNetMode) Name() string {
  return "net"
}

func (mnm *MarconiNetMode) HandleCommand(args []string) {
  if !modes.ArgsMinLenCheck(args, 1) {
    fmt.Println("USAGE: <command>")
    return
  }
  commandType := args[0]
  commandArgs := args[1:]
  if commandHandlerFunction, present := marconi_net_commands.COMMAND_MAP[commandType]; present {
    commandHandlerFunction(commandArgs)
  } else {
    fmt.Println("Invalid command " + commandType)
  }
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
  Show prompt suggestions for entering the tc sub menu
*/
func (mnm *MarconiNetMode) getTCSuggestions(line []string) []prompt.Suggest {
  return util.SimpleSubcommandCompleter(line, 1, MNET_TC_SUGGESTIONS)
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
  util.HandleFurtherCommands(marconi_net_commands.PEER, MNET_PEER_SUGGESTIONS)
}

/*
  Handle the util sub menu command
*/
func (mnm *MarconiNetMode) handleUtil(args []string) {
  util.HandleFurtherCommands(marconi_net_commands.UTIL, MNET_UTIL_SUGGESTIONS)
}

/*
  Handle the tc sub menu command
*/
func (mnm *MarconiNetMode) handleTC(args []string) {
  util.HandleFurtherCommands(marconi_net_commands.TRAFFIC_CONTROL, MNET_TC_SUGGESTIONS)
}

/*
  Handle the use network command
*/
func (mnm *MarconiNetMode) handleUseNetwork(args []string) {
  util.Logger.Info(marconi_net_commands.USE, util.ArgsToString(args))
  marconi_net_commands.UseNetwork(args)
}

/*
  Handle the create network command
*/
func (mnm *MarconiNetMode) handleCreateNetwork(args []string) {
  util.Logger.Info(marconi_net_commands.CREATE_NETWORK, util.ArgsToString(args))
  marconi_net_commands.CreateNetwork(args)
}

/*
  Handle the delete network command
*/
func (mnm *MarconiNetMode) handleDeleteNetwork(args []string) {
  util.Logger.Info(marconi_net_commands.DELETE_NETWORK, util.ArgsToString(args))
  marconi_net_commands.DeleteNetwork(args)
}

/*
  Handle the join network command
*/
func (mnm *MarconiNetMode) handleJoinNetwork(args []string) {
  util.Logger.Info(marconi_net_commands.JOIN_NETWORK, util.ArgsToString(args))
  marconi_net_commands.JoinNetwork(args)
}

/*
  Handle the get network info command
*/
func (mnm *MarconiNetMode) handleGetNetworkInfo(args []string) {
  util.Logger.Info(marconi_net_commands.GET_NETWORK_INFO, util.ArgsToString(args))
  marconi_net_commands.GetNetworkInfo(args)
}
