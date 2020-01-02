package credentials

import (
  "fmt"
  "github.com/MarconiProtocol/cli/api/middleware"
  "github.com/MarconiProtocol/cli/console/context"
  "github.com/MarconiProtocol/cli/console/modes"
  "github.com/MarconiProtocol/cli/console/modes/credentials/commands"
  "github.com/MarconiProtocol/cli/console/util"
  "github.com/MarconiProtocol/go-prompt"
)

const (
  NAME = "credential"
)

// Mode suggestions
var METH_SUGGESTIONS = []prompt.Suggest{
  {Text: credential_commands.ACCOUNT, Description: "Account related commands"},
  {Text: credential_commands.KEY, Description: "Key related commands"},
  {Text: modes.RETURN_TO_ROOT, Description: "Return to home menu"},
  {Text: modes.EXIT_CMD, Description: "Exit mcli"},
}

/*
  Menu/Mode for m-eth related operations
*/
type CredsMode struct {
  modes.BaseMode
  middlewareClient *middleware.Client
}

/*
  Create a new m-eth cmd struct,
  use this to create a CredsMode otherwise suggestion and handlers wont be properly initialized
*/
func NewCredsMode(c *context.Context) *CredsMode {
  credsMode := CredsMode{}
  credsMode.Init(c)
  credsMode.SetBaseSuggestions(METH_SUGGESTIONS)

  // suggestion and handler registrations
  credsMode.RegisterCommand(credential_commands.ACCOUNT, credsMode.getAccountSuggestions, credsMode.handleAccountSubMode)
  credsMode.RegisterCommand(credential_commands.KEY, credsMode.getKeySuggestions, credsMode.handleKeySubMode)

  // suggestion and handler registrations
  credsMode.RegisterSubCommand(credential_commands.ACCOUNT, credential_commands.UNLOCK_ACCOUNT, credsMode.getUnlockAccountSuggestions, credsMode.handleUnlockAccount)
  credsMode.RegisterSubCommand(credential_commands.ACCOUNT, credential_commands.CREATE_ACCOUNT, credsMode.getCreateAccountSuggestions, credsMode.handleCreateAccount)
  credsMode.RegisterSubCommand(credential_commands.ACCOUNT, credential_commands.LIST_ACCOUNTS, credsMode.getListAccountSuggestions, credsMode.handleListAccounts)
  credsMode.RegisterSubCommand(credential_commands.ACCOUNT, credential_commands.GET_BALANCE, credsMode.getGetBalanceSuggestions, credsMode.handleGetBalance)
  credsMode.RegisterSubCommand(credential_commands.ACCOUNT, credential_commands.SEND_TRANSACTION, credsMode.getSendTransactionSuggestions, credsMode.handleSendTransaction)
  credsMode.RegisterSubCommand(credential_commands.ACCOUNT, credential_commands.GET_TRANSACTION_RECEIPT, credsMode.getGetTransactionReceiptSuggestions, credsMode.handleGetTransactionReceipt)
  credsMode.RegisterSubCommand(credential_commands.ACCOUNT, credential_commands.EXPORT_GMRC_KEY, credsMode.getExportGMrcKeySuggestions, credsMode.handleExportGMrcKey)
  credsMode.RegisterSubCommand(credential_commands.ACCOUNT, credential_commands.USE_ACCOUNT, credsMode.getUseUserAddressSuggestions, credsMode.handleUseUserAddress)

  credsMode.RegisterSubCommand(credential_commands.KEY, credential_commands.GENERATE_MP_KEY, credsMode.getGenerateMPKeySuggestions, credsMode.handleGenerateMPKey)
  credsMode.RegisterSubCommand(credential_commands.KEY, credential_commands.USE_MPKEY, credsMode.getUseMpkKeySuggestions, credsMode.handleUseMPKey)
  credsMode.RegisterSubCommand(credential_commands.KEY, credential_commands.EXPORT_MP_KEY, credsMode.getExportMPKeySuggestions, credsMode.handleExportMPKey)
  credsMode.RegisterSubCommand(credential_commands.KEY, credential_commands.LIST_MPKEY_HASHES, credsMode.getListMpkKeyHashesSuggestions, credsMode.handleListMPKeyHashes)

  credsMode.RegisterCommand(modes.RETURN_TO_ROOT, credsMode.GetEmptySuggestions, credsMode.HandleReturnToRoot)
  credsMode.RegisterCommand(modes.EXIT_CMD, credsMode.GetEmptySuggestions, credsMode.HandleExitCommand)

  credsMode.middlewareClient = middleware.GetClient()

  return &credsMode
}

func (mm *CredsMode) CliPrefix() (string, bool) {
  return mm.Name(), false
}

func (mm *CredsMode) Name() string {
  return NAME
}

func (mm *CredsMode) HandleCommand(args []string) {
  if !modes.ArgsMinLenCheck(args, 2) {
    fmt.Println("USAGE: [account | key] <command>")
    return
  }
  commandType := args[0]
  commandArgs := args[1:]
  if commandHandlerFunction, present := credential_commands.COMMAND_MAP[commandType]; present {
    commandHandlerFunction(commandArgs)
  } else {
    fmt.Println("Invalid command " + commandType)
  }
}

/*
  Show prompt suggestions for entering the account sub menu
*/
func (mm *CredsMode) getAccountSuggestions(line []string) []prompt.Suggest {
  return util.SimpleSubcommandCompleter(line, 1, METH_ACCOUNT_SUGGESTIONS)
}

/*
  Show prompt suggestions for entering the key sub menu
*/
func (mm *CredsMode) getKeySuggestions(line []string) []prompt.Suggest {
  return util.SimpleSubcommandCompleter(line, 1, METH_KEY_SUGGESTIONS)
}

/*
  Show prompt suggestions for get balance command
*/
func (mm *CredsMode) getGetBalanceSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<0xACCOUNT_ADDRESS>", Description: "Your wallet address"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Show prompt suggestions for the send transaction command
*/
func (mm *CredsMode) getSendTransactionSuggestions(line []string) []prompt.Suggest {
  // fromAddress, toAddress, amount, password
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<0xACCOUNT_ADDRESS>", Description: "Your wallet address"}}
  case len(line) == 3:
    return []prompt.Suggest{{Text: "<0xTARGET_ADDRESS>", Description: "The target wallet address"}}
  case len(line) == 4:
    return []prompt.Suggest{{Text: "<AMOUNT_IN_MARCOS>", Description: "The amount of Marcos to send"}}
  case len(line) == 5:
    return []prompt.Suggest{{Text: "<GAS_LIMIT>", Description: "The gas limit"}}
  case len(line) == 6:
    return []prompt.Suggest{{Text: "<GAS_PRICE>", Description: "The gas price"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Show prompt suggestions for the get transaction receipt command
*/
func (mm *CredsMode) getGetTransactionReceiptSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<0xTRANSACTION_HASH>", Description: "Hash of the transaction to get receipt for"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Handle the account sub menu command
*/
func (mm *CredsMode) handleAccountSubMode(args []string) {
  util.HandleFurtherCommands(credential_commands.ACCOUNT, METH_ACCOUNT_SUGGESTIONS)
}

/*
  Handle the key sub menu command
*/
func (mm *CredsMode) handleKeySubMode(args []string) {
  util.HandleFurtherCommands(credential_commands.KEY, METH_KEY_SUGGESTIONS)
}
