#https://gist.github.com/denji/12b3a568f092ab951456
#generate private key
openssl genrsa -out server.key 2048
openssl ecparam -genkey -name secp384r1 -out server.key
#generate self-signed public key
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
