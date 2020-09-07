package main

import (
	"testing"

	"github.com/danpantry/hello-grpc/protocol"
	"google.golang.org/grpc"
)

func TestCanExecuteWorksCorrectly(t *testing.T) {
	claims := protocol.Claims{
		Permissions: []string{"foobar"},
	}

	info := grpc.UnaryServerInfo{
		FullMethod: "/Foo/Bar",
	}

	var p servicePermissionMap
	if p.CanExecute(claims, &info) == true {
		t.Error("CanExecute(): expected false, got true")
	}

	p.AddPermission("Foo", "Bar", "foobar")
	if !p.CanExecute(claims, &info) {
		t.Error("CanExecute(): expected true, got false")
	}
}
