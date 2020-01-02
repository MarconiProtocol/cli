package mkey

import (
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  "crypto/rsa"
  "crypto/sha1"
  "crypto/x509"
  "encoding/base64"
  "encoding/hex"
  "encoding/json"
  "encoding/pem"
  "fmt"
  "github.com/MarconiProtocol/marconid/core/net/vars"
  "github.com/MarconiProtocol/go-methereum-lite/accounts/keystore"
  "github.com/MarconiProtocol/go-methereum-lite/crypto"
  "github.com/pkg/errors"
  "golang.org/x/crypto/scrypt"
  "io"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"
)

const (
  NODE_PREFIX = "Nx"
)

/*
  Generate a new GoMarconi private key
*/
func GenerateGoMarconiKey() (*keystore.Key, error) {
  privateKey, err := crypto.GenerateKey()
  if err != nil {
    return nil, err
  }
  key := keystore.NewKeyFromECDSA(privateKey)
  return key, nil
}

/*
  Encrypts a GoMarconi private key and stores it in the EncryptedKeyJSONV3 format
*/
func encryptGoMarconiKey(key *keystore.Key, password string) (*keystore.EncryptedKeyJSONV3, error) {
  keyJsonBytes, err := keystore.EncryptKey(key, password, keystore.StandardScryptN, keystore.StandardScryptP)
  if err != nil {
    return nil, err
  }
  keyJson := keystore.EncryptedKeyJSONV3{}
  err = json.Unmarshal(keyJsonBytes, &keyJson)
  if err != nil {
    return nil, err
  }
  return &keyJson, nil
}

/*
  Decrypts a EncryptedKeyJSONV3 with the provided password and returns a usable key
*/
func decryptGoMarconiKey(keystoreJSON *keystore.EncryptedKeyJSONV3, password string) (*keystore.Key, error) {
  keyJsonBytes, err := json.Marshal(keystoreJSON)
  if err != nil {
    return nil, err
  }
  key, err := keystore.DecryptKey(keyJsonBytes, password)
  if err != nil {
    return nil, err
  }
  return key, nil
}

/*
  Generates a new RSA key, to be used as a Marconi Node key
*/
func generateMarconiKey() (*rsa.PrivateKey, error) {
  random := rand.Reader
  return rsa.GenerateKey(random, 2048)
}

/*
  Encrypts a Marconi Key with the provided password
*/
func encryptMarconiKey(marconiKey *rsa.PrivateKey, password string) (*EncryptedMarconiKeyJSON, error) {
  // use 32 bytes from rand.Reader as salt
  salt := make([]byte, 32)
  if _, err := io.ReadFull(rand.Reader, salt); err != nil {
    return nil, err
  }

  // scrypt.Key suggests a value of 32768 for 2017
  // lets just double that value for now
  N := 1 << 16
  r := 8
  p := 1
  keylen := 32
  // generate a 32byte key from a password
  encryptionKey, err := scrypt.Key([]byte(password), salt, N, r, p, keylen)
  if err != nil {
    return nil, err
  }

  // Get the byte data for marconiKey in ASN.1 DER encoded form
  mpkeyBytes := x509.MarshalPKCS1PrivateKey(marconiKey)

  block, err := aes.NewCipher(encryptionKey)
  if err != nil {
    return nil, err
  }
  // Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
  nonce := make([]byte, 12)
  if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
    return nil, err
  }
  aesgcm, err := cipher.NewGCM(block)
  if err != nil {
    return nil, err
  }

  // encrypt the marconiKey
  encryptedMpKeyBytes := aesgcm.Seal(nil, nonce, mpkeyBytes, nil)
  // calculate public key hash
  pubKeyHash, err := getInfohashByPubKey(&marconiKey.PublicKey)
  if err != nil {
    return nil, err
  }

  marconiKeyJSON := EncryptedMarconiKeyJSON{
    pubKeyHash,
    hex.EncodeToString(encryptedMpKeyBytes),
    hex.EncodeToString(nonce),
    hex.EncodeToString(salt),
    N,
    r,
    p,
    keylen,
  }

  return &marconiKeyJSON, nil
}

