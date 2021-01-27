#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $5)
    local CP=$(one_line_pem $6)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${ORGMSP}/$2/" \
        -e "s/\${P0PORT}/$3/" \
        -e "s/\${CAPORT}/$4/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        ../../connections/ccp-template.json 
}


ORG=citizen
ORGMSP=Citizen
P0PORT=7051
CAPORT=7054
PEERPEM=../crypto-config/peerOrganizations/citizen.vehicle.com/tlsca/tlsca.citizen.vehicle.com-cert.pem
CAPEM=../crypto-config/peerOrganizations/citizen.vehicle.com/ca/ca.citizen.vehicle.com-cert.pem

echo "$(json_ccp $ORG $ORGMSP $P0PORT $CAPORT $PEERPEM $CAPEM)" >../../connections/connection-citizen.json

ORG=lawyer
ORGMSP=Lawyer
P0PORT=8051
CAPORT=8054
PEERPEM=../crypto-config/peerOrganizations/lawyer.vehicle.com/tlsca/tlsca.lawyer.vehicle.com-cert.pem
CAPEM=../crypto-config/peerOrganizations/lawyer.vehicle.com/ca/ca.lawyer.vehicle.com-cert.pem

echo "$(json_ccp $ORG $ORGMSP $P0PORT $CAPORT $PEERPEM $CAPEM)" >../../connections/connection-lawyer.json
ORG=registryoffice
ORGMSP=RegistryOffice
P0PORT=9051
CAPORT=9054
PEERPEM=../crypto-config/peerOrganizations/registryoffice.vehicle.com/tlsca/tlsca.registryoffice.vehicle.com-cert.pem
CAPEM=../crypto-config/peerOrganizations/registryoffice.vehicle.com/ca/ca.registryoffice.vehicle.com-cert.pem

echo "$(json_ccp $ORG $ORGMSP $P0PORT $CAPORT $PEERPEM $CAPEM)" >../../connections/connection-registryoffice.json

ORG=blro
ORGMSP=BLRO
P0PORT=10051
CAPORT=10054
PEERPEM=../crypto-config/peerOrganizations/blro.vehicle.com/tlsca/tlsca.blro.vehicle.com-cert.pem
CAPEM=../crypto-config/peerOrganizations/blro.vehicle.com/ca/ca.blro.vehicle.com-cert.pem

echo "$(json_ccp $ORG $ORGMSP $P0PORT $CAPORT $PEERPEM $CAPEM)" >../../connections/connection-blro.json
