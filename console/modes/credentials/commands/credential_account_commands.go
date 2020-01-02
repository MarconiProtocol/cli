package credential_commands

import (
  "bytes"
  "encoding/json"
  "fmt"
  "github.com/MarconiProtocol/cli/api/middleware"
  "github.com/MarconiProtocol/cli/console/execution/execution_flags"
  "github.com/MarconiProtocol/cli/console/modes"
  "github.com/MarconiProtocol/cli/console/util"
  "github.com/MarconiProtocol/cli/core/configs"
  "github.com/MarconiProtocol/cli/core/mkey"
  "github.com/MarconiProtocol/go-methereum-lite/params"
  "github.com/pkg/errors"
  "io/ioutil"
  "log"
  "math/big"
  "os"
  "strconv"
  "strings"
)

// Commands
const (
  CREATE_ACCOUNT          = "create"
  UNLOCK_ACCOUNT          = "unlock"
  LIST_ACCOUNTS           = "list"
  GET_BALANCE             = "balance"
  SEND_TRANSACTION        = "send"
  GET_TRANSACTION_RECEIPT = "receipt"
  EXPORT_GMRC_KEY         = "export"
  USE_ACCOUNT             = "use"
)

const (
  METH_DATA_CHILD_DIR        = "/etc/meth/datadir"
  MIDDLEWARE_CONF_CHILD_PATH = "/etc/middleware/user_conf.json"
  EMPTY_USER_ADDRESS         = "0x0000000000000000000000000000000000000000"
)

var ACCOUNT_COMMAND_MAP = map[string]func([]string){
  CREATE_ACCOUNT:          CreateAccount,
  UNLOCK_ACCOUNT:          UnlockAccount,
  LIST_ACCOUNTS:           ListAccounts,
  GET_BALANCE:             GetBalance,
  SEND_TRANSACTION:        SendTransaction,
  GET_TRANSACTION_RECEIPT: GetTransactionReceipt,
  EXPORT_GMRC_KEY:         ExportGMrcKey,
  USE_ACCOUNT:             UseUserAddress,
}

func HandleAccountCommand(args []string) {
  if !modes.ArgsMinLenCheck(args, 1) {
    fmt.Println("USAGE: <command>")
    return
  }
  commandType := args[0]
  commandArgs := args[1:]
  if commandHandlerFunction, present := ACCOUNT_COMMAND_MAP[commandType]; present {
    commandHandlerFunction(commandArgs)
  } else {
    fmt.Println("Invalid command " + commandType)
  }
}

func CreateAccount(args []string) {
  if !modes.ArgsLenCheckWithOptional(args, 0, 2) {
    fmt.Println("Usage:", CREATE_ACCOUNT, "[Optional:", execution_flags.PASSWORD, "<password> |", execution_flags.PASSWORD_FILE, "<password file> ]")
    return
  }

  executionFlags := execution_flags.NewExecFlags(args)

  var password string
  var err error
  for {
    var cancelled bool
    password, cancelled, err = getPassword("Please enter a password for this account", executionFlags)
    if cancelled || err != nil {
      return
    }
    passwordConfirm, cancelled, err := getPassword("Please confirm the password", executionFlags)
    if cancelled || err != nil {
      return
    }
    if password == passwordConfirm {
      break
    }
    fmt.Println("Passwords did not match, please try again")
  }

  marconiAccount, err := mkey.CreateAccount(password)
  if err != nil {
    fmt.Println("Error creating an account", err)
    return
  }
  key, err := marconiAccount.GetGoMarconiKey(password)
  if err != nil {
    fmt.Println("Error fetching Marconi Key from account", err)
    return
  }
  fmt.Printf("Address:\n")
  fmt.Printf("%s\n", key.Address.String())

  // generate default middleware conf if it does not already exist
  confPath := configs.GetFullPath(MIDDLEWARE_CONF_CHILD_PATH)
  _, errConfig := os.Stat(confPath)
  if os.IsNotExist(errConfig) {
    // TODO hacky for now, need to better way to avoid Mcli from needing to write/update Middleware's config
    defaultConfig := []byte(
        "{\n" +
            "  \"meth\": {\n" +
            "    \"UserAddress\": \"0x0000000000000000000000000000000000000000\"\n" +
            "  }\n" +
            "}\n")
    errConfig = ioutil.WriteFile(confPath, defaultConfig, 0644)
    if errConfig != nil {
      log.Fatal("Failed to generate default middleware config: ", errConfig)
      return
    }
  }

  // update user_conf.json with the new User Address
  data, err := ioutil.ReadFile(confPath)
  if err == nil {
    var obj jsonObject
    err := json.Unmarshal(data, &obj)
    if err == nil && strings.Compare(obj["meth"].(jsonObject)["UserAddress"].(string), EMPTY_USER_ADDRESS) == 0 {
      obj["meth"].(jsonObject)["UserAddress"] = key.Address.String()
      newJson, err := json.Marshal(obj)
      if err == nil {
        var prettyJson bytes.Buffer
        err = json.Indent(&prettyJson, newJson, "", "  ")
        if err == nil {
          err = ioutil.WriteFile(confPath, prettyJson.Bytes(), 0644)
          if err == nil {
            //fmt.Println("Updated user_conf.json with Account Address", key.Address.String())
          }
        } else {
          fmt.Println("Failed to write new json to user_conf.json", err)
        }
      } else {
        fmt.Println("Failed to parse json for user_conf.json", err)
      }
    }
  } else {
    fmt.Println("Failed to read user_conf.json", err)
  }

  // check if exportKey argument is provided, default is exportKey = true
  exportKey := true
  if modes.ArgsLenCheck(args, 1) {
    value, err := strconv.ParseBool(args[0])
    if err == nil {
      exportKey = value
    }
  }
  if exportKey {
    keystore, err := mkey.GetAccountForAddress(key.Address.String())
    if err != nil {
      fmt.Println(err)
      return
    }

    err = keystore.ExportGoMarconiKeystore(configs.GetFullPath(METH_DATA_CHILD_DIR))
    if err != nil {
      fmt.Println("Failed to export go Marconi keystore", err)
      return
    }
  }
}

