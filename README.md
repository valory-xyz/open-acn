
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

We need a private key (secp256k1), which can be generated using the
[go-libp2p-core crypto package](https://pkg.go.dev/github.com/libp2p/go-libp2p-core/crypto#GenerateSecp256k1Key)
From this a PeerID can be derived using the
[go-libp2p-core peer package](https://pkg.go.dev/github.com/libp2p/go-libp2p-core/peer)

**This key should be used for testing purposes only**
```bash
private key: 333cdf49fbd91aed1c7bf904aebac6cd9c2d4439daa0f1c9193839345f0ce78b
public key:  02159ff7eca8359add884fe09ec15399b354d1e890bfdad62d22279dcb17e889dc
PeerID:      16Uiu2HAkvt622iDZbW6r2TSNmBE9HifZSYePbc3WU9WexjtuX5a3
```

The private key file can be copied to the `node/` folder and mounted to the images:
```bash
cp bootstrap_peer/cosmos_private_key.txt ./node
```

The ACN node can then be deployed as follows:
```bash
docker run -it -p 11000:11000 -p 9000:9000 -v $(pwd)/node:/node valory/acn-node:v0.1.0 --key-file /node/cosmos_private_key.txt --uri 0.0.0.0:9000 --uri-external 0.0.0.0:10000 --uri-delegate 0.0.0.0:11000
```

The expected output should look as follows:
```bash
14:13:32.056 DBG node/aea/api.go:184 > env_file: .acn_config package=AeaApi
14:13:32.056 DBG node/aea/api.go:216 > msgin_path:  package=AeaApi
14:13:32.056 DBG node/aea/api.go:217 > msgout_path:  package=AeaApi
14:13:32.056 DBG node/aea/api.go:218 > id: 333cdf49fbd91aed1c7bf904aebac6cd9c2d4439daa0f1c9193839345f0ce78b package=AeaApi
14:13:32.056 DBG node/aea/api.go:219 > addr:  package=AeaApi
14:13:32.056 DBG node/aea/api.go:220 > entry_peers:  package=AeaApi
14:13:32.056 DBG node/aea/api.go:221 > uri: 0.0.0.0:9000 package=AeaApi
14:13:32.056 DBG node/aea/api.go:222 > uri public: 0.0.0.0:10000 package=AeaApi
14:13:32.056 DBG node/aea/api.go:223 > uri delegate service: 0.0.0.0:11000 package=AeaApi
2022-09-15T14:13:32.056407756Z INF node/libp2p_node.go:67 > successfully initialized API to AEA!
2022-09-15T14:13:32.07263438Z INF node/dht/dhtpeer/dhtpeer.go:317 > My Peer ID is 16Uiu2HAkvt622iDZbW6r2TSNmBE9HifZSYePbc3WU9WexjtuX5a3 package=DHTPeer peerid=16Uiu2HAkvt622iDZbW6r2TSNmBE9HifZSYePbc3WU9WexjtuX5a3
2022-09-15T14:13:32.072672411Z INF node/dht/dhtpeer/dhtpeer.go:319 > successfully created libp2p node! package=DHTPeer peerid=16Uiu2HAkvt622iDZbW6r2TSNmBE9HifZSYePbc3WU9WexjtuX5a3
2022-09-15T14:13:32.072694687Z DBG node/dht/dhtpeer/dhtpeer.go:329 > Setting /aea-register/0.1.0 stream... package=DHTPeer peerid=16Uiu2HAkvt622iDZbW6r2TSNmBE9HifZSYePbc3WU9WexjtuX5a3
2022-09-15T14:13:32.072715834Z INF node/dht/dhtpeer/dhtpeer.go:465 > Load records from store ./agent_records_store_16Uiu2HAkvt622iDZbW6r2TSNmBE9HifZSYePbc3WU9WexjtuX5a3 package=DHTPeer peerid=16Uiu2HAkvt622iDZbW6r2TSNmBE9HifZSYePbc3WU9WexjtuX5a3
2022-09-15T14:13:32.07279851Z INF node/dht/dhtpeer/dhtpeer.go:375 > successfully loaded 0 agents package=DHTPeer peerid=16Uiu2HAkvt622iDZbW6r2TSNmBE9HifZSYePbc3WU9WexjtuX5a3
2022-09-15T14:13:32.072818251Z DBG node/dht/dhtpeer/dhtpeer.go:398 > Setting /aea-address/0.1.0 stream... package=DHTPeer peerid=16Uiu2HAkvt622iDZbW6r2TSNmBE9HifZSYePbc3WU9WexjtuX5a3
2022-09-15T14:13:32.072838532Z DBG node/dht/dhtpeer/dhtpeer.go:402 > Setting /aea/0.1.0 stream... package=DHTPeer peerid=16Uiu2HAkvt622iDZbW6r2TSNmBE9HifZSYePbc3WU9WexjtuX5a3
2022-09-15T14:13:32.072882157Z INF node/dht/dhtpeer/dhtpeer.go:802 > DelegateService listening for new connections... package=DHTPeer peerid=16Uiu2HAkvt622iDZbW6r2TSNmBE9HifZSYePbc3WU9WexjtuX5a3
2022-09-15T14:13:32.07291222Z INF node/dht/dhtpeer/dhtpeer.go:532 > Starting monitoring service: FileMonitoring on /acn/acn.stats package=DHTPeer peerid=16Uiu2HAkvt622iDZbW6r2TSNmBE9HifZSYePbc3WU9WexjtuX5a3
2022-09-15T14:13:32.072933499Z INF node/libp2p_node.go:148 > Peer ID: 16Uiu2HAkvt622iDZbW6r2TSNmBE9HifZSYePbc3WU9WexjtuX5a3
MULTIADDRS_LIST_START
/dns4/0.0.0.0/tcp/10000/p2p/16Uiu2HAkvt622iDZbW6r2TSNmBE9HifZSYePbc3WU9WexjtuX5a3
MULTIADDRS_LIST_END
```
