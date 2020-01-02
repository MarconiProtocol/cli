package marconid

import (
  "github.com/MarconiProtocol/cli/api/rpc"
)

type Client struct {
  host string
  port string
}

/*
  Let Marconid update network_contract_address in user_config.yml
*/
func UpdateNetworkContractAddress(address string) error {
  args := UpdateNetworkContractAddressArgs{
    NetworkContractAddress: address,
  }
  reply := UpdateNetworkContractAddressReply{}
  return rpc.GetRPCClient().Call("UserConfigService.UpdateNetworkContractAddressRPC", &args, &reply)
}

func SetNetem(interfaceName string, delay uint32, loss float32, duplicate float32, reorderProb float32, corruptProb float32) error {
  args := NetemArgs{
    InterfaceName: interfaceName,
    Delay:         delay,
    Loss:          loss,
    Duplicate:     duplicate,
    ReorderProb:   reorderProb,
    CorruptProb:   corruptProb,
  }
  reply := NetemReply{}
  return rpc.GetRPCClient().Call("TrafficControlService.SetNetemRPC", &args, &reply)
}

func SetTbf(interfaceName string, bandwidth uint64, latencyInMillis float64) error {
  args := TbfArgs{
    InterfaceName:   interfaceName,
    Bandwidth:       bandwidth,
    LatencyInMillis: latencyInMillis,
  }
  reply := TbfReply{}
  return rpc.GetRPCClient().Call("TrafficControlService.SetTbfRPC", &args, &reply)
}

func Reset(interfaceName string) error {
  args := ResetArgs{
    InterfaceName: interfaceName,
  }
  reply := ResetReply{}
  return rpc.GetRPCClient().Call("TrafficControlService.ResetRPC", &args, &reply)
}
