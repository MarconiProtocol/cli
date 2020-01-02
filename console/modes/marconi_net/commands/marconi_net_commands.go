package marconi_net_commands

import (
  "fmt"
  "github.com/MarconiProtocol/cli/api/marconid"
  "github.com/MarconiProtocol/cli/api/middleware"
  "github.com/MarconiProtocol/cli/console/modes"
  "github.com/MarconiProtocol/cli/console/util"
  "github.com/MarconiProtocol/cli/core/mkey"
  "github.com/MarconiProtocol/cli/core/processes"
  "strconv"
  "strings"
)

// Commands
const (
  PEER             = "peer"
  UTIL             = "util"
  USE              = "use"
  CREATE_NETWORK   = "create"
  DELETE_NETWORK   = "delete"
  JOIN_NETWORK     = "join"
  GET_NETWORK_INFO = "info"
  TRAFFIC_CONTROL  = "tc"
)
const (
  NETWORK_INFO_FORMAT = "%-24s %48s\n"
)

var ContractAddress = ""

// Mapping Commands to Functions
var COMMAND_MAP = map[string]func([]string){
  PEER:             HandlePeerCommand,
  UTIL:             HandleUtilCommand,
  USE:              UseNetwork,
  CREATE_NETWORK:   CreateNetwork,
  DELETE_NETWORK:   DeleteNetwork,
  JOIN_NETWORK:     JoinNetwork,
  GET_NETWORK_INFO: GetNetworkInfo,
  TRAFFIC_CONTROL:  HandleTCCommand,
}

func checkMiddlewareRunning() bool {
  statuses := processes.Instance().GetProcessRunningMap()
  if !statuses[processes.MIDDLEWARE_ID] {
    fmt.Println("Middleware is currently not running")
    return false
  }
  return true

}

func printNetworkInfo(action, networkId, admin, networkContract string) {
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
  Check if the network contract address has been set, message if it has not been set
*/
func checkNetworkSet() bool {
  networkSet := ContractAddress != ""
  if !networkSet {
    fmt.Println("Network has not been set, please set with 'use' command")
  }
  return networkSet
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

func UseNetwork(args []string) {
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", USE, "<NETWORK_CONTRACT_ADDRESS>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }
  ContractAddress = util.GetEIP55Address(args[0])
  fmt.Println("Using network contract")
  fmt.Println(ContractAddress)
}

func CreateNetwork(args []string) {
  if !checkMiddlewareRunning() {
    return
  }
  if !modes.ArgsLenCheckWithOptional(args, 0, 1) {
    fmt.Println("Usage:", CREATE_NETWORK)
    return
  }

  fmt.Println("This may take up to a minute...")
  result, err := middleware.GetClient().CreateNetwork()
  if err != nil {
    fmt.Println("Failed to create network:", err)
  } else {
    printNetworkInfo("Created a new network", result.NetworkId, result.Admin, result.NetworkContract)

    // use the network is the optional argument is provided
    if len(args) == 1 && args[0] == "use" {
      UseNetwork([]string{result.NetworkContract})
    }
  }
}

func DeleteNetwork(args []string) {
  if !checkMiddlewareRunning() {
    return
  }
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", DELETE_NETWORK, "<0xNETWORK_CONTRACT_ADDRESS>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }
  result, err := middleware.GetClient().DeleteNetwork(args[0])
  if err != nil {
    fmt.Println("Error:", err)
  } else {
    printNetworkInfo("Deleted network", result.NetworkId, result.Admin, "")
  }
}

func JoinNetwork(args []string) {
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", JOIN_NETWORK, "<0xNETWORK_CONTRACT_ADDRESS>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }

  if err := marconid.UpdateNetworkContractAddress(args[0]); err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Joined network", args[0])
  }
}

func GetNetworkInfo(args []string) {
  if !checkMiddlewareRunning() {
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
  result, err := middleware.GetClient().GetNetworkInfo(networkAddress)
  if err != nil {
    fmt.Println("Error:", err)
  } else {
    printNetworkInfo("Network Info", result.NetworkId, result.NetworkAdmin, networkAddress)
    fmt.Println("Peers:")
    if result.Peers != "" {
      peers := strings.Split(result.Peers, ",")
      for _, peer := range peers {
        // print the peer only if it exist and also active
        // GetPeerInfo return err if the peer doesn't exist, and result.Active = false if it exist but is in-active
        result, err := middleware.GetClient().GetPeerInfo(networkAddress, peer)
        if err == nil && result != nil && result.Active {
          fmt.Printf(NETWORK_INFO_FORMAT, result.IP, mkey.AddPrefixPubKeyHash(peer))
        }
      }
    }
  }
}
