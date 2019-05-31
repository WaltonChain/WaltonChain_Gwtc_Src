#!/bin/sh
if [ ! -d "./data/gwtc" ]; then
	./bin/gwtc --datadir ./data/ init ./settings/wtc.json
fi
if [ "$1" = "--mine" ]; then
./bin/gwtc --networkid 15 --datadir ./data/ --identity "wtc" $1 --etherbase $2 console
else
./bin/gwtc --networkid 15 --datadir ./data/ --identity "wtc" console
fi
