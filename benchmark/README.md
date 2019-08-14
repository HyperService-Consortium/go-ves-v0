# Benchmark


```go
const (
    signatureSize = 65
    hashSize = 64
    chainIdSize = 8
    numberSize = 32 // uint256
)
```

build:

```bash
go build ./bench_nsb
```

help:

```bash
./bench_nsb -h
Usage of ./bench_nsb:
  -con int
    	signature content size (default 400)
  -node.dep int
    	arverage depth of merkle proof nodes (default 4)
  -node.siz int
    	average size of merkle proof nodes (default 1)
  -oc int
    	off-chain transaction size(in op-intent) (default 200)
  -ses int
    	max count of go-routine (default 1)
  -txcount int
    	arverage count of tx in each op-intent (default 1)
```

The size of Op-intents = txcount * [(hashSize * 4 + chainIdSize * 2 + numberSize * 2) + oc] = txcount * (336 + oc)

The size of Merkle Proof = hashSize * node.dep + node.siz * (node.dep - 1) = (64 + node.size) * node.dep - node.size

The Size of Certificate = con + signatureSize

