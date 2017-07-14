package main

import (
	"flag"
	"fmt"
	"os"

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"jaxf-github.fanatics.corp/apparel/partner-service/pkg/db"
	"jaxf-github.fanatics.corp/apparel/partner-service/pkg/db/dbconfig"
	"jaxf-github.fanatics.corp/apparel/partner-service/pkg/endpoints"
	"jaxf-github.fanatics.corp/apparel/partner-service/pkg/pb"
	"jaxf-github.fanatics.corp/apparel/partner-service/pkg/service"
	tg "jaxf-github.fanatics.corp/apparel/partner-service/pkg/transport_grpc"
)

func main() {
	// command line flags
	grpcAddr := flag.String("grpcAddr", ":8081", "gRPC listen address")
	dom := flag.String("domain", "localhost", "domain name of service")
	certPath := flag.String("certPath", "./tls/test/test.cert.pem", "path to ssl cert file")
	keyPath := flag.String("keyPath", "./tls/test/test.key.pem", "path to ssl key file")
	sec := flag.Bool("sec", false, "use ssl cert")
	flag.Parse()

	var config *tls.Config
	host := fmt.Sprintf("%v%v", *dom, *grpcAddr)

	// set up logger
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// set up tls
	if *sec {

		cert, err := tls.LoadX509KeyPair(*certPath, *keyPath)
		if err != nil {
			err = errors.Wrap(err, "failed to create certificate")
			logger.Log("err", err)
			panic(err)
		}

		pem, err := ioutil.ReadFile(*certPath)
		if err != nil {
			err = errors.Wrap(err, "failed to create pem")
			logger.Log("err", err)
			panic(err)
		}

		certPool := x509.NewCertPool()
		ok := certPool.AppendCertsFromPEM(pem)
		if !ok {
			err = errors.New("failed to append cert from pem")
			logger.Log("err", err)
			panic(err)
		}

		config = &tls.Config{
			ServerName:               host,
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			Certificates: []tls.Certificate{cert},
			RootCAs:      certPool,
		}
	}

	// set up db
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		panic(err)
	}
	defer conn.Close()

	// Make service and endpoints
	svc := service.New(logger, db.NewPartnerServiceQuerier(conn))
	eps := endpoints.New(svc, logger)

	// Mechanical domain.
	errc := make(chan error)

	// Interrupt handler.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		ln, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errc <- err
			return
		}

		srv := tg.MakeGRPCServer(eps, logger)
		var opts []grpc.ServerOption
		if *sec {
			creds := credentials.NewTLS(config)
			fmt.Printf("%+v", &creds)
			if err != nil {
				err = errors.Wrap(err, "failed to create new server tls from file")
				logger.Log("err", err)
				panic(err)
			}
			opts = append(opts, grpc.Creds(creds))
		}
		s := grpc.NewServer(opts...)
		pb.RegisterPartnerServiceServer(s, srv)
		logger.Log("success", "it works")
		logger.Log("addr", *grpcAddr)
		errc <- s.Serve(ln)
	}()

	// Run!
	logger.Log("exit", <-errc)
}
