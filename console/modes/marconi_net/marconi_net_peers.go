package marconi_net

import (
  "github.com/MarconiProtocol/cli/console/modes/marconi_net/commands"
  "github.com/MarconiProtocol/cli/console/util"
  "github.com/MarconiProtocol/go-prompt"
)

// Mode suggestions
var MNET_PEER_SUGGESTIONS = []prompt.Suggest{
  {Text: marconi_net_commands.ADD_PEER, Description: "Add peer to a network"},
  {Text: marconi_net_commands.REMOVE_PEER, Description: "Remove peer from a network"},
  {Text: marconi_net_commands.ADD_PEER_RELATION, Description: "Add peer relation"},
  {Text: marconi_net_commands.REMOVE_PEER_RELATION, Description: "Remove peer relationship"},
  {Text: marconi_net_commands.GET_PEER_RELATIONS, Description: "Get node relationships"},
  {Text: marconi_net_commands.GET_PEER_INFO, Description: "Get node info"},
}

/*
  Show prompt suggestions for add peer command
*/
func (mnm *MarconiNetMode) getAddPeerSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<PEER_NODE_ID>", Description: "NodeID of the peer to add to the network"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Show prompt suggestions for remove peer command
*/
func (mnm *MarconiNetMode) getRemovePeerSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<PEER_NODE_ID>", Description: "NodeID of the peer to remove from the network"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Show prompt suggestions for add peer relation command
*/
func (mnm *MarconiNetMode) getAddPeerRelationSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<PEER_NODE_ID>", Description: "NodeID of the first peer in the relation"}}
  case len(line) == 3:
    return []prompt.Suggest{{Text: "<OTHER_PEER_NODE_ID>", Description: "NodeID of the second peer in the relation"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Show prompt suggestions for add peer relation command
*/
func (mnm *MarconiNetMode) getRemovePeerRelationSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<PEER_NODE_ID>", Description: "NodeID of the first peer in the relation"}}
  case len(line) == 3:
    return []prompt.Suggest{{Text: "<OTHER_PEER_NODE_ID>", Description: "NodeID of the second peer in the relation"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Show prompt suggestions for get peer relations command
*/
func (mnm *MarconiNetMode) getGetPeerRelationsSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<PEER_NODE_ID>", Description: "NodeID of the peer for which to inspect relations"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Show prompt suggestions for get peer info command
*/
func (mnm *MarconiNetMode) getGetPeerInfoSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<PEER_NODE_ID>", Description: "NodeID of the peer for which to inspect info"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Handle the add peer command
*/
func (mnm *MarconiNetMode) handleAddPeer(args []string) {
  util.Logger.Info(marconi_net_commands.PEER+" "+marconi_net_commands.ADD_PEER, util.ArgsToString(args))
  marconi_net_commands.AddPeer(args)
}

/*
  Handle the remove peer command
*/
func (mnm *MarconiNetMode) handleRemovePeer(args []string) {
  util.Logger.Info(marconi_net_commands.PEER+" "+marconi_net_commands.REMOVE_PEER, util.ArgsToString(args))
  marconi_net_commands.RemovePeer(args)
}

/*
  Handle the add peer relation command
*/
func (mnm *MarconiNetMode) handleAddPeerRelation(args []string) {
  util.Logger.Info(marconi_net_commands.PEER+" "+marconi_net_commands.ADD_PEER_RELATION, util.ArgsToString(args))
  marconi_net_commands.AddPeerRelation(args)
}

/*
  Handle the remove peer relation command
*/
func (mnm *MarconiNetMode) handleRemovePeerRelation(args []string) {
  util.Logger.Info(marconi_net_commands.PEER+" "+marconi_net_commands.REMOVE_PEER_RELATION, util.ArgsToString(args))
  marconi_net_commands.RemovePeerRelation(args)
}

/*
  Handle the get peer relations command
*/
func (mnm *MarconiNetMode) handleGetPeerRelations(args []string) {
  util.Logger.Info(marconi_net_commands.PEER+" "+marconi_net_commands.GET_PEER_RELATIONS, util.ArgsToString(args))
  marconi_net_commands.GetPeerRelations(args)
}

/*
  Handle the get peer info command
*/
func (mnm *MarconiNetMode) handleGetPeerInfo(args []string) {
  util.Logger.Info(marconi_net_commands.PEER+" "+marconi_net_commands.GET_PEER_INFO, util.ArgsToString(args))
  marconi_net_commands.GetPeerInfo(args)
}