func UnlockAccount(args []string) {
  if !modes.ArgsLenCheckWithOptional(args, 1, 2) {
    fmt.Println("Usage:", UNLOCK_ACCOUNT, "<0xACCOUNT_ADDRESS> [Optional:", execution_flags.PASSWORD, "<password> |", execution_flags.PASSWORD_FILE, " <password file> ]")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }
  executionFlags := execution_flags.NewExecFlags(args)

  password, cancelled, err := getPassword("Please enter your account password", executionFlags)
  if cancelled || err != nil {
    return
  }

  unlockInterval := 900 // 15 minutes
  if len(args) == 2 {
    interval, err := strconv.Atoi(args[1])
    if err != nil {
      fmt.Println("Failed to parse unlock interval", err)
    } else {
      unlockInterval = interval
    }
  }

  unlocked, err := middleware.GetClient().UnlockAccount(args[0], password, unlockInterval)
  if err != nil {
    fmt.Println("Error: ", err)
  } else {
    fmt.Println("Account unlocked:", unlocked)
  }
}

func ListAccounts(args []string) {
  accounts := mkey.ListAccounts()
  for idx, address := range accounts {
    fmt.Printf("%-24d %48s\n", idx, util.GetEIP55Address(address))
  }
}

func GetBalance(args []string) {
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", GET_BALANCE, "<0xACCOUNT_ADDRESS>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }

  balance, err := middleware.GetClient().GetBalance(args[0])
  if err != nil {
    fmt.Println("Error: ", err)
  } else {
    fmt.Println("Balance:")
    marcos, err := gaussStrToMarcos(balance)
    if err == nil {
      fmt.Printf("%-24s %48f\n", "In Marcos", marcos)
    }
    fmt.Printf("%-24s %48s\n\n", "In Gauss", balance)
  }
}

