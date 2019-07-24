package marconi_net

import (
  "fmt"
  "github.com/MarconiProtocol/go-prompt"
  "github.com/MarconiProtocol/cli/api/middleware"
  "github.com/MarconiProtocol/cli/console/modes"
  "github.com/MarconiProtocol/cli/core/mkey"
  "strings"
)

const (
  ADD_PEER             = "add"
  REMOVE_PEER          = "remove"
  ADD_PEER_RELATION    = "add_relation"
  REMOVE_PEER_RELATION = "remove_relation"
  GET_PEER_RELATIONS   = "relations"
  GET_PEER_INFO        = "info"
)

// Mode suggestions
var MNET_PEER_SUGGESTIONS = []prompt.Suggest{
  {Text: ADD_PEER, Description: "Add peer to a network"},
  {Text: REMOVE_PEER, Description: "Remove peer from a network"},
  {Text: ADD_PEER_RELATION, Description: "Add peer relation"},
  {Text: REMOVE_PEER_RELATION, Description: "Remove peer relationship"},
  {Text: GET_PEER_RELATIONS, Description: "Get node relationships"},
  {Text: GET_PEER_INFO, Description: "Get node info"},
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

func (mnm *MarconiNetMode) checkPeerPrerequisites() bool {
  return mnm.checkMiddlewareRunning() && mnm.checkNetworkSet()
}

/*
  Handle the add peer command
*/
func (mnm *MarconiNetMode) handleAddPeer(args []string) {
  if !mnm.checkPeerPrerequisites() {
    return
  }
  if !modes.ArgsLenCheckWithOptional(args, 1, 1) {
    fmt.Println("Usage:", ADD_PEER, "<PEER_NODE_ID>")
    return
  }
  if !modes.ArgPubKeyHashCheck(args[0]) {
    return
  }

  // check if the network does not already contain the peer
  if mnm.checkNetworkContains(mnm.contractAddress, mkey.StripPrefixPubKeyHash(args[0])) {
    fmt.Println("Network", mnm.contractAddress, "already contain peer", args[0])
    return
  }

  if modes.ArgsLenCheck(args, 2) && parseWaitForReceipt(args[1]) {
    result, err := mnm.middlewareClient.AddPeer(mnm.contractAddress, mkey.StripPrefixPubKeyHash(args[0]), true)
    r := result.(middleware.AddPeerResult)
    if err != nil {
      fmt.Println("Error:", err)
    } else {
      fmt.Println("Added a peer to the network:")
      fmt.Printf("%-24s %48s\n", "Network Id", r.NetworkId)
      fmt.Printf("%-24s %48s\n\n", "Added Peer", mkey.AddPrefixPubKeyHash(r.PubKeyHash))
    }
  } else {
    result, err := mnm.middlewareClient.AddPeer(mnm.contractAddress, mkey.StripPrefixPubKeyHash(args[0]), false)
    printTransactionHashResponse(result, err)
  }
}

/*
  Handle the remove peer command
*/
func (mnm *MarconiNetMode) handleRemovePeer(args []string) {
  if !mnm.checkPeerPrerequisites() {
    return
  }
  if !modes.ArgsLenCheckWithOptional(args, 1, 1) {
    fmt.Println("Usage:", REMOVE_PEER, "<PEER_NODE_ID>")
    return
  }
  if !modes.ArgPubKeyHashCheck(args[0]) {
    return
  }

  // check if the network contains the peer
  if !mnm.checkNetworkContains(mnm.contractAddress, mkey.StripPrefixPubKeyHash(args[0])) {
    fmt.Println("Network", mnm.contractAddress, "does not contain peer", args[0])
    return
  }

  if modes.ArgsLenCheck(args, 2) && parseWaitForReceipt(args[1]) {
    result, err := mnm.middlewareClient.RemovePeer(mnm.contractAddress, mkey.StripPrefixPubKeyHash(args[0]), true)
    r := result.(middleware.RemovePeerResult)
    if err != nil {
      fmt.Println("Error:", err)
    } else {
      fmt.Println("Removed a peer from the network:")
      fmt.Printf("%-24s %48s\n", "Network Id", r.NetworkId)
      fmt.Printf("%-24s %48s\n\n", "Removed Peer", mkey.AddPrefixPubKeyHash(r.PubKeyHash))
    }
  } else {
    result, err := mnm.middlewareClient.RemovePeer(mnm.contractAddress, mkey.StripPrefixPubKeyHash(args[0]), false)
    printTransactionHashResponse(result, err)
  }
}

/*
  Handle the add peer relation command
*/
func (mnm *MarconiNetMode) handleAddPeerRelation(args []string) {
  if !mnm.checkPeerPrerequisites() {
    return
  }
  if !modes.ArgsLenCheckWithOptional(args, 2, 1) {
    fmt.Println("Usage:", ADD_PEER_RELATION, "<PEER_NODE_ID>", "<OTHER_PEER_NODE_ID>")
    return
  }
  if !modes.ArgPubKeyHashCheck(args[0]) ||
     !modes.ArgPubKeyHashCheck(args[1]) {
    return
  }
  if args[0] == args[1] {
    fmt.Println("Adding a self-relation does not make sense")
    return
  }

  if modes.ArgsLenCheck(args, 3) && parseWaitForReceipt(args[2]) {
    result, err := mnm.middlewareClient.AddPeerRelation(mnm.contractAddress, mkey.StripPrefixPubKeyHash(args[0]), mkey.StripPrefixPubKeyHash(args[1]), true)
    r := result.(middleware.AddPeerRelationResult)
    if err != nil {
      fmt.Println("Error:", err)
    } else {
      fmt.Println("Added a new peer relationship:")
      fmt.Printf("%-24s %48s\n", "Network Id", r.NetworkId)
      fmt.Printf("%-24s %48s\n", "Peer", mkey.AddPrefixPubKeyHash(r.PubKeyHashMine))
      fmt.Printf("%-24s %48s\n\n", "Other Peer", mkey.AddPrefixPubKeyHash(r.PubKeyHashOther))
    }
  } else {
    result, err := mnm.middlewareClient.AddPeerRelation(mnm.contractAddress, mkey.StripPrefixPubKeyHash(args[0]), mkey.StripPrefixPubKeyHash(args[1]), false)
    printTransactionHashResponse(result, err)
  }
}

/*
  Handle the remove peer relation command
*/
func (mnm *MarconiNetMode) handleRemovePeerRelation(args []string) {
  if !mnm.checkPeerPrerequisites() {
    return
  }
  if !modes.ArgsLenCheckWithOptional(args, 2, 1) {
    fmt.Println("Usage:", REMOVE_PEER_RELATION, "<PEER_NODE_ID>", "<OTHER_PEER_NODE_ID>")
    return
  }
  if !modes.ArgPubKeyHashCheck(args[0]) ||
     !modes.ArgPubKeyHashCheck(args[1]) {
    return
  }

  if modes.ArgsLenCheck(args, 3) && parseWaitForReceipt(args[2]) {
    result, err := mnm.middlewareClient.RemovePeerRelation(mnm.contractAddress, mkey.StripPrefixPubKeyHash(args[0]), mkey.StripPrefixPubKeyHash(args[1]), true)
    r := result.(middleware.RemovePeerRelationResult)
    if err != nil {
      fmt.Println("Error:", err)
    } else {
      fmt.Println("Removed a peer relationship:")
      fmt.Printf("%-24s %48s\n", "Network Id", r.NetworkId)
      fmt.Printf("%-24s %48s\n", "Peer", mkey.AddPrefixPubKeyHash(r.PubKeyHashMine))
      fmt.Printf("%-24s %48s\n\n", "Other Peer", mkey.AddPrefixPubKeyHash(r.PubKeyHashOther))
    }
  } else {
    result, err := mnm.middlewareClient.RemovePeerRelation(mnm.contractAddress, mkey.StripPrefixPubKeyHash(args[0]), mkey.StripPrefixPubKeyHash(args[1]), false)
    printTransactionHashResponse(result, err)
  }
}

/*
  Handle the get peer relations command
*/
func (mnm *MarconiNetMode) handleGetPeerRelations(args []string) {
  if !mnm.checkPeerPrerequisites() {
    return
  }
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", GET_PEER_RELATIONS, "<NODE_ID>")
    return
  }
  if !modes.ArgPubKeyHashCheck(args[0]) {
    return
  }

  result, err := mnm.middlewareClient.GetPeerRelations(mnm.contractAddress, mkey.StripPrefixPubKeyHash(args[0]))
  if err != nil {
    fmt.Println("Error:", err)
  } else {
    fmt.Println("Retrieved peers:")
    peers := strings.Split(result, ",")
    for num, peer := range peers {
      fmt.Printf("Peer%-24d %48s\n", num, mkey.AddPrefixPubKeyHash(peer))
    }
  }
}

