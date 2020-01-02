package middleware

import (
  "bytes"
  "encoding/json"
  "fmt"
  "github.com/MarconiProtocol/cli/core/blockchain"
  "github.com/MarconiProtocol/cli/core/configs"
  "github.com/MarconiProtocol/cli/core/mkey"
  "math/big"
  "strconv"
)

const ETH_API_MIDDLEWARE_URL_PATH string = "/api/eth/v1"
const ETH_API_PERSONAL_URL_PATH string = "/api/personal/v1"
const MARCONI_API_MIDDLEWARE_URL_PATH string = "/api/marconi/v1"
const MIDDLEWARE_API_MIDDLEWARE_URL_PATH string = "/api/middleware/v1"

type Client struct {
  host string
  port string
}

func GetClient() *Client {
  conf := configs.LoadBaseConf()
  client := &Client{
    conf.MarconiNodeHost,
    conf.MarconiNodePort,
  }
  return client
}

/*
  Unlock the provided account address
*/
func (c *Client) UnlockAccount(address string, password string, interval int) (bool, error) {
  params := []string{
    address,
    password,
    strconv.Itoa(interval),
  }
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return false, err
  }

  payload := createJsonPayload("personal_unlockAccount", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, ETH_API_PERSONAL_URL_PATH, payload)
  if err != nil {
    return false, err
  }

  r := JSONRpcResponseBool{}
  err = json.Unmarshal(response, &r)
  if err != nil {
    return false, err
  }

  return r.Result, nil
}

/*
  Get the balance for the given address
*/
func (c *Client) GetBalance(address string) (string, error) {
  params := []string{
    address,
  }
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return "", err
  }

  payload := createJsonPayload("eth_getBalance", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, ETH_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return "", err
  }

  r := JSONRpcResponse{}
  err = json.Unmarshal(response, &r)
  if err != nil {
    return "", err
  }

  return hexStringToDecimalString(r.Result), nil
}

/*
  Register the given pub key hash
*/
func (c *Client) RegisterUser(pubKeyHash string, macHash string) (*RegisterResult, error) {
  params := map[string]string{
    "PubKeyHash": pubKeyHash,
    "MacHash":    macHash,
  }
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return nil, err
  }

  payload := createJsonPayload("registerUser", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, MARCONI_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return nil, err
  }

  r := JSONRpcResponseRegister{}
  err = json.Unmarshal(response, &r)
  if err != nil {
    return nil, err
  }

  return &r.Result, nil
}

/*
  Create a network
*/
func (c *Client) CreateNetwork() (*CreateNetworkResult, error) {
  params := map[string]string{}
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return nil, err
  }

  payload := createJsonPayload("createNetwork", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, MARCONI_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return nil, err
  }

  r := JSONRpcResponseCreateNetwork{}
  err = json.Unmarshal(response, &r)
  if err != nil {
    return nil, err
  }

  return &r.Result, nil
}

/*
  Delete a network
*/
func (c *Client) DeleteNetwork(networkId string) (*DeleteNetworkResult, error) {
  params := map[string]string{
    "NetworkId": networkId,
  }
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return nil, err
  }

  payload := createJsonPayload("deleteNetwork", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, MARCONI_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return nil, err
  }

  r := JSONRpcResponseDeleteNetwork{}
  err = json.Unmarshal(response, &r)
  if err != nil {
    return nil, err
  }

  return &r.Result, nil
}

/*
  Add a peer to a network
*/
func (c *Client) AddPeer(networkContractAddress string, peerPubKeyHash string, waitForReceipt bool) (interface{}, error) {
  params := map[string]interface{}{
    "NetworkContractAddress": networkContractAddress,
    "PeerPubKeyHash":         peerPubKeyHash,
    "WaitForReceipt":         waitForReceipt,
  }
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return nil, err
  }

  payload := createJsonPayload("addPeer", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, MARCONI_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return nil, err
  }

  if waitForReceipt {
    r := JSONRpcResponseAddPeer{}
    err = json.Unmarshal(response, &r)
    if err == nil {
      return r.Result, nil
    }
  } else {
    r := JSONRpcResponseTransactionHash{}
    err = json.Unmarshal(response, &r)
    if err == nil {
      return r.Result, nil
    }
  }
  return nil, err
}

