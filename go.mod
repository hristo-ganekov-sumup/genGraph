module github.com/hristo-ganekov-sumup/genGraph

go 1.16

replace github.com/hristo-ganekov-sumup/genGraph/internal/tfstate => ./internal/tfstate

require (
	github.com/aws/aws-sdk-go v1.40.39
	github.com/davecgh/go-spew v1.1.1 // indirect
)
