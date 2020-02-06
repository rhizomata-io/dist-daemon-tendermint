# dist-daemon-tendermint

#### Distributed Daemon framework on Tendermint


```
# Init
TMHOME=chainroot1 go run ./cmd/. init --chain-id=daemon-chain
TMHOME=chainroot2 go run ./cmd/. init --chain-id=daemon-chain
TMHOME=chainroot3 go run ./cmd/. init --chain-id=daemon-chain
```

```
# check node id
TMHOME=chainroot1 go run ./cmd/. show_node_id

```

``` 
# Run Nodes

TMHOME=chainroot1 go run ./cmd/. start  
TMHOME=chainroot2 go run ./cmd/. start --p2p.persistent_peers=7c352cbf28c81859b005e98f719a657fed713b25@127.0.0.1:26656
TMHOME=chainroot3 go run ./cmd/. start --p2p.persistent_peers=7c352cbf28c81859b005e98f719a657fed713b25@127.0.0.1:26656



```