/*
  Decrypts an encrypted Marconi Key and returns the RSA private key that can be used for cryptographic operations
*/
func decryptMarconiKey(encryptedMpKey *EncryptedMarconiKeyJSON, password string) (*rsa.PrivateKey, error) {
  encryptedKeyBytes, err := hex.DecodeString(encryptedMpKey.EncryptedMPKey)
  nonce, err := hex.DecodeString(encryptedMpKey.Nonce)
  salt, err := hex.DecodeString(encryptedMpKey.Salt)

  // get key from password
  key, err := scrypt.Key([]byte(password), salt, encryptedMpKey.N, encryptedMpKey.R, encryptedMpKey.P, encryptedMpKey.KeyLen)
  if err != nil {
    fmt.Println(err)
  }

  block, err := aes.NewCipher(key)
  if err != nil {
    return nil, err
  }
  aesgcm, err := cipher.NewGCM(block)
  if err != nil {
    return nil, err
  }

  keyBytes, err := aesgcm.Open(nil, nonce, encryptedKeyBytes, nil)
  if err != nil {
    return nil, err
  }

  // parse DER format to a native type
  mpkey, err := x509.ParsePKCS1PrivateKey(keyBytes)
  if err != nil {
    return nil, errors.New("Error parsing key bytes")
  }
  return mpkey, nil
}

func saveToFile(bytes []byte, filename string) error {
  const dirPerm = 0700
  if err := os.MkdirAll(filepath.Dir(filename), dirPerm); err != nil {
    return err
  }
  // Create and write to a temp file, rename when done
  f, err := ioutil.TempFile(filepath.Dir(filename), "."+filepath.Base(filename)+".tmp")
  if err != nil {
    return err
  }
  if _, err := f.Write(bytes); err != nil {
    f.Close()
    os.Remove(f.Name())
    return err
  }
  f.Close()
  return os.Rename(f.Name(), filename)
}

// todo : copied from marconid, do we want to read/write keys to an encrypted format?
func savePrivateKey(filename string, key *rsa.PrivateKey) {
  keyFile, err := os.Create(filename)
  if err != nil {
    fmt.Println("Failed to save key to file: ", filename)
  }
  defer keyFile.Close()

  var privateKey = &pem.Block{
    Type:  "RSA PRIVATE KEY",
    Bytes: x509.MarshalPKCS1PrivateKey(key),
  }

  err = pem.Encode(keyFile, privateKey)
}

func savePublicKey(filename string, key *rsa.PublicKey) {
  ans1Bytes, err := x509.MarshalPKIXPublicKey(key)
  if err != nil {
    fmt.Println("Failed to marshal public key", err)
  }

  var pemkey = &pem.Block{
    Type:  "PUBLIC KEY",
    Bytes: ans1Bytes,
  }

  keyFile, err := os.Create(filename)
  if err != nil {
    fmt.Println("Failed to create file for key: ", filename, err)
  }
  defer keyFile.Close()

  err = pem.Encode(keyFile, pemkey)
  if err != nil {
    fmt.Println("Failed to save key to file: ", filename, err)
  }
}

func getInfohashByPubKey(pub *rsa.PublicKey) (string, error) {
  pubKeyBytes, err := x509.MarshalPKIXPublicKey(pub)
  if err != nil {
    return "", err
  }
  sha1Bytes := sha1.Sum(pubKeyBytes)
  return hex.EncodeToString(sha1Bytes[:]), nil
}

func Generate32BitKey(path string) (key []byte, e error) {
  /* Generate a random key */
  key = make([]byte, mnet_vars.HMAC_SHA256_SIZE)
  n, err := rand.Read(key)
  if n != len(key) {
    return nil, fmt.Errorf("Error generating random key of size %d: %s", len(key), err)
  }

  /* Base64 encode the key */
  key_base64 := make([]byte, base64.StdEncoding.EncodedLen(len(key)))
  base64.StdEncoding.Encode(key_base64, key)

  /* Write the base64 encoded key */
  err = saveToFile(key_base64, path)
  if err != nil {
    return nil, fmt.Errorf("Error writing base64 encoded key to keyfile: %s", err)
  }

  return key, nil
}

func AddPrefixPubKeyHash(pubkeyhash string) string {
  return NODE_PREFIX + pubkeyhash
}

func StripPrefixPubKeyHash(pubkeyhash string) string {
  if len(pubkeyhash) == 42 {
    return strings.TrimPrefix(pubkeyhash, NODE_PREFIX)
  }
  return pubkeyhash
}
