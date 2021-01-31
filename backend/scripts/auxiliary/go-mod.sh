cd ../../chaincode

CC_DIR=$PWD

CC_NAMES="lawyer_cc registryoffice_cc blro_cc land_cc transfer_cc"

for CC in $CC_NAMES; do
    echo "Installing Go dependencies in "$CC
    cd $CC
    go mod vendor
    cd ..
done
echo "Installing Go dependencies complete!"
