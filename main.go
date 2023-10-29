package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	_ = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}
