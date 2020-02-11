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
TMHOME=chainroot2 go run ./cmd/. start --p2p.persistent_peers=00596aca3135db332f3605a2b49e20ac00ef9052@127.0.0.1:26656
TMHOME=chainroot3 go run ./cmd/. start --p2p.persistent_peers=00596aca3135db332f3605a2b49e20ac00ef9052@127.0.0.1:26656



```