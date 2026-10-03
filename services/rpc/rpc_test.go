/*
Merlin is a post-exploitation command and control framework.

This file is part of Merlin.
Copyright (C) 2026  Russel Van Tuyl

Merlin is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as
published by the Free Software Foundation, either version 3 of the License, or any later version.

Merlin is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with Merlin.  If not, see <http://www.gnu.org/licenses/>.
*/

package rpc

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/Ne0nd0g/merlin-cli/rpc"
)

// fakeMerlinServer implements just enough of the Merlin gRPC service for the CLI's Connect() to
// succeed: Register (returns a UUID, after validating the auth metadata) and Listen (holds the
// stream open until the client disconnects). Everything else is UnimplementedMerlinServer.
type fakeMerlinServer struct {
	pb.UnimplementedMerlinServer
	id       string
	password string
}

func (f *fakeMerlinServer) Register(ctx context.Context, _ *emptypb.Empty) (*pb.ID, error) {
	// The CLI's unary interceptor must have added the password as "authorization" metadata.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("authorization")) == 0 || md.Get("authorization")[0] != f.password {
		return nil, context.Canceled
	}
	return &pb.ID{Id: f.id}, nil
}

func (f *fakeMerlinServer) Listen(_ *pb.ID, stream pb.Merlin_ListenServer) error {
	<-stream.Context().Done() // keep the stream open until the client goes away
	return nil
}

// testTLSCert returns a self-signed cert valid for 127.0.0.1 / localhost for the in-process server.
func testTLSCert(t *testing.T) tls.Certificate {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		t.Fatalf("serial: %v", err)
	}
	tmpl := x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: "merlin-cli-test"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}
	der, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}
}

// TestConnect starts an in-process TLS gRPC server implementing a minimal Merlin service, then
// drives the real Connect() against it to verify the client dials (grpc.NewClient), authenticates
// via the metadata interceptor, and completes the Register handshake.
func TestConnect(t *testing.T) {
	service = nil // reset the package singleton so NewRPCService builds fresh

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer func() { _ = lis.Close() }()

	const password = "test-password"
	wantID := uuid.New().String()

	cert := testTLSCert(t)
	srv := grpc.NewServer(grpc.Creds(credentials.NewServerTLSFromCert(&cert)))
	pb.RegisterMerlinServer(srv, &fakeMerlinServer{id: wantID, password: password})
	go func() { _ = srv.Serve(lis) }()
	defer srv.Stop()

	svc, err := NewRPCService(password, false, "", "", "")
	if err != nil {
		t.Fatalf("NewRPCService: %v", err)
	}

	if err = svc.Connect(lis.Addr().String()); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	if svc.id.String() != wantID {
		t.Fatalf("registered id = %s, want %s", svc.id, wantID)
	}
}
