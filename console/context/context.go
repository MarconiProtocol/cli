package context

import (
  "errors"
  "fmt"
  "github.com/MarconiProtocol/go-prompt"
  "github.com/MarconiProtocol/cli/console/mode_interface"
  "github.com/MarconiProtocol/cli/console/util"
)

// A context is used by the console for state and initialization
// TODO : finish scratchpad implementation
type Context struct {
  registeredModes map[string]mode_interface.Mode
  modeSuggestions []prompt.Suggest
  currentMode     mode_interface.Mode
  blackboard      map[string]interface{}
}

// Create a new context
func NewContext() *Context {
  context := Context{}
  context.registeredModes = make(map[string]mode_interface.Mode)
  context.blackboard = make(map[string]interface{})
  return &context
}

// Register a mode that will be usable within the context
func (c *Context) RegisterMode(m mode_interface.Mode, desc string) {
  c.registeredModes[m.Name()] = m
  c.modeSuggestions = append(c.modeSuggestions, prompt.Suggest{Text: m.Name(), Description: desc})
}

// Select the active mode, based on the mode name
func (c *Context) SelectMode(modeName string) (mode_interface.Mode, error) {
  if mode, exists := c.registeredModes[modeName]; exists {
    c.currentMode = mode
    return mode, nil
  }
  return nil, errors.New(fmt.Sprintf("Could not find a registered mode with name [%s], did not select mode", modeName))
}

// Returns the current active mode
func (c *Context) GetCurrentMode() (mode_interface.Mode, error) {
  if c.currentMode == nil {
    return nil, errors.New("No mode currently selected in context")
  }
  return c.currentMode, nil
}

// Returns a slice of Prompt.Suggest objects containing suggestions for the currently registered modes
func (c *Context) ModeSuggestionsFunc(line []string) []prompt.Suggest {
  return util.SimpleSubcommandCompleter(line, 1, c.modeSuggestions)
}

// When the active mode is set to nil, the cli will exit
func (c *Context) Exit() {
  fmt.Println("Bye!")
  c.currentMode = nil
}
