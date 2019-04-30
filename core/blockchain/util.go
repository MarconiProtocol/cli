package blockchain

import (
  "gitlab.neji.vm.tc/marconi/go-ethereum/common"
  "gitlab.neji.vm.tc/marconi/go-ethereum/core/types"
  "gitlab.neji.vm.tc/marconi/cli/core/mkey"
  "math/big"

  "fmt"
)

const (
  CHAIN_ID = 161027
)

// Create a transaction with no data payload
func CreateTransaction(nonce uint64, toAddressStr string, amount *big.Int, gasLimit uint64, gasPrice *big.Int) *types.Transaction {
  toAddress := common.HexToAddress(toAddressStr)
  transaction := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, []byte{})
  return transaction
}

// Sign a given transaction with the Marconi account
func SignTransaction(mKeyStore *mkey.MarconiAccount, transaction *types.Transaction, password string) (*types.Transaction, error) {
  key, err := mKeyStore.GetGoMarconiKey(password)
  if err != nil {
    fmt.Println("Error loading GoMarconi Key", err)
    return nil, err
  }

  signedTransaction, err := types.SignTx(transaction, types.NewEIP155Signer(big.NewInt(CHAIN_ID)), key.PrivateKey)
  if err != nil {
    return nil, err
  }
  return signedTransaction, nil
}
