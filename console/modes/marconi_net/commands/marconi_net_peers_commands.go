package marconi_net_commands

import (
  "fmt"
  "github.com/MarconiProtocol/cli/api/middleware"
  "github.com/MarconiProtocol/cli/console/modes"
  "github.com/MarconiProtocol/cli/core/mkey"
  "strings"
)

// Commands
const (
  ADD_PEER             = "add"
  REMOVE_PEER          = "remove"
  ADD_PEER_RELATION    = "add_relation"
  REMOVE_PEER_RELATION = "remove_relation"
  GET_PEER_RELATIONS   = "relations"
  GET_PEER_INFO        = "info"
)

var PEER_COMMAND_MAP = map[string]func([]string){
  ADD_PEER:             AddPeer,
  REMOVE_PEER:          RemovePeer,
  ADD_PEER_RELATION:    AddPeerRelation,
  REMOVE_PEER_RELATION: RemovePeerRelation,
  GET_PEER_RELATIONS:   GetPeerRelations,
  GET_PEER_INFO:        GetPeerInfo,
}

func HandlePeerCommand(args []string) {
  if !modes.ArgsMinLenCheck(args, 1) {
    fmt.Println("USAGE: <command>")
    return
  }
  commandType := args[0]
  commandArgs := args[1:]
  if commandHandlerFunction, present := PEER_COMMAND_MAP[commandType]; present {
    commandHandlerFunction(commandArgs)
  } else {
    fmt.Println("Invalid command " + commandType)
  }
}

func checkPeerPrerequisites() bool {
  return checkMiddlewareRunning() && checkNetworkSet()
}

func checkNetworkContains(network string, peer string) bool {
  result, err := middleware.GetClient().GetPeerInfo(network, peer)
  if err == nil && result != nil {
    return result.Active
  }
  return false
}

func AddPeer(args []string) {
  if !checkPeerPrerequisites() {
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
  if checkNetworkContains(ContractAddress, mkey.StripPrefixPubKeyHash(args[0])) {
    fmt.Println("Network", ContractAddress, "already contain peer", args[0])
    return
  }

  if modes.ArgsLenCheck(args, 2) && parseWaitForReceipt(args[1]) {
    result, err := middleware.GetClient().AddPeer(ContractAddress, mkey.StripPrefixPubKeyHash(args[0]), true)
    r := result.(middleware.AddPeerResult)
    if err != nil {
      fmt.Println("Error:", err)
    } else {
      fmt.Println("Added a peer to the network:")
      fmt.Printf("%-24s %48s\n", "Network Id", r.NetworkId)
      fmt.Printf("%-24s %48s\n\n", "Added Peer", mkey.AddPrefixPubKeyHash(r.PubKeyHash))
    }
  } else {
    result, err := middleware.GetClient().AddPeer(ContractAddress, mkey.StripPrefixPubKeyHash(args[0]), false)
    printTransactionHashResponse(result, err)
  }
}

func RemovePeer(args []string) {
  if !checkPeerPrerequisites() {
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
  if !checkNetworkContains(ContractAddress, mkey.StripPrefixPubKeyHash(args[0])) {
    fmt.Println("Network", ContractAddress, "does not contain peer", args[0])
    return
  }

  if modes.ArgsLenCheck(args, 2) && parseWaitForReceipt(args[1]) {
    result, err := middleware.GetClient().RemovePeer(ContractAddress, mkey.StripPrefixPubKeyHash(args[0]), true)
    r := result.(middleware.RemovePeerResult)
    if err != nil {
      fmt.Println("Error:", err)
    } else {
      fmt.Println("Removed a peer from the network:")
      fmt.Printf("%-24s %48s\n", "Network Id", r.NetworkId)
      fmt.Printf("%-24s %48s\n\n", "Removed Peer", mkey.AddPrefixPubKeyHash(r.PubKeyHash))
    }
  } else {
    result, err := middleware.GetClient().RemovePeer(ContractAddress, mkey.StripPrefixPubKeyHash(args[0]), false)
    printTransactionHashResponse(result, err)
  }
}

func AddPeerRelation(args []string) {
  if !checkPeerPrerequisites() {
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
    result, err := middleware.GetClient().AddPeerRelation(ContractAddress, mkey.StripPrefixPubKeyHash(args[0]), mkey.StripPrefixPubKeyHash(args[1]), true)
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
    result, err := middleware.GetClient().AddPeerRelation(ContractAddress, mkey.StripPrefixPubKeyHash(args[0]), mkey.StripPrefixPubKeyHash(args[1]), false)
    printTransactionHashResponse(result, err)
  }
}

func RemovePeerRelation(args []string) {
  if !checkPeerPrerequisites() {
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
    result, err := middleware.GetClient().RemovePeerRelation(ContractAddress, mkey.StripPrefixPubKeyHash(args[0]), mkey.StripPrefixPubKeyHash(args[1]), true)
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
    result, err := middleware.GetClient().RemovePeerRelation(ContractAddress, mkey.StripPrefixPubKeyHash(args[0]), mkey.StripPrefixPubKeyHash(args[1]), false)
    printTransactionHashResponse(result, err)
  }
}

func GetPeerRelations(args []string) {
  if !checkPeerPrerequisites() {
    return
  }
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", GET_PEER_RELATIONS, "<NODE_ID>")
    return
  }
  if !modes.ArgPubKeyHashCheck(args[0]) {
    return
  }

  result, err := middleware.GetClient().GetPeerRelations(ContractAddress, mkey.StripPrefixPubKeyHash(args[0]))
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

func GetPeerInfo(args []string) {

  if !checkPeerPrerequisites() {
    return
  }
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", GET_PEER_INFO, "<NODE_ID>")
    return
  }
  if !modes.ArgPubKeyHashCheck(args[0]) {
    return
  }

  result, err := middleware.GetClient().GetPeerInfo(ContractAddress, mkey.StripPrefixPubKeyHash(args[0]))
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
