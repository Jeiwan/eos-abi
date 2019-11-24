## eos-abi

ABI encoder and decoder for EOS. Uses the official ABI library [abieos](https://github.com/EOSIO/abieos).

### Usage
1. Get latest `abieos` version:
    ```
    rm -rf abieos
    make update-abieos
    ```
1. Edit `abieos/CMakeLists.txt` and replace:
    ```
    add_library(abieos MODULE src/abieos.cpp)
    ```
    With:
    ```
    add_library(abieos SHARED src/abieos.cpp)
    ```
1. Compile `abieos`:
    ```
    make compile-abieos
    ```