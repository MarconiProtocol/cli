# Overview
The Marconi CLI Client (mCLI) is the main entry point to interacting with the Marconi software components. It is a user friendly command line tool that provides the functionality of a wallet as well as interaction with the Marconi Net contracts.

## Quick Links
- [Using mCLI](#usage)
- [Design](#design)

## Usage
Marconi CLI (mCLI) client can be started with the console, used to help to start background processes or to download packages.

To start mCLI in console mode, it needs to be run with the console runtime flag.
```
$ ./mcli -console
>
```  
To start any meaningful interaction with the functionality of the different installed modes the client must first be switched to a mode. This can be done with the `jump` command.
```
> j credential
credential>
```  
At anytime, `tab` can be pressed to display a suggestions list.  
A more extensive example of using mCLI to create an account, create a Marconi subnet and more can be found in our [wiki](https://github.com/MarconiProtocol/wiki/wiki/Setup-Instructions)

### Execution Mode

mCLI can also be run with the execution flag. This allows you to run any commands without entering into console mode.
Commands can be run in the form: <mode> <command>

Current Modes Supported:
- credential
- net

(process related commands are handled by daemon mode)

```
$ ./mcli --mode=exec --command="<mode><command>"
$ ./mcli --mode=exec --command="credential account list"
$ ./mcli --mode=exec --command="net create"
```
Multiple commands can be run at once, separate each command with a ;
```
$ ./mcli --mode=exec --command="credential account create --password test; credential account list"
```

## Modes
At this moment mCLI is released with the following modes:
- [marconi_credential](#credential)
- [marconi_net](#net)
- [marconi_process](#process)
***

### credential
This mode helps users create and manage the different aspects of a Marconi Account.
```
credential>
             account  Account related commands  
             key      Key related commands      
             home     Return to home menu       
             exit     Exit mcli                 
```


#### account
Account is a `credential` submode, with the following commands
```
credential> account
                     create   Create account                                         
                     unlock   Unlock account                                         
                     list     List accounts                                          
                     send     Send a transaction                                     
                     balance  Get balance for an account                             
                     receipt  Get receipt for a transaction                          
                     export   Export GO Marconi Keystore associate with an account
```

##### account create
Creates a new Marconi account.
```
credential> account create [Optional: --password <PASSWORD> | --password-file <PASSWORD_FILE>]
```
Optional
- `--password <PASSWORD>`  Password can optionally be provided on the command line (if not, the user will be prompted)
- `--password-file <PASSWORD_FILE>` Path to a file containing the password that can be optionally provided. (if not, the user will be prompted)

##### account unlock
Unlocks a given Marconi account.  
```
credential> account unlock <0xACCOUNT_ADDRESS> [Optional: --password <PASSWORD> | --password-file <PASSWORD_FILE>]
```
 - `<0xACCOUNT_ADDRESS>`   The Marconi address of the account to unlock.  
 Optional:
- `--password <PASSWORD>`  Password can optionally be provided on the command line (if not, the user will be prompted)
- `--password-file <PASSWORD_FILE>` Path to a file containing the password that can be optionally provided. (if not, the user will be prompted)

##### account list
Lists all available accounts to use from mCLI.
```
credential> account list
```   

##### account send
Sends Marcos from your account to a target account.
```
credential> account send <0xACCOUNT_ADDRESS> <0xTARGET_ADDRESS> <AMOUNT_IN_MARCOS> <GAS_LIMIT> <GAS_PRICE> [Optional: --password <PASSWORD> | --password-file <PASSWORD_FILE> | --skip-prompts]
```  
 - `<0xACCOUNT_ADDRESS>`   Your Marconi address to send Marcos from.  
 - `<0xTARGET_ADDRESS>`    The target Marconi address to send Marcos to.  
 - `<AMOUNT_IN_MARCOS>`    The amount of Marcos to send.  
 - `<GAS_LIMIT>`           The upper limit in gas to spend on this value transfer tx.  
 - `<GAS_PRICE>`           The price per unit of gas in Gauss you are willing to pay.  
 
 Optional:
- `--password <PASSWORD>`  Password can optionally be provided on the command line (if not, the user will be prompted)
- `--password-file <PASSWORD_FILE>` Path to a file containing the password that can be optionally provided. (if not, the user will be prompted)
- `--skip-prompts`          Indicates that prompts should be skipped (ie. Confirmations) and defaults used instead

##### account balance
Checks the Marcos balance of a given account.
```
credential> account balance <0xACCOUNT_ADDRESS>
```  
 - `<0xACCOUNT_ADDRESS>`   The Marconi address whose balance will be returned.

##### account receipt
Returns the receipt of a transaction by transaction hash.
```
credential> account receipt <0xTRANSACTION_HASH>
```  
 - `<0xTRANSACTION_HASH>`   The transaction hash of the transaction whose receipt will be returned.  

##### account export
Exports the Marconi keystore file stored in the account file
```
credential> account export <0xACCOUNT_ADDRESS> <GO-MARCONI_DATA_DIR_PATH>
```  
 - `<0xACCOUNT_ADDRESS>`        The Marconi address of the account to be exported.  
 - `<GO-MARCONI_DATA_DIR_PATH>` The directory where the Marconi account file is stored.  

#### credential> key
Key is a `credential` submode, with the following commands  
```
credential> key
                 generate  Generate nodekey                        
                 use       Set nodekey to use with other commands  
                 export    Export nodekey                          
                 list      List nodekeys                           
```

##### key generate
Generates a new Marconi node key.
```
credential> key generate <0xACCOUNT_ADDRESS> [Optional: --password <PASSWORD> | --password-file <PASSWORD_FILE>]
```
- `<0xACCOUNT_ADDRESS>`   The Marconi address to generate a new node key for.

Optional:
 - `--password <PASSWORD>`  Password can optionally be provided on the command line (if not, the user will be prompted)
 - `--password-file <PASSWORD_FILE>` Path to a file containing the password that can be optionally provided. (if not, the user will be prompted)

##### key use 
Set the Marconi node key to use with other commands.
```
credential> key use <0xACCOUNT_ADDRESS> [Optional: --password <PASSWORD> | --password-file <PASSWORD_FILE> | --node-key <NODE_KEY> (Default 0) | --skip-prompts]
```
- `<0xACCOUNT_ADDRESS>`   The Marconi address to acquire node key from

Optional:
 - `--password <PASSWORD>`  Password can optionally be provided on the command line (if not, the user will be prompted)
 - `--node-key <NODE_KEY>`  Set the node key to be used (other wise user will be prompted for it)
 - `--skip-prompts` Optional flag to indicate if prompts should be skipped and defaults used instead (if --node-key is not present 0 will be set for this value)
 
##### key export
Exports a node key from the account file to public/private key files.
```
credential> key export <0xACCOUNT_ADDRESS> [Optional: --password <PASSWORD> | --password-file <PASSWORD_FILE>]
```
- `<0xACCOUNT_ADDRESS>`   The Marconi address whose node keys to export.

Optional:
 - `--password <PASSWORD>`  Password can optionally be provided on the command line (if not, the user will be prompted)
 - `--password-file <PASSWORD_FILE>` Path to a file containing the password that can be optionally provided. (if not, the user will be prompted)



### net
This mode helps users manage a Marconi subnet.  
**NOTE: net use command must first be used to set the Marconi subnet on which the other commands will operate on**
```
net>
      peer    Peer related commands                   
      util    Utility commands                        
      use     Set network to use with other commands  
      create  Create new network                      
      delete  Delete existing network                 
      join    Join an existing network                
      info    Get network info
      home    Return to home menu
      exit    Exit mcli                 

```

#### use
Set the Marconi subnet to the one provided. All Marconi Net commands will be operated on the subnet set by this command.
```
net> use <0xNETWORK_CONTRACT_ADDRESS>
```
- `<0xNETWORK_CONTRACT_ADDRESS>`   The address of the network contract to be interacted with.

#### peer
Peer is a `net` submode, with the following commands
```
net> peer
           add              Add peer to a network       
           remove           Remove peer from a network  
           add_relation     Add peer relation           
           remove_relation  Remove peer relationship    
           relations        Get node relationships      
           info             Get node info               

```

##### peer add
Adds a peer to the Marconi subnet
```
net> peer add <PEER_NODE_ID>
```
- `<PEER_NODE_ID>`  The node id of the peer that will be added to the Marconi subnet.

##### peer remove
Removes a peer from the Marconi subnet
```
net> peer remove <PEER_NODE_ID>
```
- `<PEER_NODE_ID>`  The node id of the peer that will be removed from the Marconi subnet.

##### peer add_relation
Adds a relationship between two peers in the Marconi subnet. This relationship dictates that a mPipe will be created between the two nodes.
```
net> peer add_relation <PEER_NODE_ID> <OTHER_PEER_NODE_ID>
```
- `<PEER_NODE_ID>`        The node id of one peer in the relationship.
- `<OTHER_PEER_NODE_ID>`  The node id of the the other peer in the relationship.

##### peer remove_relation
Removes an existing relationship between two peers in the Marconi subnet. This will cause the mPipe between the two nodes to be destroyed.
```
net> peer remove_relation <PEER_NODE_ID> <OTHER_PEER_NODE_ID>
```
- `<PEER_NODE_ID>`        The node id of one peer in the relationship.
- `<OTHER_PEER_NODE_ID>`  The node id of the the other peer in the relationship.

##### peer relations
Prints all relationships the given peer is in.
```
net> peer relations <PEER_NODE_ID>
```
- `<PEER_NODE_ID>`        The node id of the peer whose relationships are to be printed.

##### peer info
Prints information about this peer. Ex: The peer's relationships, assigned IP address...
```
net> peer info <PEER_NODE_ID>
```
- `<PEER_NODE_ID>`        The node id of the peer whose details are to be printed.


#### util
Util is a `net` submode, with the following commands
```
net> util
           generate_32bitkey  Generate a 32 bit key  
           register           Register a nodeID      
```

##### util generate_32bitkey
Generates a 32 bit key. Can be used as a subnet identifer
```
net> util generate_32bitkey [Optional: --path <PATH> | --skip-prompts]
```

Optional:
- `--path <PATH>` The path to save the key to, if not present (and --skip-prompts not present) the user will be prompted
- `--skip-prompts`  Indicates that all prompts should be skipped and defaults used if flags aren't set

#### create
Deploys a new network contract, effectively creating a new Marconi subnet. The address invoking this function call will become the admin of this new network.
```
net> create
```

#### delete
Removes the reference of this network contract from the network manager contract.
```
net> delete <0xNETWORK_CONTRACT_ADDRESS>
```
- `<0xNETWORK_CONTRACT_ADDRESS>`   The address of the network contract to be removed.

#### join
Join updates the configuration to join the specified network. This will not work if the node is not actually a part of the network.
```
net> join <0xNETWORK_CONTRACT_ADDRESS>
```
- `<0xNETWORK_CONTRACT_ADDRESS>`   The address of the network contract that this node is a part of.

#### info
Prints information about a specific network.
```
net> info <0xNETWORK_CONTRACT_ADDRESS>
```
- `<0xNETWORK_CONTRACT_ADDRESS>`   The address of the network contract to print information for.


### process
Used to start processes as background daemons.
```
process>
      start    Start a process and manage its lifetime                 
      stop     Stop a managed process  
      restart  Restart a managed process                      
      list     List running processes
      version  Show version of each component                     
      home     Return to home menu
      exit     Exit mcli
```

#### start
Starts the process specified
```
process> start <process_name>
```
- `<process_name> is one of gmeth, middleware or marconid`

#### stop
Stops the process specified
```
process> stop <process_name>
```
- `<process_name> is one of gmeth, middleware or marconid`

#### restart
Restarts the process specified
```
process> restart <process_name>
```
- `<process_name> is one of gmeth, middleware or marconid`

#### version
Shows the version of each component
```
process> version
```

#### list
Lists running processes and their status (ie. STOPPED/RUNNING)
```
process> list
```

## Design
The mCLI is comprised of the following components:
- [REPL Console](#repl-console)
- [Package Manager](#package-manager)
- [Process Manager](#process-manager)
- [Middleware Client](#middleware-client)

### REPL Console
The capabilities of the REPL Console (Read-Evaluate-Print Loop) is largely dependent on the plugins the users choose to install with mCLI.   
The goal is to have a multipurpose CLI client whose functionality can be expanded on easily.  
Integration with plugins is not yet complete but will be released soon. However the current iteration of the codebase has been designed with this goal kept in mind.

### Package Manager
The package manager helps to download extra dependencies or packages that mCLI can run. The downloaded packages are configured through `packages_conf.json`.  
Here is a sample snippet:
```
{
  "Version": "0.1.878",
  "Packages" : [
    {
      "Id": "marconid",
      "Dir": "./",
      "Source": "https://download.marconi.org/deployment/components/marconid/0.1.1063/marconid_linux.tar.gz",
      "Version": "0.1.1063",
      "VersionFile": "etc/marconid/version.txt"
    },

    ...

  ]
}

```
In this snippet, the package manager is configured to download the `marconid` package and will extract it to the directory `./`

### Process Manager
The process manager runs processes
Here is a sample snippet:
```
{
  "Version": "0.0.1",
  "Processes" : [   
    {
      "Id": "marconid",
      "Dependencies": ["meth", "middleware"],
      "Dir": "./bin",
      "Source": "",
      "Version": "",
      "Command": "./marconid",
      "Arguments": ["/opt/marconi/etc/marconid/l2.key", "/opt/marconi/etc/marconid/block/basebeacon_cluster1"],
      "LogFilename": "marconid.log",
      "WaitForCompletion": false,
      "PidFilename": "marconid.pid"
    }
  ]
}

```
In this snippet, the process manager is configured to run `./marconid` in the `./bin` directory with the arguments `/opt/marconi/etc/marconid/l2.key`, `/opt/marconi/etc/marconid/block/basebeacon_cluster1` and it's output to stdout is logged to `marconid.log`.

### Middleware Client
The middleware client is used to interface with the Marconi middleware. The middleware client sends JSON RPC over http to the locally running middleware process.
