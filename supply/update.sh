docker exec -it cli peer chaincode install -n supply -v 1.1 -l "golang" -p "github.com/supply/go"
sleep 10
docker exec -it cli peer chaincode upgrade -n supply -v 1.1 -o orderer.example.com:7050 -C mychannel -l "golang" -p "github.com/supply/go" -c '{"Args":[""]}' -P "OR ('Org1MSP.member','Org2MSP.member')"
sleep 10