# Universal WebAssembly Virtual Machine

UWAVM is a WebAssembly VM in Go. 
UWAVM provide tools to decode wasm binary files that compiled by golang, javascript, c/c++, rust.

## Quick start

### Requirements
* Go 1.12.x or later

### Build
1. Clone the repository
```
git clone https://github.com/BeDreamCoder/uwavm.git
```
2. Compile wasm contract
```
make wasm
```
3. Build uwavm executable program
```
make build
```

### Run
```
cd output
```

#### Deploy contract
```
./uwavm contract deploy -n erc20 -l go -a '{"initSupply":"1000000"}' -p ./wasm/erc20.wasm -c alice
```

#### Query contract
```
./uwavm contract query -n erc20 -l go -m query -a '{"action":"balanceOf","address":"alice"}' -c alice
```

#### Invoke contract
```
./uwavm contract invoke -n erc20 -l go -m invoke -a '{"action":"transfer","to":"bob","amount":"100"}' -c alice
```
