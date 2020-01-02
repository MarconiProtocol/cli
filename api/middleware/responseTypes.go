package middleware

/*
  Various RPC Response types
*/

// JSON RPC Error Codes
const (
  RpcErrorCode_ParseError     = -32700
  RpcErrorCode_InvalidRequest = -32600
  RpcErrorCode_MethodNotFound = -32601
  RpcErrorCode_InvalidParams  = -32602
  RpcErrorCode_InternalError  = -32603
)

type JSONRpcResponse struct {
  Jsonrpc string `json:"jsonrpc"`
  Id      int    `json:"id"`
  Result  string `json:"result"`
}

type JSONRpcResponseError struct {
  Jsonrpc string    `json:"jsonrpc"`
  Id      int       `json:"id"`
  Error   *RpcError `json:"error"`
}

type RpcError struct {
  Code    int    `json:"code"`
  Message string `json:"message"`
}

type JSONRpcResponseInt struct {
  Jsonrpc string `json:"jsonrpc"`
  Id      int    `json:"id"`
  Result  int    `json:"result"`
}

type JSONRpcResponseBool struct {
  Jsonrpc string `json:"jsonrpc"`
  Id      int    `json:"id"`
  Result  bool   `json:"result"`
}

type JSONRpcResponseReceipt struct {
  Jsonrpc string  `json:"jsonrpc"`
  Id      int     `json:"id"`
  Result  Reciept `json:"result"`
}

type Reciept struct {
  BlockHash         string
  BlockNumber       string
  ContractAddress   string
  CumulativeGasUsed string
  From              string
  GasUsed           string
  Logs              []string
  LogsBloom         string
  To                string
  TransactionHash   string
  TransactionIndex  string
}

type JSONRpcResponseCreateNetwork struct {
  Jsonrpc string              `json:"jsonrpc"`
  Id      int                 `json:"id"`
  Result  CreateNetworkResult `json:"result"`
}

type JSONRpcResponseRegister struct {
  Jsonrpc string         `json:"jsonrpc"`
  Id      int            `json:"id"`
  Result  RegisterResult `json:"result"`
}

type RegisterResult struct {
  PubKeyHash string
}

type CreateNetworkResult struct {
  NetworkId       string
  NetworkContract string
  Admin           string
}

type JSONRpcResponseDeleteNetwork struct {
  Jsonrpc string              `json:"jsonrpc"`
  Id      int                 `json:"id"`
  Result  DeleteNetworkResult `json:"result"`
}

type DeleteNetworkResult struct {
  NetworkId string
  Admin     string
}

type JSONRpcResponseAddPeer struct {
  Jsonrpc string        `json:"jsonrpc"`
  Id      int           `json:"id"`
  Result  AddPeerResult `json:"result"`
}

type JSONRpcResponseRemovePeer struct {
  Jsonrpc string           `json:"jsonrpc"`
  Id      int              `json:"id"`
  Result  RemovePeerResult `json:"result"`
}

type AddPeerResult struct {
  NetworkId  string
  PubKeyHash string
}

type RemovePeerResult struct {
  NetworkId  string
  PubKeyHash string
}

type JSONRpcResponseAddPeerRelation struct {
  Jsonrpc string                `json:"jsonrpc"`
  Id      int                   `json:"id"`
  Result  AddPeerRelationResult `json:"result"`
}

type JSONRpcResponseRemovePeerRelation struct {
  Jsonrpc string                   `json:"jsonrpc"`
  Id      int                      `json:"id"`
  Result  RemovePeerRelationResult `json:"result"`
}

type AddPeerRelationResult struct {
  NetworkId       string
  PubKeyHashMine  string
  PubKeyHashOther string
}

type RemovePeerRelationResult struct {
  NetworkId       string
  PubKeyHashMine  string
  PubKeyHashOther string
}

type JSONRpcResponseTransactionHash struct {
  Jsonrpc string                `json:"jsonrpc"`
  Id      int                   `json:"id"`
  Result  TransactionHashResult `json:"result"`
}

type TransactionHashResult struct {
  TransactionHash string
}

type JSONRpcResponsePeerInfo struct {
  Jsonrpc string         `json:"jsonrpc"`
  Id      int            `json:"id"`
  Result  PeerInfoResult `json:"result"`
}

type JSONRpcResponsePeers struct {
  Jsonrpc string `json:"jsonrpc"`
  Id      int    `json:"id"`
  Result  string `json:"result"`
}

type PeerInfoResult struct {
  NetworkId  string `json:"0"`
  PubKeyHash string `json:"1"`
  Peers      string `json:"2"`
  IP         string `json:"3"`
  Active     bool   `json:"4"`
}

type NetworkInfoResult struct {
  NetworkId    string
  NetworkAdmin string
  Peers        string
}
