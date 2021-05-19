# btc-csv

This program allows to parse the bitcoin blockchain and write the nodes and relationships of the created Tranaction Graph to csv.

## Getting Started

The following section describes how to setup the environment to be able to run btc-csv.

### Prerequisits

1. Install Golang and set up your environment.

2. Get this program:
```
$ go get -u github.com/wallerprogramm/btc-csv
```
6. Create a file with name .env in `$GOPATH/src/github.com/wallerprogramm/btc-csv` with the folloing content:
```
# The package
PKG=btc-csv
# The folder where the BTC blockchain is stored
DATADIR=/path/to/bitcoinfolder
#often: DATADIR=/home/username/.bitcoin
# The folder where the utxo leveldb is stored
DBDIR=/path/to/leveldb/utxo
```

## Run btc-csv

### Prerequisits

Choose from which height to which height should be parsed. To do so enter main.go and change height to the wished end block height. If the starting block height should be changed to n, be aware that the `utxo` leveldb with block height (n-1) has to be provided. To change the start height change the 0 in `btc.LoadFile(0, height, grapho.Build, c)` to the wished start height.

Install all dependecies of btc-csv:
```
$ cd $GOPATH/github.com/wallerprogramm/btc-csv
$ go get ./...
```
Install btc-csv:
```
$ cd $GOPATH/github.com/wallerprogramm/btc-csv
$ go install
```

### Run it

1. Make sure that all the csv files that should not be changed are moved away from `$GOPATH/github.com/wallerprogramm/btc-csv`.
2. Make sure that the leveldb is either inexistent (if start height = 0) or the needed leveldb is in place (if start height > 0)
3. Run btc-csv:
```
$ $GOBIN/btc-csv
```
4. Wait for a couple of hours.

## Remove duplicate addresses from addresses.csv

```
$ cd $GOPATH/github.com/wallerprogramm/btc-csv
$ sort -u addresses.csv -o addresses.csv -T /path/to/folder/for/tmp/files
```

## Import the nodes and relationships into Neo4j

### Prerequisits

1. Install Neo4j Desktop (download it [here](https://neo4j.com/download-center/)) if not already installed.

### Run it

1. Start Neo4j Desktop
2. Create a new graph (choose a name, pw: password).
3. Change the memory settings of the graph under Settings to the following:
```
dbms.memory.heap.initial_size=6G
dbms.memory.heap.max_size=6G
dbms.memory.pagecache.size=6G
```
4. Run the following command in the Neo4j Desktop Terminal:
```
$ export DATA=/path/to/folder/containing/csv-files/
$ export HEADERS=/path/to/folder/containing/csv-headers/

$ ./bin/neo4j-admin import \
	--mode=csv \
	--database=btc.db \
	--nodes:Address $HEADERS/addresses-header.csv,$DATA/addresses.csv \
	--nodes:Block $HEADERS/blocks-header.csv,$DATA/blocks.csv \
	--nodes:Transaction $HEADERS/transactions-header.csv,$DATA/transactions.csv \
	--relationships:IS_BEFORE $HEADERS/before_rel-header.csv,$DATA/before_rel.csv \
	--relationships:BELONGS_TO $HEADERS/belongs_to_rel-header.csv,$DATA/belongs_to_rel.csv \
	--relationships:RECEIVES $HEADERS/receives_rel-header.csv,$DATA/receives_rel.csv \
	--relationships:SENDS $HEADERS/sends_rel-header.csv,$DATA/sends_rel.csv \
	--ignore-missing-nodes=true \
	--ignore-duplicate-nodes=true \
	--multiline-fields=true \
	--high-io=true
```
5. Change the used database in graph Settings by adding:
```
dbms.active_database=btc.db
```
