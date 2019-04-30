package credentials

import (
  "bytes"
  "encoding/json"
  "fmt"
  "github.com/pkg/errors"
  "gitlab.neji.vm.tc/marconi/go-ethereum/params"
  "gitlab.neji.vm.tc/marconi/go-prompt"
  "gitlab.neji.vm.tc/marconi/cli/console/modes"
  "gitlab.neji.vm.tc/marconi/cli/console/util"
  "gitlab.neji.vm.tc/marconi/cli/core/configs"
  "gitlab.neji.vm.tc/marconi/cli/core/mkey"
  "io/ioutil"
  "math/big"
  "strconv"
  "strings"
)

const (
  CREATE_ACCOUNT          = "create"
  UNLOCK_ACCOUNT          = "unlock"
  LIST_ACCOUNTS           = "list"
  GET_BALANCE             = "balance"
  SEND_TRANSACTION        = "send"
  GET_TRANSACTION_RECEIPT = "receipt"
  EXPORT_GMRC_KEY         = "export"

  METH_DATA_CHILD_DIR = "/etc/meth/datadir"
  MIDDLEWARE_CONF_CHILD_PATH = "/etc/middleware/middleware_conf.json"
  EMPTY_USER_ADDRESS = "0x0000000000000000000000000000000000000000"
)

type jsonObject = map[string]interface{}

// Mode suggestions
var METH_ACCOUNT_SUGGESTIONS = []prompt.Suggest{
  {Text: CREATE_ACCOUNT, Description: "Create account"},
  {Text: UNLOCK_ACCOUNT, Description: "Unlock account"},
  {Text: LIST_ACCOUNTS, Description: "List accounts"},
  {Text: SEND_TRANSACTION, Description: "Send a transaction"},
  {Text: GET_BALANCE, Description: "Get balance for an account"},
  {Text: GET_TRANSACTION_RECEIPT, Description: "Get receipt for a transaction"},
  {Text: EXPORT_GMRC_KEY, Description: "Export Go Marconi Keystore associated with an account"},
}

/*
  Show prompt suggestions for create account command
*/
func (mm *CredsMode) getCreateAccountSuggestions(line []string) []prompt.Suggest {
  return []prompt.Suggest{}
}

/*
  Show prompt suggestions for the unlock account command
*/
func (mm *CredsMode) getUnlockAccountSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<0xACCOUNT_ADDRESS>", Description: "The GoMarconi account to be unlocked"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Show prompt suggestions for listing all available accounts
*/
func (mm *CredsMode) getListAccountSuggestions(line []string) []prompt.Suggest {
  return []prompt.Suggest{}
}

/*
  Handle the create account command
*/
func (mm *CredsMode) handleCreateAccount(args []string) {
  if !modes.ArgsLenCheckWithOptional(args, 0, 1) {
    fmt.Println("Usage:", CREATE_ACCOUNT)
    return
  }

  var password string
  for {
    var cancelled bool
    password, cancelled = getPassword("Please enter a password for this account")
    if cancelled {
      return
    }
    passwordConfirm, cancelled := getPassword("Please confirm the password")
    if cancelled {
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

  // update middleware_conf.json with the new User Address
  confPath := configs.GetFullPath(MIDDLEWARE_CONF_CHILD_PATH)
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
            //fmt.Println("Updated middleware_conf.json with Account Address", key.Address.String())
          }
        } else {
          fmt.Println("Failed to write new json to middleware_conf.json", err)
        }
      } else {
        fmt.Println("Failed to parse json for middleware_conf.json", err)
      }
    }
  } else {
    fmt.Println("Failed to read middleware_conf.json", err)
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

/*
  Handle the unlock account command
*/
func (mm *CredsMode) handleUnlockAccount(args []string) {
  if !modes.ArgsLenCheckWithOptional(args, 1, 1) {
    fmt.Println("Usage:", UNLOCK_ACCOUNT, "<0xACCOUNT_ADDRESS>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }

  password, cancelled := getPassword("Please enter your account password")
  if cancelled {
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

  unlocked, err := mm.middlewareClient.UnlockAccount(args[0], password, unlockInterval)
  if err != nil {
    fmt.Println("Error: ", err)
  } else {
    fmt.Println("Account unlocked:", unlocked)
  }
}

/*
  Handle the list accounts command
*/
func (mm *CredsMode) handleListAccounts(args []string) {
  accounts := mkey.ListAccounts()
  fmt.Println("Accounts:")
  for idx, address := range accounts {
    fmt.Printf("%-24d %48s\n", idx, util.GetEIP55Address(address))
  }
}

/*
  Handle the get balance command
*/
func (mm *CredsMode) handleGetBalance(args []string) {
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", GET_BALANCE, "<0xACCOUNT_ADDRESS>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) {
    return
  }

  balance, err := mm.middlewareClient.GetBalance(args[0])
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

/*
  Handle the send transaction command
*/
func (mm *CredsMode) handleSendTransaction(args []string) {
  if !modes.ArgsLenCheck(args, 5) {
    fmt.Println("Usage:", SEND_TRANSACTION, "<0xACCOUNT_ADDRESS> <0xOTHER_ADDRESS> <AMOUNT_IN_MARCOS> <GAS_LIMIT> <GAS_PRICE>")
    return
  }
  if !modes.ArgAddressCheck(args[0]) ||
     !modes.ArgAddressCheck(args[1]) ||
     !modes.ArgFloatCheck(args[2]) ||
     !modes.ArgUInt64Check(args[3]) ||
     !modes.ArgBigIntCheck(args[4]) {
    return
  }

  amountToSend, _ := strconv.ParseFloat(args[2], 64)
  if amountToSend < 0.0 {
    fmt.Println("Cannot send a negative amount.")
    return
  }

  balance, err := mm.middlewareClient.GetBalance(args[0])
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
  confirmed := modes.GetConfirmationInput()
  if !confirmed {
    fmt.Println("Transaction was cancelled")
    return
  }

  fmt.Print(": ")
  password, cancelled := getPassword("Please enter your account password")
  if cancelled {
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

  nonce, err := mm.middlewareClient.GetTransactionCount(args[0])
  if err != nil {
    fmt.Println("Error:", err)
  } else {
    fmt.Println("This may take up to a minute...")
    txHash, err := mm.middlewareClient.SendTransaction(
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

/*
  Handle the get transaction receipt command
*/
func (mm *CredsMode) handleGetTransactionReceipt(args []string) {
  if !modes.ArgsLenCheck(args, 1) {
    fmt.Println("Usage:", GET_TRANSACTION_RECEIPT, "<transaction_hash>")
    return
  }
  if !modes.ArgTxHashCheck(args[0]) {
    return
  }

  receipt, err := mm.middlewareClient.GetTransactionReceipt(args[0])
  if err != nil {
    fmt.Println("Error:", err)
  } else {
    fmt.Println("Receipt:", receipt)
  }
}

/*
  Handle export command
*/
func (mm *CredsMode) handleExportGMrcKey(args []string) {
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
