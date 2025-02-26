package v1rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

/*const (
	maxRequestBytes   = 1.5 * 1024 * 1024
	grpcOverheadBytes = 512 * 1024
	maxStreams        = 1<<32 - 1 // math.MaxUint32와 동일
	maxSendBytes      = 1<<31 - 1 // math.MaxInt32와 동일
)*/

const (
	defaultMaxRequestBytes   = 1.5 * 1024 * 1024
	defaultGrpcOverheadBytes = 512 * 1024
	defaultMaxStreams        = 1<<32 - 1 // math.MaxUint32와 동일
	defaultMaxSendBytes      = 1<<31 - 1 // math.MaxInt32와 동일
)

var (
	address = ":50052"
)

func init() {
	// TODO: Prometheus 적용 예정
}

// 값이 없거나 잘못된 경우 defaultVal 을 반환한다.
func getEnvInt(key string, defaultVal int) int {
	s := os.Getenv(key)
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("Invalid value for %s: %v. Using default: %d", key, err, defaultVal)
		return defaultVal
	}
	return val
}

// 기존 코드에서는 net.Listen 실패 시 log.Fatalf로 바로 종료했으나,
// 수정된 버전에서는 에러를 반환해 상위에서 graceful하게 처리할 수 있게 함.

// 예시 Unary 인터셉터 (로깅용)
// gRPC 요청을 받을 때마다 요청 메서드와 에러 정보를 로깅함.
func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("Received request for %s", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("Method %s error: %v", info.FullMethod, err)
	}
	return resp, err
}

func Server() error {
	// 기존 코드:
	// lis, err := net.Listen("tcp", address)
	// if err != nil {
	//     log.Fatalf("failed to listen: %v", err)
	// }
	// 수정: 에러 발생 시 에러 반환
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	// TLS 설정 예시 (TLS가 필요하면 주석 풀고 사용)
	/*
		certFile := "server.crt"
		keyFile := "server.key"
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			return fmt.Errorf("failed to load TLS credentials: %w", err)
		}
	*/

	// 환경 변수로 옵션 값을 오버라이드할 수 있음
	maxRecvMsgSize := getEnvInt("GRPC_MAX_RECV_MSG_SIZE", int(defaultMaxRequestBytes+defaultGrpcOverheadBytes))
	maxSendMsgSize := getEnvInt("GRPC_MAX_SEND_MSG_SIZE", defaultMaxSendBytes)
	maxConcurrentStreams := getEnvInt("GRPC_MAX_CONCURRENT_STREAMS", defaultMaxStreams)

	// gRPC 서버 옵션 설정
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(maxRecvMsgSize),
		grpc.MaxSendMsgSize(maxSendMsgSize),
		grpc.MaxConcurrentStreams(uint32(maxConcurrentStreams)),
		grpc.UnaryInterceptor(loggingInterceptor),
		// grpc.Creds(creds), // TLS 사용 시 활성화
	}

	grpcServer := grpc.NewServer(opts...)

	// 기존 서비스 등록 (주석 처리된 상태)
	// RegisterJobsManSrv(grpcServer)
	// RegisterHelloWorldManSrv(grpcServer)
	// 수정: 헬스 체크 서비스 등록
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// 기존: Reflection 서비스 등록
	reflection.Register(grpcServer)

	log.Printf("gRPC server started, address: %s", address)

	// graceful shutdown 처리 추가
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("Received signal: %v. Initiating graceful shutdown...", sig)
		// GracefulStop은 현재 처리 중인 요청을 모두 완료한 후 서버를 중지함.
		grpcServer.GracefulStop()
	}()

	// 서버 시작
	serveErr := grpcServer.Serve(lis)
	if serveErr != nil {
		if !strings.Contains(serveErr.Error(), "use of closed network connection") {
			log.Printf("gRPC server returned with error: %v", serveErr)
		} else {
			log.Printf("gRPC server is shut down")
		}
	}
	return serveErr
}
