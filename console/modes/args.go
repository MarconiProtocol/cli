package modes

import (
  "fmt"
  "math/big"
  "os"
  "strconv"
  "strings"

  "github.com/MarconiProtocol/cli/console/util"
)

// TODO in the future, we should have more validations here
// so we can catch user input errors

func ArgsLenCheck(args []string, requiredLen int) bool {
  return len(args) == requiredLen
}

func ArgsLenCheckWithOptional(args []string, requiredLen int, optionalLen int) bool {
  return len(args) == requiredLen || len(args) == (requiredLen+optionalLen)
}

func eip55AddressCheck(arg string) bool {
  // don't run check on address that has only upper or lowercase
  if !(strings.ContainsAny(arg, "ABCDEF") && strings.ContainsAny(arg, "abcdef")) {
    return true
  }

  if arg != util.GetEIP55Address(arg) {
    fmt.Println("Mixed-case address failed EIP-55 checksum verification")
    return false
  }

  return true
}

func ArgAddressCheck(arg string) bool {
  if !(strings.HasPrefix(arg, util.ADDRESS_PREFIX) && len(strings.TrimPrefix(arg, util.ADDRESS_PREFIX)) == 40) {
    fmt.Println("Account address should start with 0x and have a length of 42")
    return false
  }

  if !eip55AddressCheck(arg) {
    return false
  }

  return true
}

func ArgTxHashCheck(arg string) bool {
  if !(strings.HasPrefix(arg, util.TRANSACTION_PREFIX) && len(strings.TrimPrefix(arg, util.TRANSACTION_PREFIX)) == 64) {
    fmt.Println("A tx hash should start with 0x and have a length of 64")
    return false
  }
  return true
}

func ArgPubKeyHashCheck(arg string) bool {
  if !(strings.HasPrefix(arg, util.PUB_KEY_PREFIX) && len(strings.TrimPrefix(arg, util.PUB_KEY_PREFIX)) == 40) {
    fmt.Println("A pub key hash should start with Nx and have a length of 42")
    return false
  }
  return true
}

func ArgDirExistsCheck(arg string) bool {
  if _, err := os.Stat(arg); os.IsNotExist(err) {
    fmt.Println(fmt.Sprint("Path: ", arg, "does not exist"))
    return false
  }
  return true
}

func ArgFloatCheck(arg string) bool {
  if _, err := strconv.ParseFloat(arg, 64); err != nil {
    fmt.Println("Argument", arg, "could not be parsed to a float")
    return false
  }
  return true
}

func ArgUInt64Check(arg string) bool {
  if _, err := strconv.ParseUint(arg, 10, 64); err != nil {
    fmt.Println("Argument", arg, "could not be parsed to a uint64")
    return false
  }
  return true
}

func ArgBigIntCheck(arg string) bool {
  var i big.Int
  if _, success := i.SetString(arg, 10); !success {
    fmt.Println("Argument", arg, "could not be parsed to a big int")
    return false
  }
  return true
}

func ArgBoolCheck(arg string) bool {
  if _, err := strconv.ParseBool(arg); err != nil {
    fmt.Println("Argument", arg, "could not be parsed to a boolean")
    return false
  }
  return true
}
