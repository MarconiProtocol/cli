package util

import (
  "github.com/MarconiProtocol/cli/console/execution/execution_flags"
  "github.com/MarconiProtocol/go-prompt"
  mlog "github.com/MarconiProtocol/log"
  "golang.org/x/crypto/sha3"

  "fmt"
  "strings"
  "unicode"
)

const (
  ADDRESS_PREFIX     = "0x"
  TRANSACTION_PREFIX = "0x"
  PUB_KEY_PREFIX     = "Nx"
)

var Logger *mlog.Mlog // log commands submitted by users

func SplitInput(input string, keepLastElement bool) []string {
  inputSlice := strings.Split(input, " ")

  // trim and remove empty arguments
  inputSliceView := inputSlice[:0]
  for i, e := range inputSlice {
    e = strings.TrimSpace(e)

    // we need to include the last element because in case it's a space, we
    // need to keep it so that the completion machinery knows to get
    // completions from the next word, instead of trying to complete the
    // current word
    if e != "" ||
      (keepLastElement && i == len(inputSlice)-1) {

      inputSliceView = append(inputSliceView, e)
    }
  }
  return inputSliceView
}

func SimpleSubcommandCompleter(line []string, index int, suggestions []prompt.Suggest) []prompt.Suggest {
  if len(line) != index+1 {
    return []prompt.Suggest{}
  }
  return prompt.FilterHasPrefix(suggestions, line[index], true)
}

func getEIP55Case(r rune, n byte) rune {
  if n >= 8 {
    return unicode.ToUpper(r)
  } else {
    return unicode.ToLower(r)
  }
}

// Looks for a specific value in the array, and returns its location returns -1 if not found
func FindIndexOfValueInArray(arr []string, value string) int {
  for index, a := range arr {
    if a == value {
      return index
    }
  }
  return -1 // If nothing is found, return -1
}

func GetEIP55Address(address string) string {
  // assume address properly prefixed and is the correct length
  unprefixed := strings.TrimPrefix(address, ADDRESS_PREFIX)

  hash := sha3.NewLegacyKeccak256()
  hash.Write([]byte(strings.ToLower(unprefixed)))
  hashed := hash.Sum(nil)

  var builder strings.Builder
  builder.Grow(42)
  builder.WriteString(ADDRESS_PREFIX)

  for i, c := range unprefixed {
    nybble := hashed[i/2]

    if i%2 == 0 {
      nybble >>= 4
    } else {
      nybble &= 0x0f
    }

    builder.WriteRune(getEIP55Case(c, nybble))
  }

  return builder.String()
}

func HandleFurtherCommands(command string, suggestions []prompt.Suggest) {
  fmt.Println("Usage:", command, "has further commands")
  for _, c := range suggestions {
    fmt.Printf("    %-32s : %-32s\n", c.Text, c.Description)
  }
}

// Removes the password
func scrubArgs(args []string) []string {
  passwordFlagIndex := FindIndexOfValueInArray(args, execution_flags.PASSWORD)
  scrubbedArgs := make([]string, len(args))
  copy(scrubbedArgs, args)
  if passwordFlagIndex != -1 {
    if passwordFlagIndex+1 < len(args) {
      return append(scrubbedArgs[:passwordFlagIndex], scrubbedArgs[passwordFlagIndex+2:]...)
    }
    return append(scrubbedArgs[:passwordFlagIndex], scrubbedArgs[passwordFlagIndex+1:]...)
  }
  return args
}

func ArgsToString(args []string) string {
  scrubbedArgs := scrubArgs(args)
  str := strings.Join(scrubbedArgs, " ")
  if str == "" {
    return str
  } else {
    return " " + str
  }
}
