OUT        = ../../internal/proto
SWAGGER    = ../api/swagger
INCLUDEDIR = $(CURDIR)/includes
PROTOS     = $(shell find . -type f -name '*.proto')
APP_NAME   = api-v1-users

# Proto Dependencies
PROTODEPS += github.com/gogo/protobuf
PROTODEPS += github.com/grpc-ecosystem/grpc-gateway
PROTODEPS += github.com/mwitkow/go-proto-validators

all: codegen

.PHONY: codegen
codegen: $(PROTODEPS) | $(PROTOS)

# Links dependencies into includes directory
.PHONY: $(PROTODEPS)
$(PROTODEPS): %:
ifeq ($(wildcard $(INCLUDEDIR)/$*),)
	$(info linking $*)
	mkdir -p $(INCLUDEDIR)/$*
	ln -s  $(shell go list -m -f '{{ .Dir }}' $*)/* $(INCLUDEDIR)/$*
endif

# Generate protocol buffer files
.PHONY: $(PROTOS)
$(PROTOS): %.proto:
	@mkdir -p $(OUT)
	@protoc \
  	-I . \
	-I $(INCLUDEDIR) \
	-I $(INCLUDEDIR)/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I $(INCLUDEDIR)/github.com/grpc-ecosystem/grpc-gateway \
	--go_out=plugins=grpc:$(OUT) \
	--govalidators_out=gogoimport=true:$(OUT) \
	$*.proto
	$(MAKE) gen-proxy
	$(MAKE) gen-swagger

gen-proxy:
	protoc -I/usr/local/include -I. \
	-I $(INCLUDEDIR) \
	-I $(INCLUDEDIR)/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I $(INCLUDEDIR)/github.com/grpc-ecosystem/grpc-gateway \
	--grpc-gateway_out=logtostderr=true:$(OUT) \
    $(APP_NAME)/server.proto

gen-swagger:
	@mkdir -p $(OUT)
	protoc -I/usr/local/include -I. \
	-I $(INCLUDEDIR) \
	-I $(INCLUDEDIR)/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I $(INCLUDEDIR)/github.com/grpc-ecosystem/grpc-gateway \
	--swagger_out=logtostderr=true:$(SWAGGER) \
    $(APP_NAME)/server.proto
