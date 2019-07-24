package mode_interface

import (
  "github.com/MarconiProtocol/go-prompt"
)

/*
  Interface for CLI Modes
*/
type Mode interface {
  // Returns suggestions for your cli command based on the current prompt
  Completer(d prompt.Document) []prompt.Suggest
  // Handle the typed in prompt selection
  HandleSelection(currCmd Mode, selection string, args []string)
  // Defines a custom cli prompt prefixes
  CliPrefix() (string, bool)
  // The name of the mode, will be used as the string in mode change commands
  Name() string
}
