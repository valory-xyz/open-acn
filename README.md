
The `libp2p_node` is an integral part of the ACN.

## ACN - Agent Communication Network

The agent communication network (ACN) provides a system for [agents](https://github.com/valory-xyz/open-aea) 
to find each other and communicate, solely based on their wallet addresses. 
It addresses the message delivery problem.

For more details check out the [docs](https://valory-xyz.github.io/open-aea/acn/).

## Development

To run all tests run:

``` bash
make test
```

To lint:

``` bash
make lint
```

For mocks generation:
check https://github.com/golang/mock

## Messaging patterns

Interaction protocol
___
ACN
___
TCP/UDP/...
___

### Messaging patterns inwards ACN:


Connection (`p2p_libp2p_client`) > Delegate Client > Relay Peer > Peer (Discouraged!)

Connection (`p2p_libp2p_client`)  > Delegate Client > Peer

Connection (`p2p_libp2p`) > Relay Peer > Peer

Connection (`p2p_libp2p`) > Peer


### Messaging patterns outwards ACN


Peer > Relay Peer > Delegate Client > Connection (`p2p_libp2p_client`) (Discouraged!)

Peer > Relay Peer > Connection (`p2p_libp2p`)

Peer > Delegate Client > Connection (`p2p_libp2p_client`)

Peer > Connection (`p2p_libp2p`)


In total 4*4 = 16 patterns (practically: 3*3 = 9 patterns)

## Guarantees

ACN should guarantee total ordering of messages for all agent pairs, independent of the type of connection and ACN messaging pattern used.

## Advanced feature (post `v1`):

Furthermore, there is the agent mobility. An agent can move between entry-points (Relay Peer/Peer/Delegate Client). The ACN must ensure that all messaging patterns maintain total ordering of messages for agent pairs during the move.

## ACN protocols

The ACN has the following protocols:

- register
- lookup
- unregister (dealt with by DHT defaults)
- DHT default protocols in libp2p
- message delivery protocol

## Dockerfile


#### Key and peerID generation
We need a private key (secp256k1), which can be generated using the
[go-libp2p-core crypto package](https://pkg.go.dev/github.com/libp2p/go-libp2p-core/crypto#GenerateSecp256k1Key)
From this a PeerID can be derived using the
[go-libp2p-core peer package](https://pkg.go.dev/github.com/libp2p/go-libp2p-core/peer)


#### Example: starting up a boostrap (genesis) node on a local network
**This key should be used for testing purposes only!**
```bash
private key: 54562eb807d2f80df8151db0a394cac72e16435a5f64275c277cae70308e8b24
public key:  037ed15dcee3a317e590cbdd28768ad8e2d29960b3e5d4eccca14bc94f83747f09
PeerID:      16Uiu2HAmA3cBbvMtLjqnkmoLBFuzmsVYnGfNLvr5Ws3Ey7JeRFAa
```

These can be passed using an environment variables file (`.boostrap_node`):
```bash
AEA_P2P_ID=54562eb807d2f80df8151db0a394cac72e16435a5f64275c277cae70308e8b24
AEA_P2P_URI=0.0.0.0:9000
AEA_P2P_URI_PUBLIC=0.0.0.0:10000
AEA_P2P_DELEGATE_URI=0.0.0.0:11000
AEA_P2P_URI_MONITORING=0.0.0.0:8080
```

then the boostrap node can be started as follows: 
```bash
docker run --network=host --env-file .boostrap_node -it valory/acn-node:v0.1.0 --config-from-env
```

The expected output should look as follows:
```bash
WARNING: Published ports are discarded when using host network mode
10:41:19.319 DBG node/aea/api.go:184 > env_file: .acn_config package=AeaApi
10:41:19.320 DBG node/aea/api.go:216 > msgin_path:  package=AeaApi
10:41:19.320 DBG node/aea/api.go:217 > msgout_path:  package=AeaApi
10:41:19.320 DBG node/aea/api.go:218 > id: 54562eb807d2f80df8151db0a394cac72e16435a5f64275c277cae70308e8b24 package=AeaApi
10:41:19.320 DBG node/aea/api.go:219 > addr:  package=AeaApi
10:41:19.320 DBG node/aea/api.go:220 > entry_peers:  package=AeaApi
10:41:19.320 DBG node/aea/api.go:221 > uri: 0.0.0.0:9000 package=AeaApi
10:41:19.320 DBG node/aea/api.go:222 > uri public: 0.0.0.0:10000 package=AeaApi
10:41:19.320 DBG node/aea/api.go:223 > uri delegate service: 0.0.0.0:11000 package=AeaApi
2022-09-28T10:41:19.320547746Z INF node/libp2p_node.go:67 > successfully initialized API to AEA!
2022-09-28T10:41:19.340368199Z INF node/dht/dhtpeer/dhtpeer.go:317 > My Peer ID is 16Uiu2HAmMC2tJMRaRTeWSESv8mArbq6jipJCD4adSBcBLsbc7cSL package=DHTPeer peerid=16Uiu2HAmMC2tJMRaRTeWSESv8mArbq6jipJCD4adSBcBLsbc7cSL
2022-09-28T10:41:19.340423964Z INF node/dht/dhtpeer/dhtpeer.go:319 > successfully created libp2p node! package=DHTPeer peerid=16Uiu2HAmMC2tJMRaRTeWSESv8mArbq6jipJCD4adSBcBLsbc7cSL
2022-09-28T10:41:19.340509775Z DBG node/dht/dhtpeer/dhtpeer.go:329 > Setting /aea-register/0.1.0 stream... package=DHTPeer peerid=16Uiu2HAmMC2tJMRaRTeWSESv8mArbq6jipJCD4adSBcBLsbc7cSL
2022-09-28T10:41:19.340532593Z INF node/dht/dhtpeer/dhtpeer.go:465 > Load records from store ./agent_records_store_16Uiu2HAmMC2tJMRaRTeWSESv8mArbq6jipJCD4adSBcBLsbc7cSL package=DHTPeer peerid=16Uiu2HAmMC2tJMRaRTeWSESv8mArbq6jipJCD4adSBcBLsbc7cSL
2022-09-28T10:41:19.340612489Z INF node/dht/dhtpeer/dhtpeer.go:375 > successfully loaded 0 agents package=DHTPeer peerid=16Uiu2HAmMC2tJMRaRTeWSESv8mArbq6jipJCD4adSBcBLsbc7cSL
2022-09-28T10:41:19.340630719Z DBG node/dht/dhtpeer/dhtpeer.go:398 > Setting /aea-address/0.1.0 stream... package=DHTPeer peerid=16Uiu2HAmMC2tJMRaRTeWSESv8mArbq6jipJCD4adSBcBLsbc7cSL
2022-09-28T10:41:19.340646646Z DBG node/dht/dhtpeer/dhtpeer.go:402 > Setting /aea/0.1.0 stream... package=DHTPeer peerid=16Uiu2HAmMC2tJMRaRTeWSESv8mArbq6jipJCD4adSBcBLsbc7cSL
2022-09-28T10:41:19.34068517Z INF node/dht/dhtpeer/dhtpeer.go:802 > DelegateService listening for new connections... package=DHTPeer peerid=16Uiu2HAmMC2tJMRaRTeWSESv8mArbq6jipJCD4adSBcBLsbc7cSL
2022-09-28T10:41:19.340715804Z INF node/dht/dhtpeer/dhtpeer.go:532 > Starting monitoring service: Prometheus at 8080 package=DHTPeer peerid=16Uiu2HAmMC2tJMRaRTeWSESv8mArbq6jipJCD4adSBcBLsbc7cSL
2022-09-28T10:41:19.34073922Z INF node/libp2p_node.go:148 > Peer ID: 16Uiu2HAmMC2tJMRaRTeWSESv8mArbq6jipJCD4adSBcBLsbc7cSL
MULTIADDRS_LIST_START
/dns4/0.0.0.0/tcp/10000/p2p/16Uiu2HAmMC2tJMRaRTeWSESv8mArbq6jipJCD4adSBcBLsbc7cSL
MULTIADDRS_LIST_END
```

#### Example: adding an entry peer on a local network

We create an environment variables file (`.entry_node`), containing the genesis node maddr as the entry point:
```bash
AEA_P2P_ID=f4261323cf7c42f4e5113fad6bd30c9ea71d0dbe4f34f72217e18c703cda4011
AEA_P2P_URI=0.0.0.0:9001
AEA_P2P_URI_PUBLIC=0.0.0.0:10001
AEA_P2P_DELEGATE_URI=0.0.0.0:11001
AEA_P2P_URI_MONITORING=0.0.0.0:8081
AEA_P2P_ENTRY_URIS=/dns4/0.0.0.0/tcp/9000/p2p/16Uiu2HAmMC2tJMRaRTeWSESv8mArbq6jipJCD4adSBcBLsbc7cSL
```

With the bootstrap node running in a different terminal, we can start the entry node:
```bash
docker run --network=host --env-file .entry_node -it valory/acn-node:v0.1.0 --config-from-env
```

The output one may expect is similar to that of starting up the bootstrap node.
Note that if you're not using a local network the correct port exposure and forwarding should be set.
