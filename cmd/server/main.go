package main

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/fs"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	pb "rusprofile-wrapper/internal/rpc_server"
	"strconv"
	"time"
)

//go:embed swagger-ui
var content embed.FS

const (
	firstURLPref = "https://www.rusprofile.ru/ajax.php?query=%s&action=search&cacheKey=%.12f"
	secondURLPref = "https://www.rusprofile.ru"
)

var (
	port = flag.Int("port", 50051, "The server port")
	grpcServerEndpoint = flag.String("grpc-server-endpoint",  "localhost:50051", "gRPC server endpoint")
	endpointPort = flag.String("endpoint-port",  ":9090", "endpoint port")
)

func addHeaders(r *http.Request) {
	r.Header.Set("Accept-Language", "en-US,en;q=0.5")
	r.Header.Set("Accept", "application/json, text/plain, */*")
	r.Header.Set("TE", "Trailers")
	r.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
}

type JSONResponse struct {
	UL []struct {
		Name              string      `json:"name"`
		Link              string      `json:"link"`
		INN               string      `json:"inn"`
		CEO               string      `json:"ceo_name"`
	} `json:"ul"`
	IPCount int    `json:"ip_count"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type server struct {
	pb.UnimplementedCompanyInfoServiceServer
}

func (s *server) FetchCompanyInfo(ctx context.Context, in *pb.CompanyRequest) (*pb.CompanyResponse, error) {
	INN := strconv.FormatUint(in.GetINN(), 10)
	if len(INN) != 10 && len(INN) != 12 && len(INN) != 5 {
		return nil, fmt.Errorf("wrong INN len!")
	}
	client := &http.Client{Timeout: 10 * time.Second}

	firstURL := fmt.Sprintf(firstURLPref, INN, rand.Float64())
	response, err := doRequest(*client, firstURL, in.GetINN())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to request first url: %s", firstURL)
	}
	defer response.Body.Close()

	r := JSONResponse{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "failed to decode first response")
	}
	if len(r.UL) < 1 {
		return nil, fmt.Errorf("nothing was found, INN: %v", in.GetINN())
	}

	secondURL := fmt.Sprintf("%s%s", secondURLPref, r.UL[0].Link)
	response, err = doRequest(*client, secondURL, in.GetINN())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to request second url: %s", secondURL)
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create html parser")
	}
	kppSelection := doc.Find("[id='clip_kpp']")
	KPP, err := strconv.ParseUint(kppSelection.Text(), 10, 32)

	return &pb.CompanyResponse{
		INN:               in.GetINN(),
		KPP:               uint32(KPP),
		CompanyName:       r.UL[0].Name,
		DirectorFullName:  r.UL[0].CEO,
	}, nil
}

func doRequest(client http.Client, url string, INN uint64) (*http.Response, error) {
	rand.Seed(time.Now().Unix())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err,"failed to create new request with INN: %s", INN)
	}
	addHeaders(req)

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("http status code is not 200: %d INN: %s", response.StatusCode, INN)
	}
	return response, nil
}

func (s *server) runServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gRPCServer := grpc.NewServer()
	pb.RegisterCompanyInfoServiceServer(gRPCServer, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := gRPCServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runHTTP() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	gRPCmux := runtime.NewServeMux()

	err := pb.RegisterCompanyInfoServiceHandlerFromEndpoint(ctx, gRPCmux, *grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalf("failed to register http handler, %v", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/inn/", gRPCmux)
	mux.HandleFunc("/openapi.json", swaggerHandler)

	fsys, _ := fs.Sub(content, "swagger-ui")
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(fsys))))

	log.Printf("http server listening at %v", *endpointPort)
	if err := http.ListenAndServe(*endpointPort, mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile("internal/rpc_server/server.swagger.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  err.Error(),
			"result": nil,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func main() {
	flag.Parse()
	s := server{}
	go s.runServer()
	runHTTP()
}