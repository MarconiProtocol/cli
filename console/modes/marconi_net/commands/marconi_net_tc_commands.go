package marconi_net_commands

import (
  "fmt"
  "github.com/MarconiProtocol/cli/api/marconid"
  "github.com/MarconiProtocol/cli/console/modes"
  "os/exec"
  "strconv"
  "strings"
)

const (
  INFO          = "info"
  SET_BANDWIDTH = "set_bandwidth"
  SET_CORRUPT   = "set_corrupt"
  SET_DELAY     = "set_delay"
  SET_DUPLICATE = "set_duplicate"
  SET_LOSS      = "set_loss"
  SET_REORDER   = "set_reorder"
  RESET         = "reset"
)

var TC_COMMAND_MAP = map[string]func([]string){
  INFO:          HandleInfo,
  SET_BANDWIDTH: HandleSetBandwidth,
  SET_CORRUPT:   HandleSetCorrupt,
  SET_DELAY:     HandleSetDelay,
  SET_DUPLICATE: HandleSetDuplicate,
  SET_LOSS:      HandleSetLoss,
  SET_REORDER:   HandleSetReorder,
  RESET:         HandleReset,
}

func HandleTCCommand(args []string) {
  if !modes.ArgsMinLenCheck(args, 1) {
    fmt.Println("USAGE: <command>")
    return
  }
  commandType := args[0]
  commandArgs := args[1:]
  if commandHandlerFunction, present := TC_COMMAND_MAP[commandType]; present {
    commandHandlerFunction(commandArgs)
  } else {
    fmt.Println("Invalid command " + commandType)
  }
}

func HandleInfo(args []string) {
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", INFO, "<NET_INTERFACE>")
    return
  }
  cmd := "tc qdisc show dev " + args[0]
  if output, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
    fmt.Println(err)
  } else {
    fmt.Println(strings.Trim(string(output), "\n"))
  }
}

func HandleSetBandwidth(args []string) {
  if !modes.ArgsLenCheck(args, 3) {
    fmt.Println("Usage:", SET_BANDWIDTH, "<NET_INTERFACE>, <BANDWIDTH>, <LATENCY>")
    return
  }

  interfaceName := args[0]
  bandwidth, err := strconv.ParseUint(args[1], 10, 64)
  if err != nil {
    fmt.Println(err)
    return
  }
  latency, err := strconv.ParseFloat(args[2], 64)
  if err != nil {
    fmt.Println(err)
    return
  }
  if err := marconid.SetTbf(interfaceName, bandwidth, latency); err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Success")
  }
}

func HandleSetCorrupt(args []string) {
  if !modes.ArgsLenCheck(args, 2) {
    fmt.Println("Usage:", SET_CORRUPT, "<NET_INTERFACE>, <CORRUPT_RATE>")
    return
  }

  interfaceName := args[0]
  temp, err := strconv.ParseFloat(args[1], 32)
  if err != nil {
    fmt.Println(err)
    return
  }
  corruptProb := float32(temp)
  if err := marconid.SetNetem(interfaceName, 0, 0, 0, 0, corruptProb); err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Success")
  }
}

func HandleSetDelay(args []string) {
  if !modes.ArgsLenCheck(args, 2) {
    fmt.Println("Usage:", SET_DELAY, "<NET_INTERFACE>, <DELAY>")
    return
  }

  interfaceName := args[0]
  temp, err := strconv.ParseUint(args[1], 10, 32)
  if err != nil {
    fmt.Println(err)
    return
  }
  delay := uint32(temp)
  if err := marconid.SetNetem(interfaceName, delay, 0, 0, 0, 0); err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Success")
  }
}

func HandleSetDuplicate(args []string) {
  if !modes.ArgsLenCheck(args, 2) {
    fmt.Println("Usage:", SET_DUPLICATE, "<NET_INTERFACE>, <DUPLICATE_RATE>")
    return
  }

  interfaceName := args[0]
  temp, err := strconv.ParseFloat(args[1], 32)
  if err != nil {
    fmt.Println(err)
    return
  }
  duplicate := float32(temp)
  if err := marconid.SetNetem(interfaceName, 0, 0, duplicate, 0, 0); err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Success")
  }
}

func HandleSetLoss(args []string) {
  if !modes.ArgsLenCheck(args, 2) {
    fmt.Println("Usage:", SET_LOSS, "<NET_INTERFACE>, <LOSS_RATE>")
    return
  }

  interfaceName := args[0]
  temp, err := strconv.ParseFloat(args[1], 32)
  if err != nil {
    fmt.Println(err)
    return
  }
  loss := float32(temp)
  if err := marconid.SetNetem(interfaceName, 0, loss, 0, 0, 0); err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Success")
  }
}

func HandleSetReorder(args []string) {
  if !modes.ArgsLenCheck(args, 2) {
    fmt.Println("Usage:", SET_REORDER, "<NET_INTERFACE>, <REORDER_RATE>")
    return
  }

  interfaceName := args[0]
  temp, err := strconv.ParseFloat(args[1], 32)
  if err != nil {
    fmt.Println(err)
    return
  }
  reorderProb := float32(temp)
  if err := marconid.SetNetem(interfaceName, 0, 0, 0, reorderProb, 0); err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Success")
  }
}

func HandleReset(args []string) {
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", RESET, "<NET_INTERFACE>")
    return
  }

  if err := marconid.Reset(args[0]); err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Success")
  }
}
