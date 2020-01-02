package credentials

import (
  "github.com/MarconiProtocol/cli/console/modes/credentials/commands"
  "github.com/MarconiProtocol/cli/console/util"
  "github.com/MarconiProtocol/go-prompt"
)

// Mode suggestions
var METH_KEY_SUGGESTIONS = []prompt.Suggest{
  {Text: credential_commands.GENERATE_MP_KEY, Description: "Generate nodekey"},
  {Text: credential_commands.USE_MPKEY, Description: "Set nodekey to use with other commands"},
  {Text: credential_commands.EXPORT_MP_KEY, Description: "Export nodekey"},
  {Text: credential_commands.LIST_MPKEY_HASHES, Description: "List nodekeys"},
}

/*
  Show prompt suggestions for generate Marconi Node Private key command
*/
func (mm *CredsMode) getGenerateMPKeySuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<0xACCOUNT_ADDRESS>", Description: "The GoMarconi account you wish to generate an nodekey for"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Show prompt suggestions for export Marconi Node Private key command
*/
func (mm *CredsMode) getExportMPKeySuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<0xACCOUNT_ADDRESS>", Description: "The GoMarconi account you wish to export nodekey for"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Show prompt suggestions for export marconi keystore command
*/
func (mm *CredsMode) getExportGMrcKeySuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<0xACCOUNT_ADDRESS>", Description: "The GoMarconi account you wish to export the go Marconi keystore for"}}
  case len(line) == 3:
    return []prompt.Suggest{{Text: "<GO-MARCONI_DATA_DIR_PATH>", Description: "The GoMarconi data directory path"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Show prompt suggestions for use userAddress command
*/
func (mm *CredsMode) getUseUserAddressSuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<0xACCOUNT_ADDRESS>", Description: "The account address you wish to use"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Show prompt suggestions for generating a 32 bit key
*/
func (mm *CredsMode) getGenerate32BitKeySuggestions(line []string) []prompt.Suggest {
  return []prompt.Suggest{}
}

/*
  Show prompt suggestions for showing mpkeys for an account
*/
func (mm *CredsMode) getListMpkKeyHashesSuggestions(line []string) []prompt.Suggest {
  return []prompt.Suggest{{Text: "<0xACCOUNT_ADDRESS>", Description: "The GoMarconi account whose nodekey you wish to view"}}
}

/*
  Show prompt suggestions for copying mpkeys for use with marconid
*/
func (mm *CredsMode) getUseMpkKeySuggestions(line []string) []prompt.Suggest {
  switch {
  case len(line) == 2:
    return []prompt.Suggest{{Text: "<0xACCOUNT_ADDRESS>", Description: "The GoMarconi account that has the nodekeys"}}
  default:
    return []prompt.Suggest{}
  }
}

/*
  Handle the  generate Marconi Node Private key command
*/
func (mm *CredsMode) handleGenerateMPKey(args []string) {
  util.Logger.Info(credential_commands.KEY+" "+credential_commands.GENERATE_MP_KEY, util.ArgsToString(args))
  credential_commands.GenerateMPKey(args)
}

/*
  Handle the export Marconi Node Private keys command
*/
func (mm *CredsMode) handleExportMPKey(args []string) {
  util.Logger.Info(credential_commands.KEY+" "+credential_commands.EXPORT_MP_KEY, util.ArgsToString(args))
  credential_commands.ExportMPKey(args)
}

func (mm *CredsMode) handleListMPKeyHashes(args []string) {
  util.Logger.Info(credential_commands.KEY+" "+credential_commands.LIST_MPKEY_HASHES, util.ArgsToString(args))
  credential_commands.ListMPKeyHashes(args)
}

func (mm *CredsMode) handleUseMPKey(args []string) {
  util.Logger.Info(credential_commands.KEY+" "+credential_commands.USE_MPKEY, util.ArgsToString(args))
  credential_commands.UseMPKey(args)
}
