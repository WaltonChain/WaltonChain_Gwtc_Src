#!/bin/sh
if [ ! -d "./data" ]; then
	 ./bin/gwtc --datadir ./data/ init ./settings/wtc.json
fi
if [ "$1" = "--mine" ]; then
	nohup ./bin/gwtc --networkid 15 --datadir ./data/ --identity "wtc" $1 --etherbase $2 > gwtc.log 2>&1 &
else
	nohup ./bin/gwtc --networkid 15 --datadir ./data/ --identity "wtc" > gwtc.log 2>&1 &
fi
