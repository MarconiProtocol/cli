package credentials

import (
  "fmt"

  "gitlab.neji.vm.tc/marconi/go-prompt"
  "gitlab.neji.vm.tc/marconi/cli/api/middleware"
  "gitlab.neji.vm.tc/marconi/cli/console/context"
  "gitlab.neji.vm.tc/marconi/cli/console/modes"
  "gitlab.neji.vm.tc/marconi/cli/console/util"
)

const (
  ACCOUNT = "account"
  KEY     = "key"
)

// Mode suggestions
var METH_SUGGESTIONS = []prompt.Suggest{
  {Text: ACCOUNT, Description: "Account related commands"},
  {Text: KEY, Description: "Key related commands"},
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
  credsMode.RegisterCommand(ACCOUNT, credsMode.getAccountSuggestions, credsMode.handleAccountSubMode)
  credsMode.RegisterCommand(KEY, credsMode.getKeySuggestions, credsMode.handleKeySubMode)

  // suggestion and handler registrations
  credsMode.RegisterSubCommand(ACCOUNT, UNLOCK_ACCOUNT, credsMode.getUnlockAccountSuggestions, credsMode.handleUnlockAccount)
  credsMode.RegisterSubCommand(ACCOUNT, CREATE_ACCOUNT, credsMode.getCreateAccountSuggestions, credsMode.handleCreateAccount)
  credsMode.RegisterSubCommand(ACCOUNT, LIST_ACCOUNTS, credsMode.getListAccountSuggestions, credsMode.handleListAccounts)
  credsMode.RegisterSubCommand(ACCOUNT, GET_BALANCE, credsMode.getGetBalanceSuggestions, credsMode.handleGetBalance)
  credsMode.RegisterSubCommand(ACCOUNT, SEND_TRANSACTION, credsMode.getSendTransactionSuggestions, credsMode.handleSendTransaction)
  credsMode.RegisterSubCommand(ACCOUNT, GET_TRANSACTION_RECEIPT, credsMode.getGetTransactionReceiptSuggestions, credsMode.handleGetTransactionReceipt)
  credsMode.RegisterSubCommand(ACCOUNT, EXPORT_GMRC_KEY, credsMode.getExportGMrcKeySuggestions, credsMode.handleExportGMrcKey)

  credsMode.RegisterSubCommand(KEY, GENERATE_MP_KEY, credsMode.getGenerateMPKeySuggestions, credsMode.handleGenerateMPKey)
  credsMode.RegisterSubCommand(KEY, USE_MPKEY, credsMode.getUseMpkKeySuggestions, credsMode.handleUseMPKey)
  credsMode.RegisterSubCommand(KEY, EXPORT_MP_KEY, credsMode.getExportMPKeySuggestions, credsMode.handleExportMPKey)
  credsMode.RegisterSubCommand(KEY, LIST_MPKEY_HASHES, credsMode.getListMpkKeyHashesSuggestions, credsMode.handleListMPKeyHashes)

  credsMode.RegisterCommand(modes.RETURN_TO_ROOT, credsMode.GetEmptySuggestions, credsMode.HandleReturnToRoot)
  credsMode.RegisterCommand(modes.EXIT_CMD, credsMode.GetEmptySuggestions, credsMode.HandleExitCommand)

  credsMode.middlewareClient = middleware.GetClient()

  return &credsMode
}

func (mm *CredsMode) CliPrefix() (string, bool) {
  return mm.Name(), false
}

func (mm *CredsMode) Name() string {
  return "credential"
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
  util.HandleFurtherCommands(ACCOUNT, METH_ACCOUNT_SUGGESTIONS)
}

/*
  Handle the key sub menu command
*/
func (mm *CredsMode) handleKeySubMode(args []string) {
  util.HandleFurtherCommands(KEY, METH_KEY_SUGGESTIONS)
}

func getPassword(outputPrompt string) (string, bool) {
  fmt.Print(outputPrompt, ": ")

  earlyExit := make(chan struct{}, 1)
  cancelled := false

  cancel := func(*prompt.Buffer) {
    cancelled = true
    earlyExit <- struct{}{}
  }

  password := prompt.Input("", modes.PasswordCompleter,
    prompt.OptionHiddenInput(),
    prompt.OptionAddKeyBind(prompt.KeyBind{Key: prompt.ControlC, Fn: cancel}),
    prompt.OptionSetEarlyExit(earlyExit))

  return password, cancelled
}
