package mkey

import (
  "encoding/json"
  "errors"
  "fmt"
  "gitlab.neji.vm.tc/marconi/go-ethereum/accounts/keystore"
  "io/ioutil"
  "os"
  "path/filepath"
  "strconv"
  "strings"

  "gitlab.neji.vm.tc/marconi/cli/core/configs"
)

const (
  ACCOUNT_CHILD_DIR            = "/accounts"
  ACCOUNT_FILE_PREFIX          = "Account_Key"
  MARCONI_PRIVATE_KEY_FILENAME = "mpkey"
  MAX_MARCONI_PRIVATE_KEYS     = 16
  MARCONI_KEY_CHILD_DIR        = "/etc/marconid/keys"
)

/*
  MarconiAccount contains information for the GoMarconi Keystore and MarconiNode Keystores
*/
type MarconiAccount struct {
  filename     string                      // note: filename will not be marshalled to bytes
  GMrcKeystore keystore.EncryptedKeyJSONV3 `json:"gMrcKeystore"`
  MarconiKeys  []EncryptedMarconiKeyJSON   `json:"marconiKeys"`
}

/*
  Struct equivalent of the JSON format used to store an Encrypted Marconi Node Private Key
*/
type EncryptedMarconiKeyJSON struct {
  PublicKeyHash  string `json:"pubKeyHash"`
  EncryptedMPKey string `json:"encryptedMPkey"`
  Nonce          string `json:"nonce"`
  Salt           string `json:"salt"`
  N              int    `json:"n"`
  R              int    `json:"r"`
  P              int    `json:"p"`
  KeyLen         int    `json:"keylen"`
}

/*
  Create a new account, encrypted with password
*/
func CreateAccount(password string) (*MarconiAccount, error) {
  mKeystore, err := NewAccount(password)
  if err != nil {
    return nil, err
  }
  return mKeystore, nil
}

/*
  Returns the MarconiAccount associated with the provided address.
*/
func GetAccountForAddress(address string) (*MarconiAccount, error) {
  var matchingAccount *MarconiAccount
  f := func(m *MarconiAccount) bool {
    if strings.EqualFold(address, m.GMrcKeystore.Address) ||
      strings.EqualFold(strings.TrimPrefix(address, "0x"), m.GMrcKeystore.Address) {
      matchingAccount = m
      return false
    }
    return true
  }
  iterateAccounts(f)

  if matchingAccount == nil {
    return nil, errors.New(fmt.Sprintf("No keystore found with address %s", address))
  }
  return matchingAccount, nil
}

/*
  Returns the addresses of the Marconi accounts stored under the accounts directory
*/
func ListAccounts() []string {
  accounts := []string{}
  f := func(m *MarconiAccount) bool {
    accounts = append(accounts, m.GMrcKeystore.Address)
    return true
  }
  iterateAccounts(f)
  return accounts
}

type MpkeyListItem struct {
  Idx        int
  Pubkeyhash string
}

func GetMPKeyHashesForAddress(address string) ([]MpkeyListItem, error) {
  mpkeyList := []MpkeyListItem{}
  account, err := GetAccountForAddress(address)
  if err == nil {
    for idx, mpkey := range account.MarconiKeys {
      mpkeyList = append(mpkeyList, MpkeyListItem{idx, mpkey.PublicKeyHash})
    }
  }
  return mpkeyList, err
}

/*
  Helper function to iterate the accounts stored in the accounts directory,
  invoking the visitor function on each account.

  If the visitor function returns true, iterateAccounts continues to iterate, otherwise it will stop
*/
func iterateAccounts(visitor func(*MarconiAccount) bool) {
  // read the files in the account directory
  path := configs.GetFullPath(ACCOUNT_CHILD_DIR)
  files, err := ioutil.ReadDir(path)
  if err != nil {
    fmt.Println(err)
  }
  // check for account files
  for _, f := range files {
    // if the file matches the account file prefix, try to load it
    if strings.HasPrefix(f.Name(), ACCOUNT_FILE_PREFIX) {
      path := filepath.Join(path, f.Name())
      // try to loadKeystore file
      ks, err := LoadAccount(path)
      if err != nil {
        // ignore erroneous loads
        //fmt.Println(fmt.Sprintf("Error loading keystore file: %s, ignoring...", path))
        continue
      }
      if !visitor(ks) {
        break
      }
    }
  }
}

/*
  Load a Marconi account from a specific Marconi account file
*/
func LoadAccount(filename string) (*MarconiAccount, error) {
  mAcccountJSON, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }

  mAccount := MarconiAccount{}
  err = json.Unmarshal(mAcccountJSON, &mAccount)
  if err != nil {
    return nil, err
  }
  mAccount.filename = filename

  return &mAccount, nil
}

