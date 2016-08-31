#!/bin/sh

# server 后面的IP为允许服务器在不方便指定域名时,直接以IP访问.
# rsapire 生成用于对参数进行签名
help(){
	echo "example: ccs server  IP:10.1.9.34,IP:1.2.3.4"
	echo "example: ccs client  or: css client .+, or: css client info,system,plugin" 
	echo "example: ccs ca"
	echo "example: ccs rsapair"
}

if [ $# -lt 1 ] 
then
	help
	exit 0
fi

PASSWORD="2qif49lt"
EMAIL="11she_232@163.com"

if [ $1 = "server" ] 
then
	if [ $# -lt 2 ] 
	then
		help
		exit 0
	fi
	openssl genrsa -out new-server-key.pem 4096
	openssl req -subj "/CN=server" -sha256 -new -key new-server-key.pem -out new-server.csr
	echo "extendedKeyUsage =critical,serverAuth\n subjectAltName = $2" > extfile.cnf
	openssl x509 -req -days 3650 -sha256 -in new-server.csr -CA ca-cert.pem \
	-CAkey ca-key.pem -CAcreateserial -out new-server-cert.pem -extfile extfile.cnf \
	-passin pass:$PASSWORD

	echo "rm .csr"
	rm -v new-server.csr
elif [ $1 = "client" ]
then
	# 字符串可以是正则，如权限全开：.+
	ext='extendedKeyUsage =critical,clientAuth\n 1.2.3.4=ASN1:UTF8String:'
	if [ $# -gt 1 ]
	then
		ext="${ext}$2"
	else
		ext="${ext}ping info"
	fi
	echo "$ext"
	echo "${ext}" > extfile.cnf
	openssl genrsa -out new-client-key.pem 4096
	openssl req -subj '/CN=client' -new -key new-client-key.pem -out new-client.csr
	openssl x509 -req -days 3650 -sha256 -in new-client.csr -CA ca-cert.pem \
	-CAkey ca-key.pem -CAcreateserial -out new-client-cert.pem -extfile extfile.cnf \
	-passin pass:$PASSWORD
	echo "rm .csr"
	rm -v new-client.csr
elif [ $1 = "ca" ]
then
	openssl genrsa -aes256 -out new-ca-key.pem 4096
	openssl req -new -x509 -days 3650 -key new-ca-key.pem -sha256 -out new-ca-cert.pem \
	-subj "/C=CN/ST=ShangHai/L=ShangHai/O=Agent Inc/OU=Dev/CN=Agent CA Host/emailAddress=${EMAIL}" -passout pass:${PASSWORD}
elif [ $1 = "rsapair" ]
then
	openssl genrsa -out new-rsa-key.pem 1024
	openssl rsa -in new-rsa-key.pem -pubout -out new-rsa-pub.pem
else
	help
	exit 0
fi
