.PHONY: gen-proto-go gen-proto-py tidy style

gen-proto-go:
	@echo "Generating Go gRPC code..."
	@protoc --proto_path=api/causal/v1alpha1 \
		--go_out=./api/causal/v1alpha1 --go_opt=paths=source_relative \
		--go-grpc_out=./api/causal/v1alpha1 --go-grpc_opt=paths=source_relative \
		causal.proto

gen-proto-py:
	@echo "Generating Python gRPC code..."
	@python3 -m grpc_tools.protoc --proto_path=api/causal/v1alpha1 \
		--python_out=worker \
		--grpc_python_out=worker \
		causal.proto

tidy:
	@echo "Running tidy..."
	@go mod tidy

style:
	@echo "Running formatter"
	@goimports -l -w ./
