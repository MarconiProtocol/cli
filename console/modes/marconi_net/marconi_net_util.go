package marconi_net

import (
  "github.com/MarconiProtocol/cli/console/modes/marconi_net/commands"
  "github.com/MarconiProtocol/cli/console/util"
  "github.com/MarconiProtocol/go-prompt"
)

// Mode suggestions
var MNET_UTIL_SUGGESTIONS = []prompt.Suggest{
  {Text: marconi_net_commands.GENERATE_32BITKEY, Description: "Generate a 32 bit key"},
  {Text: marconi_net_commands.GET_MPIPE_PORT, Description: "Get the shared MPipe port"},
  {Text: marconi_net_commands.REGISTER, Description: "Register a nodeID"},
  {Text: marconi_net_commands.START_NETFLOW, Description: "Start netflow monitoring"},
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

func (mnm *MarconiNetMode) getGetMpipePortSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<PEER_NODE_ID>", Description: "NodeID of the first peer"}}
  case len(line) == 3:
    return []prompt.Suggest{{Text: "<OTHER_PEER_NODE_ID>", Description: "NodeID of the second peer"}}
  default:
    return []prompt.Suggest{}
  }
}

func (mnm *MarconiNetMode) getStartNetflowSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<COLLECTOR_IP>", Description: "IP of the collector to send netflow data to."}}
  case len(line) == 3:
    return []prompt.Suggest{{Text: "<COLLECTOR_PORT>", Description: "Port of the collector to send netflow data to."}}
  case len(line) == 4:
    return []prompt.Suggest{{Text: "<INTERFACE_ID>", Description: "Interface to monitor."}}
  case len(line) == 5:
    return []prompt.Suggest{{Text: "[LOGGING_DIR]", Description: "Optional directory to output local netflow data to."}}
  default:
    return []prompt.Suggest{}
  }
}

func (mnm *MarconiNetMode) handleGenerate32BitKey(args []string) {
  util.Logger.Info(marconi_net_commands.UTIL+" "+marconi_net_commands.GENERATE_32BITKEY, util.ArgsToString(args))
  marconi_net_commands.Generate32BitKey(args)
}

/*
  Handle the register command
*/
func (mnm *MarconiNetMode) handleRegisterUser(args []string) {
  util.Logger.Info(marconi_net_commands.UTIL+" "+marconi_net_commands.REGISTER, util.ArgsToString(args))
  marconi_net_commands.Register(args)
}

/*
  Handle the get mpipe port command
*/
func (mnm *MarconiNetMode) handleGetMPipePort(args []string) {
  util.Logger.Info(marconi_net_commands.UTIL+" "+marconi_net_commands.GET_MPIPE_PORT, util.ArgsToString(args))
  marconi_net_commands.GetMPipePort(args)
}

/*
  Handle the start netflow command
*/
func (mnm *MarconiNetMode) handleStartNetflow(args []string) {
  util.Logger.Info(marconi_net_commands.UTIL+" "+marconi_net_commands.START_NETFLOW, util.ArgsToString(args))
  marconi_net_commands.StartNetflow(args)
}