/*
  Remove a peer to a network
*/
func (c *Client) RemovePeer(networkContractAddress string, peerPubKeyHash string, waitForReceipt bool) (interface{}, error) {
  params := map[string]interface{}{
    "NetworkContractAddress": networkContractAddress,
    "PeerPubKeyHash":         peerPubKeyHash,
    "WaitForReceipt":         waitForReceipt,
  }
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return nil, err
  }

  payload := createJsonPayload("removePeer", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, MARCONI_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return nil, err
  }

  if waitForReceipt {
    r := JSONRpcResponseRemovePeer{}
    err = json.Unmarshal(response, &r)
    if err == nil {
      return r.Result, nil
    }
  } else {
    r := JSONRpcResponseTransactionHash{}
    err = json.Unmarshal(response, &r)
    if err == nil {
      return r.Result, nil
    }
  }
  return nil, err
}

/*
  Add an edge between peers in a network
*/
func (c *Client) AddPeerRelation(networkContractAddress string, peerPubKeyHash string, otherPeerPubKeyHash string, waitForReceipt bool) (interface{}, error) {
  params := map[string]interface{}{
    "NetworkContractAddress": networkContractAddress,
    "PeerPubKeyHash":         peerPubKeyHash,
    "OtherPeerPubKeyHash":    otherPeerPubKeyHash,
    "WaitForReceipt":         waitForReceipt,
  }
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return nil, err
  }

  payload := createJsonPayload("addPeerRelation", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, MARCONI_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return nil, err
  }

  if waitForReceipt {
    r := JSONRpcResponseAddPeerRelation{}
    err = json.Unmarshal(response, &r)
    if err == nil {
      return r.Result, nil
    }
  } else {
    r := JSONRpcResponseTransactionHash{}
    err = json.Unmarshal(response, &r)
    if err == nil {
      return r.Result, nil
    }
  }
  return nil, err
}

/*
  Remove an edge between peers in a network
*/
func (c *Client) RemovePeerRelation(networkContractAddress string, peerPubKeyHash string, otherPeerPubKeyHash string, waitForReceipt bool) (interface{}, error) {
  params := map[string]interface{}{
    "NetworkContractAddress": networkContractAddress,
    "PeerPubKeyHash":         peerPubKeyHash,
    "OtherPeerPubKeyHash":    otherPeerPubKeyHash,
    "WaitForReceipt":         waitForReceipt,
  }
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return nil, err
  }

  payload := createJsonPayload("removePeerRelation", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, MARCONI_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return nil, err
  }

  if waitForReceipt {
    r := JSONRpcResponseRemovePeerRelation{}
    err = json.Unmarshal(response, &r)
    if err == nil {
      return r.Result, nil
    }
  } else {
    r := JSONRpcResponseTransactionHash{}
    err = json.Unmarshal(response, &r)
    if err == nil {
      return r.Result, nil
    }
  }
  return nil, err
}

/*
  Inspect a peer's edges in a network
*/
func (c *Client) GetPeerRelations(networkContractAddress string, pubKeyHash string) (string, error) {
  params := map[string]string{
    "NetworkContractAddress": networkContractAddress,
    "PubKeyHash":             pubKeyHash,
  }
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return "", err
  }

  payload := createJsonPayload("getPeerRelations", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, MARCONI_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return "", err
  }

  r := JSONRpcResponse{}
  err = json.Unmarshal(response, &r)
  if err != nil {
    return "", err
  }

  return r.Result, nil
}

/*
  Inspect a peer's info in a network
*/
func (c *Client) GetPeerInfo(networkContractAddress string, pubKeyHash string) (*PeerInfoResult, error) {
  params := map[string]string{
    "NetworkContractAddress": networkContractAddress,
    "PubKeyHash":             pubKeyHash,
  }
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return nil, err
  }

  payload := createJsonPayload("getPeerInfo", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, MARCONI_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return nil, err
  }

  r := JSONRpcResponsePeerInfo{}
  err = json.Unmarshal(response, &r)
  if err != nil {
    return nil, err
  }

  return &r.Result, nil
}

/*
  Construct and return the network related information
*/
func (c *Client) GetNetworkInfo(networkContractAddress string) (*NetworkInfoResult, error) {
  id, err := c.GetInfoFromNetwork(networkContractAddress, "getNetworkId")
  if err != nil {
    id = ""
    fmt.Println("Failed to get network id", err)
  }

  admin, err := c.GetInfoFromNetwork(networkContractAddress, "getNetworkAdmin")
  if err != nil {
    admin = ""
    fmt.Println("Failed to get network admin", err)
  }

  peers, err := c.GetInfoFromNetwork(networkContractAddress, "getPeers")
  if err != nil {
    peers = ""
    fmt.Println("Failed to get peers", err)
  }

  networkInfo := NetworkInfoResult{NetworkId: id, NetworkAdmin: admin, Peers: peers}

  return &networkInfo, nil
}