/*
  Handle the get peer info command
*/
func (mnm *MarconiNetMode) handleGetPeerInfo(args []string) {
  if !mnm.checkPeerPrerequisites() {
    return
  }
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", GET_PEER_INFO, "<NODE_ID>")
    return
  }
  if !modes.ArgPubKeyHashCheck(args[0]) {
    return
  }

  result, err := mnm.middlewareClient.GetPeerInfo(mnm.contractAddress, mkey.StripPrefixPubKeyHash(args[0]))
  if err != nil {
    fmt.Println("Error:", err)
  } else {
    fmt.Println("Peer Info:")
    fmt.Printf("%-24s %48s\n", "Network Id", result.NetworkId)
    //fmt.Printf("%-24s %48s\n", "NodeID", mkey.AddPrefixPubKeyHash(result.PubKeyHash))
    fmt.Printf("%-24s %48s\n", "IP", result.IP)
    fmt.Printf("%-24s %48v\n", "Active", result.Active)

    fmt.Println("Peers:")
    peers := strings.Split(result.Peers, ",")
    for num, peer := range peers {
      fmt.Printf("%-24d %48s\n", num, mkey.AddPrefixPubKeyHash(peer))
    }

  }
}

func (mnm *MarconiNetMode) checkNetworkContains(network string, peer string) bool {
  result, err := mnm.middlewareClient.GetPeerInfo(network, peer)
  if err == nil && result != nil {
    return result.Active
  }
  return false
}
