#!/bin/bash

#source /home/rogelio/.wasmedge/env
# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
# ancond start --log_level $LOGLEVEL --trace
#export GIN_MODE=release

./main --ipfshost=https://ipfs.xdv.digital/  --rootkey=baguqeeramdqve3hvfjccnvisofsg74pr7vh2iko5fzdqefbo2moa3hza73vq --peeraddr /ip4/67.205.135.10/tcp/4001/p2p/12D3KooWAFuNRzKGnc5hTCeYDrWuaUdjG9Py77fUCe3FJ3xUQ26c --cors true  --origins=http://172.23.247.211/:3000,http//localhost:3000 
