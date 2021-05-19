# btc-neo4j

This program allows to parse the bitcoin blockchain and write nodes and relationships according to the following schema to a neo4j graph.

## Getting Started

The following section describes how to setup the environment to be able to run btc-neo4j.

### Prerequisits

1. Install Neo4j Desktop (download it [here](https://neo4j.com/download-center/))

2. Install Golang and set up your environment.

3. Install Seabolt (check [here](https://github.com/neo4j-drivers/seabolt)). There is a description for Ubuntu. For Arch/Manjaro you need to install CMake >= 3.12 and OpenSSL Development Libraries (must include static libraries).

4. Install neo4j-go-driver (instructions can be found [here](https://github.com/neo4j/neo4j-go-driver)

5. Get this program:
```
go get -u github.com/wallerprogramm/btc-neo4j
```
6. Create a file with name .env in `$GOPATH/src/github.com/wallerprogramm/btc-neo4j` with the folloing content:
```
# The package
PKG=btc-neo4j
# The folder where the BTC blockchain is stored
DATADIR=/path/to/bitcoinfolder
#often: DATADIR=/home/username/.bitcoin
```

## Run btc-neo4j

### Prerequisits

Install all dependecies of btc-neo4j:
```
$cd $GOPATH/github.com/wallerprogramm/btc-neo4j
$go get ./...
```
Install btc-neo4j:
```
$cd $GOPATH/github.com/wallerprogramm/btc-neo4j
$go install
```

### Run it

1. Start Neo4j Desktop
2. Create a new graph (choose a name, pw: blob). If you have a that contains already a part of the data, you can use this one. Comment the following out in graph Settings:
```
dbms.memory.heap.initial_size=512m
dbms.memory.heap.max_size=1G
dbms.memory.pagecache.size=512m
```
3. Start the graph.
4. Run btc-neo4j:
```
$ $GOBIN/btc-neo4j
```
