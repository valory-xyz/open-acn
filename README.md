
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

We need a private key, we can be created using open-aea: 

```fish
aea create bootstrap_peer --local

cd bootstrap_peer && \
        aea generate-key cosmos && \
        aea add-key cosmos && \
        aea get-multiaddress cosmos && \
        aea get-public-key cosmos && \
        cd ../
```

And the expected output should look as follows:
```
Initializing AEA project 'bootstrap_peer'
Creating project directory './bootstrap_peer'
Creating config file aea-config.yaml
Adding default packages ...
Adding protocol 'open_aea/signing:1.0.0:bafybeiambqptflge33eemdhis2whik67hjplfnqwieoa6wblzlaf7vuo44'...
Successfully added protocol 'open_aea/signing:1.0.0'.
16Uiu2HAmA3cBbvMtLjqnkmoLBFuzmsVYnGfNLvr5Ws3Ey7JeRFAa
02d9383bc6e37e4d6b56cecb93c1e2796407c07decaa3ae598b6f762946a464353
```

To establish a delegate connection the public key is required, which is displayed on the last line.
The private key file can be copied to the `node/` folder and mounted to the images:
```fish
cp bootstrap_peer/cosmos_private_key.txt ./node
```

The ACN node can then be deployed as follows:
```fish
docker run -it -p 11000:11000 -p 9000:9000 -v (pwd)/node:/node valory/acn-node:v0.1.0 --key-file /node/cosmos_private_key.txt --uri 0.0.0.0:9000 --uri-external 0.0.0.0:10000 --uri-delegate 0.0.0.0:11000
```


