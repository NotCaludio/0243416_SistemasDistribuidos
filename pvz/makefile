# Esto nos ayuda a posicionar nuestros config files en una carpeta dentro de nuestro proyecto

CONFIG_PATH=${HOME}/sync/Documents/septimo_semestre/computoDistribuido/0243416_SistemasDistribuidos/pvz/.pvz/

.PHONY: init

init:
	mkdir -p ${CONFIG_PATH}

.PHONY: gencert
# gencert
# First creates the bare certificate, it is just the base certificate that others will differ from
# Then creates the server certificate, this allows our server certification
# Finally we create the client certificate this allows two way authentication
gencert:
	cfssl gencert \
		-initca certificates/ca-csr.json | cfssljson -bare ca
	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=certificates/ca-config.json \
		-profile=server \
		certificates/server-csr.json | cfssljson -bare server
	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=certificates/ca-config.json \
		-profile=client \
		certificates/client-csr.json | cfssljson -bare client
	mv *.pem *.csr ${CONFIG_PATH}

compile:
	protoc api/v1/*.proto \
					--go_out=.\
					--go_opt=paths=source_relative \
					--proto_path=.
test:
	go test -race ./...
compile_rpc:
	protoc api/v1/*.proto \
	--go_out=. \
	--go_opt=paths=source_relative \
    --go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \