package middleware

import (
  "bytes"
  "encoding/json"
  "fmt"
  "github.com/pkg/errors"
  "io/ioutil"
  "math/big"
  "net/http"
  "strings"
)

// Construct a JSON RPC payload given a RPC method name and the params
func createJsonPayload(method string, params string) []byte {
  payloadStr := "{\"jsonrpc\":\"2.0\", \"id\":1, \"method\":\"" + method + "\", \"params\":" + params + " }"
  return []byte(payloadStr)
}

// Send a HTTP POST request containing a JSON RPC
func sendJsonRpcOverHttp(host string, port string, path string, jsonPayloadBytes []byte) ([]byte, error) {
  url := host + ":" + port + path
  request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayloadBytes))
  request.Header.Set("Content-Type", "application/json")

  client := new(http.Client)
  response, err := client.Do(request)
  if err != nil {
    return []byte{}, errors.New(fmt.Sprintf("Could not connect to the middleware at %s:%s", host, port))
  }
  defer response.Body.Close()

  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    return []byte{}, errors.New("Could not read the response body")
  }

  rpcError := JSONRpcResponseError{}
  err = json.Unmarshal(body, &rpcError)
  if err != nil {
    return nil, err
  }
  if rpcError.Error != nil {
    errMsg := parseRpcError(rpcError.Error)
    return nil, errors.New(errMsg)
  }

  return body, nil
}

// Convert a string hexadecimal number to a string decimal number
func hexStringToDecimalString(hexString string) string {
  if val, parsed := new(big.Int).SetString(hexString, 0); parsed {
    return val.String()
  }
  return ""
}

// Given a RpcError object, return a human readable string message
func parseRpcError(rpcError *RpcError) string {
  switch rpcError.Code {
    case RpcErrorCode_ParseError:
      return fmt.Sprintf("Parse Error: %s", rpcError.Message)
    case RpcErrorCode_InvalidRequest:
      return fmt.Sprintf("Invalid Request: %s", rpcError.Message)
    case RpcErrorCode_MethodNotFound:
      return fmt.Sprintf("Method Not Found: %s", rpcError.Message)
    case RpcErrorCode_InvalidParams:
      return fmt.Sprintf("Invalid Params: %s", rpcError.Message)
    case RpcErrorCode_InternalError:
      // special case for authentication error
      const AuthenticationErrorMsg = "authentication"
      if strings.Contains(rpcError.Message, AuthenticationErrorMsg) {
        return "Please unlock your account first using credential mode"
      }
      return fmt.Sprintf("Server Internal Error: %s", rpcError.Message)
    default:
      return fmt.Sprintf("Code: %d, Message: %s", rpcError.Code, rpcError.Message)
  }
}
