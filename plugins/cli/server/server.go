package server

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/google/uuid"
	"github.com/goxgen/goxgen/utils/mapper"
	"github.com/rs/cors"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
)

type MapperKey int
type LoggerKey int

const (
	mapperKey MapperKey = iota
	loggerKey LoggerKey = iota
)

type Server struct {
	cliContext *cli.Context
	mapper     *mapper.Mapper
	logger     *zap.Logger
}

// New creates a new server instance
func New(ctx *cli.Context) (s *Server, err error) {
	s = &Server{
		cliContext: ctx,
	}
	s.mapper, err = s.prepareMapper(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare mapper: %w", err)
	}
	s.logger, err = s.prepareLogger(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare logger: %w", err)
	}
	return s, nil
}

type Data struct {
	Host                     string
	Port                     string
	HTTPS                    bool
	AppPath                  string
	GraphqlURIPath           string
	GraphqlURL               string
	GraphqlPlaygroundEnabled bool
	GraphqlPlaygroundURIPath string
	GraphqlPlaygroundUrl     string
}

type Constructor func(ctx *cli.Context) (*handler.Server, error)

func (s *Server) ListenAndServe(serverConstructor Constructor) error {

	defer s.logger.Sync()

	data := s.GetDataFromCliContext()

	httpHandler, err := serverConstructor(s.cliContext)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	if data.GraphqlPlaygroundEnabled {
		mux.Handle(data.GraphqlPlaygroundURIPath, playground.Handler("Playground", data.GraphqlURIPath))
		s.logger.Info("Serving graphql playground", zap.String("url", data.GraphqlPlaygroundUrl))
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	mux.Handle(data.GraphqlURIPath, c.Handler(s.commonMiddleware(httpHandler)))
	s.logger.Info("Serving graphql", zap.String("url", data.GraphqlURL))

	srv := &http.Server{
		Addr:    data.Port,
		Handler: mux,
	}

	return srv.ListenAndServe()
}

// TestServer creates a new test server instance
func (s *Server) TestServer(ctx *cli.Context, serverConstructor Constructor) (testSrv *httptest.Server, cancel func()) {
	tempDB := os.TempDir() + "/" + uuid.New().String() + ".db"

	err := ctx.Set("DatabaseSourceName", "file:"+tempDB+"?mode=rwc&cache=shared&_fk=1")
	if err != nil {
		panic(err)
	}

	err = ctx.Set("DatabaseDriver", "sqlite")
	if err != nil {
		panic(err)
	}

	srv, err := serverConstructor(ctx)
	if err != nil {
		panic(err)
	}
	testSrv = httptest.NewServer(s.commonMiddleware(srv))

	// Cleanup the temp db
	cancel = func() {
		err := os.Remove(tempDB)
		if err != nil {
			panic(err)
		}
	}
	return testSrv, cancel
}

// commonMiddleware is a middleware that adds a key-value pair to the context
func (s *Server) commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Add key-value pair to context
		ctx := r.Context()
		ctx = context.WithValue(ctx, mapperKey, s.mapper)
		ctx = context.WithValue(ctx, loggerKey, s.logger)

		// Update the request with the new context
		r = r.WithContext(ctx)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// prepareMapper creates a new mapper instance
func (s *Server) prepareMapper(ctx *cli.Context) (*mapper.Mapper, error) {
	mpr, err := mapper.New()
	if err != nil {
		return nil, err
	}
	return mpr, nil
}

// prepareLogger creates a new logger instance
// Use the DevMode flag to enable development mode
// Use the LogLevel flag to set the log level
// Using zap logger
func (s *Server) prepareLogger(ctx *cli.Context) (*zap.Logger, error) {
	var log *zap.Logger
	if ctx.Bool("DevMode") {
		log, _ = zap.NewDevelopment()
	} else {
		log, _ = zap.NewProduction()
	}

	providedLevel := ctx.String("LogLevel")
	level, err := zapcore.ParseLevel(providedLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}
	log = log.WithOptions(zap.IncreaseLevel(level))
	return log, nil
}

// GetMapper returns the mapper instance from the context
func GetMapper(ctx context.Context) *mapper.Mapper {
	return ctx.Value(mapperKey).(*mapper.Mapper)
}

// GetLogger returns the logger instance from the context
func GetLogger(ctx context.Context) *zap.Logger {
	return ctx.Value(loggerKey).(*zap.Logger)
}

// GetDataFromCliContext returns the data from the cli context for the server
func (s *Server) GetDataFromCliContext() *Data {
	data := &Data{}
	data.HTTPS = s.cliContext.Bool("HTTPS")
	data.Host = s.cliContext.String("Host")
	data.Port = ":" + strconv.Itoa(s.cliContext.Int("Port"))
	data.AppPath = strings.Trim(s.cliContext.String("AppPath"), "/")
	if data.AppPath == "" {
		data.AppPath = "/"
	} else {
		data.AppPath = "/" + data.AppPath + "/"
	}
	data.GraphqlURL = s.cliContext.String("GraphqlURL")
	data.GraphqlURIPath = s.cliContext.String("GraphqlURIPath")
	data.GraphqlPlaygroundURIPath = s.cliContext.String("GraphqlPlaygroundURIPath")
	proto := "http://"
	data.GraphqlPlaygroundEnabled = s.cliContext.Bool("GraphqlPlaygroundEnabled")

	if data.HTTPS {
		proto = "https://"
	}

	if data.GraphqlURL == "" {
		data.GraphqlURL += proto + data.Host + data.Port + data.GraphqlURIPath
	}

	if data.GraphqlPlaygroundEnabled {
		data.GraphqlPlaygroundUrl = proto + data.Host + data.Port + data.GraphqlPlaygroundURIPath
	}

	return data
}
