package credentials

import (
  "github.com/MarconiProtocol/cli/console/modes/credentials/commands"
  "github.com/MarconiProtocol/cli/console/util"
  "github.com/MarconiProtocol/go-prompt"
)

// Mode suggestions
var METH_ACCOUNT_SUGGESTIONS = []prompt.Suggest{
  {Text: credential_commands.CREATE_ACCOUNT, Description: "Create account"},
  {Text: credential_commands.UNLOCK_ACCOUNT, Description: "Unlock account"},
  {Text: credential_commands.LIST_ACCOUNTS, Description: "List accounts"},
  {Text: credential_commands.SEND_TRANSACTION, Description: "Send a transaction"},
  {Text: credential_commands.GET_BALANCE, Description: "Get balance for an account"},
  {Text: credential_commands.GET_TRANSACTION_RECEIPT, Description: "Get receipt for a transaction"},
  {Text: credential_commands.EXPORT_GMRC_KEY, Description: "Export Go Marconi Keystore associated with an account"},
  {Text: credential_commands.USE_ACCOUNT, Description: "Use account address"},
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
  util.Logger.Info(credential_commands.ACCOUNT+" "+credential_commands.CREATE_ACCOUNT, util.ArgsToString(args))
  credential_commands.CreateAccount(args)
}

/*
  Handle the unlock account command
*/
func (mm *CredsMode) handleUnlockAccount(args []string) {
  util.Logger.Info(credential_commands.ACCOUNT+" "+credential_commands.UNLOCK_ACCOUNT, util.ArgsToString(args))
  credential_commands.UnlockAccount(args)
}

/*
  Handle the list accounts command
*/
func (mm *CredsMode) handleListAccounts(args []string) {
  util.Logger.Info(credential_commands.ACCOUNT+" "+credential_commands.LIST_ACCOUNTS, util.ArgsToString(args))
  credential_commands.ListAccounts(args)
}

/*
  Handle the get balance command
*/
func (mm *CredsMode) handleGetBalance(args []string) {
  util.Logger.Info(credential_commands.ACCOUNT+" "+credential_commands.GET_BALANCE, util.ArgsToString(args))
  credential_commands.GetBalance(args)
}

/*
  Handle the send transaction command
*/
func (mm *CredsMode) handleSendTransaction(args []string) {
  util.Logger.Info(credential_commands.ACCOUNT+" "+credential_commands.SEND_TRANSACTION, util.ArgsToString(args))
  credential_commands.SendTransaction(args)
}

/*
  Handle the get transaction receipt command
*/
func (mm *CredsMode) handleGetTransactionReceipt(args []string) {
  util.Logger.Info(credential_commands.ACCOUNT+" "+credential_commands.GET_TRANSACTION_RECEIPT, util.ArgsToString(args))
  credential_commands.GetTransactionReceipt(args)
}

/*
  Handle export command
*/
func (mm *CredsMode) handleExportGMrcKey(args []string) {
  util.Logger.Info(credential_commands.ACCOUNT+" "+credential_commands.EXPORT_GMRC_KEY, util.ArgsToString(args))
  credential_commands.ExportGMrcKey(args)
}

/*
  Handle use command
*/
func (mm *CredsMode) handleUseUserAddress(args []string) {
  util.Logger.Info(credential_commands.ACCOUNT+" "+credential_commands.USE_ACCOUNT, util.ArgsToString(args))
  credential_commands.UseUserAddress(args)
}