/*
  Inspect a info of a network
*/
func (c *Client) GetInfoFromNetwork(networkContractAddress string, method string) (string, error) {
  params := map[string]string{
    "NetworkContractAddress": networkContractAddress,
  }
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return "", err
  }

  payload := createJsonPayload(method, string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, MARCONI_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return "", err
  }

  r := JSONRpcResponse{}
  err = json.Unmarshal(response, &r)
  if err != nil {
    return "", err
  }

  return r.Result, nil
}

/*
  Send a transaction
*/
func (c *Client) SendTransaction(password string, nonce uint64, fromAddress string, toAddress string, amount *big.Int, gasLimit uint64, gasPrice *big.Int) (string, error) {
  keystore, err := mkey.GetAccountForAddress(fromAddress)
  if err != nil {
    return "", err
  }

  transaction := blockchain.CreateTransaction(nonce, toAddress, amount, gasLimit, gasPrice)
  signedTransaction, err := blockchain.SignTransaction(keystore, transaction, password)
  if err != nil {
    return "", err
  }

  var byteBuffer bytes.Buffer
  err = signedTransaction.EncodeRLP(&byteBuffer)
  if err != nil {
    return "", err
  }
  rlpEncodedTransaction := fmt.Sprintf("0x%x", byteBuffer.Bytes())
  params := []string{rlpEncodedTransaction}
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return "", err
  }

  payload := createJsonPayload("eth_sendRawTransaction", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, ETH_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return "", err
  }

  // parse response from http request to a certain type of response struct obj
  r := JSONRpcResponse{}
  err = json.Unmarshal(response, &r)
  if err != nil {
    return "", err
  }

  return r.Result, nil
}

/*
  Get the number of transactions made by the provided address
*/
func (c *Client) GetTransactionCount(address string) (int, error) {
  params := []string{address}
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return 0, err
  }

  payload := createJsonPayload("eth_getTransactionCount", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, ETH_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return 0, err
  }

  // parse response from http request to a certain type of response struct obj
  r := JSONRpcResponseInt{}
  err = json.Unmarshal(response, &r)
  if err != nil {
    return 0, err
  }

  return r.Result, nil
}

/*
  Get transaction receipt by transaction hash
*/
func (c *Client) GetTransactionReceipt(transactionHash string) (*Reciept, error) {
  params := []string{transactionHash}
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return nil, err
  }

  payload := createJsonPayload("eth_getTransactionReceipt", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, ETH_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return nil, err
  }

  // parse response from http request to a certain type of response struct obj
  r := JSONRpcResponseReceipt{}
  err = json.Unmarshal(response, &r)
  if err != nil {
    return nil, err
  }

  return &r.Result, nil
}

/*
  Let middleware update userAddress in user_conf.json
*/
func (c *Client) UpdateUserAddress(userAddress string) (bool, error) {
  params := []string{
    userAddress,
  }
  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return false, err
  }

  payload := createJsonPayload("updateUserAddress", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, MIDDLEWARE_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return false, err
  }

  r := JSONRpcResponseBool{}
  err = json.Unmarshal(response, &r)
  if err != nil {
    return false, err
  }

  return r.Result, nil
}

/*
  Call middleware to start netflow monitor
 */
func (c *Client) StartNetflow(collectorIp string, collectorPort string, bridgeId string, loggingDirectory string) (error) {
  params := map[string]string{
    "collectorIp" : collectorIp,
    "collectorPort" : collectorPort,
    "interface" : bridgeId,
    "loggingDirectory" : loggingDirectory,
  }

  paramsBytes, err := json.Marshal(params)
  if err != nil {
    return err
  }
  payload := createJsonPayload("startNetflow", string(paramsBytes))
  response, err := sendJsonRpcOverHttp(c.host, c.port, MARCONI_API_MIDDLEWARE_URL_PATH, payload)
  if err != nil {
    return err
  }
  r := JSONRpcResponse{}
  err = json.Unmarshal(response, &r)
  if err != nil {
    return err
  }

  return nil
}