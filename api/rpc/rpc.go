package rpc

import (
  "bytes"
  "github.com/MarconiProtocol/cli/core/configs"
  "github.com/gorilla/rpc/v2/json2"
  "net/http"
  "sync"
)

const (
  RPC_PATH = "/rpc/m/request"
)

type Client struct {
  URL string
}

var once sync.Once
var client *Client

func GetRPCClient() *Client {
  once.Do(func() {
    conf := configs.LoadBaseConf()
    client = &Client{
      URL: conf.MarconiNodeHost + ":" + conf.MarconidRPCPort + RPC_PATH,
    }
  })
  return client
}

func (c *Client) Call(method string, args interface{}, reply interface{}) error {
  message, err := json2.EncodeClientRequest(method, args)
  if err != nil {
    return err
  }

  resp, err := http.Post(c.URL, "application/json", bytes.NewReader(message))
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  return json2.DecodeClientResponse(resp.Body, &reply)
}
