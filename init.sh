go mod tidy
go build main.go
./main --rootkey bafyreibqiybccahqa3vhpxluhglojfh654x3bfgrbdy73pgywaxam5sdim --roothash 7f00cfd66db986fdb5ab62e93c0bb592e02392b51ba0d1f1d1a64e2e765b1f5c

# Remote peer connection command
# ./main --peers /ip4/190.140.158.159/tcp/7702/p2p/12D3KooWSEdJA3WovMnEZdh2oNgBsGJX85rCjzgcjczdJ92sg2WW --rootkey bafyreibqiybccahqa3vhpxluhglojfh654x3bfgrbdy73pgywaxam5sdim --roothash 7f00cfd66db986fdb5ab62e93c0bb592e02392b51ba0d1f1d1a64e2e765b1f5c --sync true