func SendTransaction(args []string) {
  if !modes.ArgsLenCheckWithOptionalRange(args, 5, 2, 3) {
    fmt.Println("Usage:", SEND_TRANSACTION, "<0xACCOUNT_ADDRESS> <0xOTHER_ADDRESS> <AMOUNT_IN_MARCOS> <GAS_LIMIT> <GAS_PRICE> [Optional:", execution_flags.PASSWORD, "<password> |", execution_flags.PASSWORD_FILE, "<password file> |", execution_flags.SKIP_PROMPT_USE_DEFAULTS, "]")
    return
  }
  if !modes.ArgAddressCheck(args[0]) ||
      !modes.ArgAddressCheck(args[1]) ||
      !modes.ArgFloatCheck(args[2]) ||
      !modes.ArgUInt64Check(args[3]) ||
      !modes.ArgBigIntCheck(args[4]) {
    return
  }

  executionFlags := execution_flags.NewExecFlags(args)

  amountToSend, _ := strconv.ParseFloat(args[2], 64)
  if amountToSend < 0.0 {
    fmt.Println("Cannot send a negative amount.")
    return
  }

  balance, err := middleware.GetClient().GetBalance(args[0])
  if err != nil {
    fmt.Println("Failed to get balance: ", err)
    return
  }

  if gaussBalance, _ := gaussStrToMarcos(balance); gaussBalance.Cmp(big.NewFloat(amountToSend)) < 0 {
    fmt.Println("Insufficient balance, please double check the amount.")
    return
  }

  // Transaction summary
  fmt.Println("Please confirm the transaction:")
  fmt.Printf("%-16s: %48s\n", "From Address", args[0])
  fmt.Printf("%-16s: %48s\n", "To Address", args[1])
  fmt.Printf("%-16s: %48s\n", "Marcos to Send", args[2])
  fmt.Printf("%-16s: %48s\n", "Gas Limit", args[3])
  fmt.Printf("%-16s: %48s\n", "Gas Price", args[4])

  // If the Skip flag isn't set, then ask for confirmation
  if !executionFlags.CheckSkipPromptsFlagSet() {
    confirmed := modes.GetConfirmationInput()
    if !confirmed {
      fmt.Println("Transaction was cancelled")
      return
    }
  }

  fmt.Print(": ")
  password, cancelled, err := getPassword("Please enter your account password", executionFlags)
  if cancelled || err != nil {
    return
  }

  // amount is expected in ether
  amount, _ := strconv.ParseFloat(args[2], 64)
  amountInEtherFloat := new(big.Float).SetFloat64(amount)
  conversionFactor := new(big.Float).SetInt(big.NewInt(params.Ether))

  // NOTE: int value of wei is an approximation due to floating point operations
  amountInWeiFloat := new(big.Float).Mul(amountInEtherFloat, conversionFactor)
  amountInWeiInt := new(big.Int)
  amountInWeiFloat.Int(amountInWeiInt)

  gasLimit, _ := strconv.ParseUint(args[3], 10, 64)
  gasPrice, _ := new(big.Int).SetString(args[4], 10)

  nonce, err := middleware.GetClient().GetTransactionCount(args[0])
  if err != nil {
    fmt.Println("Error:", err)
  } else {
    fmt.Println("This may take up to a minute...")
    txHash, err := middleware.GetClient().SendTransaction(
      password, // password
      uint64(nonce),
      args[0],        // fromAddress
      args[1],        // toAddress
      amountInWeiInt, // amount in wei
      gasLimit,
      gasPrice,
    )
    if err != nil {
      fmt.Println("Error:", err)
    } else {
      fmt.Println("Transaction Hash:", txHash)
    }
  }
}

func GetTransactionReceipt(args []string) {
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", GET_TRANSACTION_RECEIPT, "<transaction_hash>")
    return
  }
  if !modes.ArgTxHashCheck(args[0]) {
    return
  }

  receipt, err := middleware.GetClient().GetTransactionReceipt(args[0])
  if err != nil {
    fmt.Println("Error:", err)
  } else {
    fmt.Println("Receipt:", receipt)
  }
}

func ExportGMrcKey(args []string) {
  if !modes.ArgsLenCheckWithOptional(args, 1, 1) {
    fmt.Println("Usage:", EXPORT_GMRC_KEY, "<0xACCOUNT_ADDRESS> <GO-MARCONI_DATA_DIR_PATH>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }

  exportDir := configs.GetFullPath(METH_DATA_CHILD_DIR)
  if modes.ArgsLenCheck(args, 2) {
    exportDir = args[1]
  }

  if !modes.ArgDirExistsCheck(exportDir) {
    fmt.Println(exportDir, "doesn't exist")
    return
  }

  keystore, err := mkey.GetAccountForAddress(args[0])
  if err != nil {
    fmt.Println(err)
    return
  }

  err = keystore.ExportGoMarconiKeystore(exportDir)
  if err != nil {
    fmt.Println("Failed to export go Marconi keystore", err)
    return
  }
  fmt.Println("go Marconi Keystore successfully exported")
  return
}

func UseUserAddress(args []string) {

  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", ACCOUNT, USE_ACCOUNT, "<0xACCOUNT_ADDRESS>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }

  success, err := middleware.GetClient().UpdateUserAddress(args[0])
  if err != nil {
    fmt.Println("Error: ", err)
  } else if !success {
    fmt.Println("Use account address failed")
  } else {
    fmt.Println("Using account address", args[0])
  }
}

func gaussStrToMarcos(amountInGaussStr string) (*big.Float, error) {
  var amountInGauss *big.Float
  var success bool
  amountInGauss, success = new(big.Float).SetString(amountInGaussStr)
  if !success {
    return nil, errors.New("Error! Provided wei string could not be converted to a number")
  }
  amountInMarcos := new(big.Float).Quo(amountInGauss, big.NewFloat(params.Ether))
  return amountInMarcos, nil
}
