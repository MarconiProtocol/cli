package marconi_net

import (
  "github.com/MarconiProtocol/cli/console/modes/marconi_net/commands"
  "github.com/MarconiProtocol/cli/console/util"
  "github.com/MarconiProtocol/go-prompt"
)

// Mode suggestions
var MNET_TC_SUGGESTIONS = []prompt.Suggest{
  {Text: marconi_net_commands.INFO, Description: "Show the qdisc attached to this network interface"},
  {Text: marconi_net_commands.SET_BANDWIDTH, Description: "Set bandwidth control"},
  {Text: marconi_net_commands.SET_CORRUPT, Description: "Set the corrupt rate of packets"},
  {Text: marconi_net_commands.SET_DELAY, Description: "Set the delay of packets"},
  {Text: marconi_net_commands.SET_DUPLICATE, Description: "Set the duplicate rate of packets"},
  {Text: marconi_net_commands.SET_LOSS, Description: "Set the loss rate of packets"},
  {Text: marconi_net_commands.SET_REORDER, Description: "Set the reorder rate of packets"},
  {Text: marconi_net_commands.RESET, Description: "Reset the qdisc for this interface"},
}

/*
 Show prompt suggestions for info command
*/
func (mnm *MarconiNetMode) getInfoSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<NET_INTERFACE>", Description: "The network interface (mpipe or mbridge) to check"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
 Show prompt suggestions for set_bandwidth command
*/
func (mnm *MarconiNetMode) getSetBandwidthSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<NET_INTERFACE>", Description: "Network interface (mpipe or mbridge)"}}
  case len(line) == 3:
    return []prompt.Suggest{{Text: "<BANDWIDTH>", Description: "Bandwidth in bytes/s"}}
  case len(line) == 4:
    return []prompt.Suggest{{Text: "<LATENCY>", Description: "Latency in millisecond"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
 Show prompt suggestions for set_corrupt command
*/
func (mnm *MarconiNetMode) getSetCorruptSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<NET_INTERFACE>", Description: "Network interface (mpipe or mbridge)"}}
  case len(line) == 3:
    return []prompt.Suggest{{Text: "<CORRUPT_RATE>", Description: "Corrupt rate of packets in percentage (e.g., 20)"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
 Show prompt suggestions for set_delay command
*/
func (mnm *MarconiNetMode) getSetDelaySuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<NET_INTERFACE>", Description: "Network interface (mpipe or mbridge)"}}
  case len(line) == 3:
    return []prompt.Suggest{{Text: "<DELAY>", Description: "Delay of packets in milliseconds (e.g., 150)"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
 Show prompt suggestions for set_duplicate command
*/
func (mnm *MarconiNetMode) getSetDuplicateSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<NET_INTERFACE>", Description: "Network interface (mpipe or mbridge)"}}
  case len(line) == 3:
    return []prompt.Suggest{{Text: "<DUPLICATE_RATE>", Description: "Duplicate rate of packets in percentage (e.g., 20)"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
 Show prompt suggestions for set_loss command
*/
func (mnm *MarconiNetMode) getSetLossSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<NET_INTERFACE>", Description: "Network interface (mpipe or mbridge)"}}
  case len(line) == 3:
    return []prompt.Suggest{{Text: "<LOSS_RATE>", Description: "Loss rate of packets in percentage (e.g., 20)"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
 Show prompt suggestions for set_reorder command
*/
func (mnm *MarconiNetMode) getSetReorderSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<NET_INTERFACE>", Description: "Network interface (mpipe or mbridge)"}}
  case len(line) == 3:
    return []prompt.Suggest{{Text: "<REORDER_RATE>", Description: "Reorder rate of packets in percentage (e.g., 20)"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
 Show prompt suggestions for reset command
*/
func (mnm *MarconiNetMode) getResetSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<NET_INTERFACE>", Description: "Network interface (mpipe or mbridge)"}}
  default:
    return []prompt.Suggest{}
  }
}

func (mnm *MarconiNetMode) handleInfo(args []string) {
  util.Logger.Info(marconi_net_commands.TRAFFIC_CONTROL+" "+marconi_net_commands.INFO, util.ArgsToString(args))
  marconi_net_commands.HandleInfo(args)
}

func (mnm *MarconiNetMode) handleSetBandwidth(args []string) {
  util.Logger.Info(marconi_net_commands.TRAFFIC_CONTROL+" "+marconi_net_commands.SET_BANDWIDTH, util.ArgsToString(args))
  marconi_net_commands.HandleSetBandwidth(args)
}

func (mnm *MarconiNetMode) handleSetCorrupt(args []string) {
  util.Logger.Info(marconi_net_commands.TRAFFIC_CONTROL+" "+marconi_net_commands.SET_CORRUPT, util.ArgsToString(args))
  marconi_net_commands.HandleSetCorrupt(args)
}

func (mnm *MarconiNetMode) handleSetDelay(args []string) {
  util.Logger.Info(marconi_net_commands.TRAFFIC_CONTROL+" "+marconi_net_commands.SET_DELAY, util.ArgsToString(args))
  marconi_net_commands.HandleSetDelay(args)
}

func (mnm *MarconiNetMode) handleSetDuplicate(args []string) {
  util.Logger.Info(marconi_net_commands.TRAFFIC_CONTROL+" "+marconi_net_commands.SET_DUPLICATE, util.ArgsToString(args))
  marconi_net_commands.HandleSetDuplicate(args)
}

func (mnm *MarconiNetMode) handleSetLoss(args []string) {
  util.Logger.Info(marconi_net_commands.TRAFFIC_CONTROL+" "+marconi_net_commands.SET_LOSS, util.ArgsToString(args))
  marconi_net_commands.HandleSetLoss(args)
}

func (mnm *MarconiNetMode) handleSetReorder(args []string) {
  util.Logger.Info(marconi_net_commands.TRAFFIC_CONTROL+" "+marconi_net_commands.SET_REORDER, util.ArgsToString(args))
  marconi_net_commands.HandleSetReorder(args)
}

func (mnm *MarconiNetMode) handleReset(args []string) {
  util.Logger.Info(marconi_net_commands.TRAFFIC_CONTROL+" "+marconi_net_commands.RESET, util.ArgsToString(args))
  marconi_net_commands.HandleReset(args)
}