/*
  Generates a new Marconi account that will be encrypted with the given password
*/
func NewAccount(password string) (*MarconiAccount, error) {
  // Generate a new GoMarconi key
  key, err := GenerateGoMarconiKey()
  if err != nil {
    return nil, err
  }
  // Encrypt the key into an encrypted JSON format
  keyJSON, err := encryptGoMarconiKey(key, password)
  if err != nil {
    return nil, err
  }

  // Create a MarconiAccount and save it to disk
  mAccount := &MarconiAccount{}
  mAccount.filename = generateMarconiAccountFilename(keyJSON.Address)
  mAccount.GMrcKeystore = *keyJSON
  mAccount.MarconiKeys = []EncryptedMarconiKeyJSON{}
  mAccount.saveAccount()

  return mAccount, nil
}

/*
  Generates a new MarconiKey and adds it to the account
*/
// TODO: add a limit to the amount of keys we can add
func (m *MarconiAccount) GenerateMarconiKey(password string) (*EncryptedMarconiKeyJSON, error) {

  if len(m.MarconiKeys) >= MAX_MARCONI_PRIVATE_KEYS {
    return nil, errors.New("Error: Reached maximum number of Marconi Private Keys that can be stored in one account")
  }

  // Generate and encrypt new Marconi private key
  mpkey, err := generateMarconiKey()
  if err != nil {
    return nil, err
  }
  encryptedMPKeyJson, err := encryptMarconiKey(mpkey, password)
  if err != nil {
    return nil, err
  }

  // Add new Marconi private key to account and update the account file
  m.MarconiKeys = append(m.MarconiKeys, *encryptedMPKeyJson)
  m.saveAccount()

  return encryptedMPKeyJson, nil
}

/*
  Exports MarconiKeys stored in the account to mpkey files in the keystore directory
*/
func (m *MarconiAccount) ExportMarconiKeys(password string) error {
  exportDir := filepath.Join(configs.GetFullPath(ACCOUNT_CHILD_DIR), filepath.Base(m.filename)+"_mpkeys")

  count := 0
  for idx, encryptedMarconiKey := range m.MarconiKeys {
    marconiKey, err := decryptMarconiKey(&encryptedMarconiKey, password)
    if err != nil {
      return err
    }

    if err := os.MkdirAll(exportDir, 0700); err != nil {
      return err
    }

    savePrivateKey(filepath.Join(exportDir, MARCONI_PRIVATE_KEY_FILENAME+strconv.Itoa(idx)), marconiKey)
    savePublicKey(filepath.Join(exportDir, MARCONI_PRIVATE_KEY_FILENAME+strconv.Itoa(idx)+".pub"), &marconiKey.PublicKey)

    count++
  }
  if count > 0 {
    fmt.Println("\nExported", count, "keys to:")
    fmt.Println(exportDir)
  }
  return nil
}

func (m *MarconiAccount) UseMarconiKey(mpkeyId string, password string) error {
  if !strings.HasPrefix(mpkeyId, MARCONI_PRIVATE_KEY_FILENAME) {
    return errors.New("Entered nodekey is incorrect")
  }
  idxStr := strings.TrimPrefix(mpkeyId, MARCONI_PRIVATE_KEY_FILENAME)
  idx, err := strconv.Atoi(idxStr)
  if err != nil {
    return err
  }
  encryptedMarconiKey := m.MarconiKeys[idx]
  marconiKey, err := decryptMarconiKey(&encryptedMarconiKey, password)
  if err != nil {
    return err
  }
  path := configs.GetFullPath(MARCONI_KEY_CHILD_DIR)
  if err := os.MkdirAll(path, 0700); err != nil {
    return err
  }
  savePrivateKey(filepath.Join(path, MARCONI_PRIVATE_KEY_FILENAME), marconiKey)
  savePublicKey(filepath.Join(path, MARCONI_PRIVATE_KEY_FILENAME+".pub"), &marconiKey.PublicKey)
  //fmt.Println("Marconi Key", mpkeyId, "exported to:", path)
  return nil
}

/*
  Returns the GoMarconi key stored in the account
*/
func (m *MarconiAccount) GetGoMarconiKey(password string) (*keystore.Key, error) {
  key, err := decryptGoMarconiKey(&m.GMrcKeystore, password)
  if err != nil {
    return nil, err
  }
  return key, nil
}

/*
  Export and save the go marconi keystore to a path. Keys are already encrypted
  so can be written directly to disk. It's the same format as normal geth keys.
*/

func (m *MarconiAccount) ExportGoMarconiKeystore(path string) error {
  bytes, err := json.Marshal(m.GMrcKeystore)
  if err != nil {
    return err
  }
  filename := filepath.Join(path, "keystore", filepath.Base(m.filename))
  absFilename, err := filepath.Abs(filename)
  fmt.Println("Exported to:")
  fmt.Printf("%s\n", absFilename)
  return saveToFile(bytes, filename)
}

/*
  Saves the account data to disk
*/
func (m *MarconiAccount) saveAccount() error {
  bytes, _ := json.Marshal(m)
  return saveToFile(bytes, m.filename)
}

/*
  Generate a new account filename based on the provided address and the current timestamp
*/
func generateMarconiAccountFilename(address string) string {
  filename := fmt.Sprintf("%s-0x%s", ACCOUNT_FILE_PREFIX, address)
  path := filepath.Join(configs.GetFullPath(ACCOUNT_CHILD_DIR), filename)
  return path
}
