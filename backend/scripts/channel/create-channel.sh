#!/bin/bash
echo "Creating channel..."
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/lran.com/orderers/orderer.lran.com/msp/tlscacerts/tlsca.lran.com-cert.pem
CORE_PEER_LOCALMSPID=CitizenMSP
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/citizen.lran.com/peers/peer0.citizen.lran.com/tls/ca.crt
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/citizen.lran.com/users/Admin@citizen.lran.com/msp
CORE_PEER_ADDRESS=peer0.citizen.lran.com:7051
CHANNEL_NAME=mainchannel
CORE_PEER_TLS_ENABLED=true
ORDERER_SYSCHAN_ID=syschain

peer channel create -o orderer.lran.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt

cat log.txt

#peer channel create -o orderer.lran.com:7050 /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/lran.com/orderers/orderer.lran.com/msp/tlscacerts/tlsca.lran.com-cert.pem -c mainchannel -f ./channel-artifacts/channel.tx
sleep 10
echo
echo "Channel created, joining Citizen..."
peer channel join -b mainchannel.block
