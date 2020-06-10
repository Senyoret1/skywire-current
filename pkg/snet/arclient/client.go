// Package arclient implements address resolver client
package arclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/SkycoinProject/dmsg/cipher"
	"github.com/SkycoinProject/skycoin/src/util/logging"
	"github.com/libp2p/go-reuseport"
	"nhooyr.io/websocket"

	"github.com/SkycoinProject/skywire-mainnet/internal/httpauth"
)

var log = logging.MustGetLogger("arclient")

const (
	bindSTCPRPath    = "/bind/stcpr"
	bindSUDPRPath    = "/bind/sudpr"
	resolveSTCPRPath = "/resolve/"
	resolveSTCPHPath = "/resolve_hole_punch/"
	resolveSUDPRPath = "/resolve_sudpr/"
	resolveSUDPHPath = "/resolve_sudph/"
	wsPath           = "/ws"
	addrChSize       = 1024
)

var (
	// ErrNoEntry means that there exists no entry for this PK.
	ErrNoEntry = errors.New("no entry for this PK")
	// ErrNotConnected is returned when PK is not connected.
	ErrNotConnected = errors.New("this PK is not connected")
)

// Error is the object returned to the client when there's an error.
type Error struct {
	Error string `json:"error"`
}

// APIClient implements DMSG discovery API client.
type APIClient interface {
	io.Closer
	LocalAddr() string
	BindSTCPR(ctx context.Context, port string) error
	BindSUDPR(ctx context.Context, port string) error
	ResolveSTCPR(ctx context.Context, pk cipher.PubKey) (string, error)
	ResolveSTCPH(ctx context.Context, pk cipher.PubKey) (string, error)
	ResolveSUDPR(ctx context.Context, pk cipher.PubKey) (string, error)
	ResolveSUDPH(ctx context.Context, pk cipher.PubKey) (string, error)
	WS(ctx context.Context, dialCh <-chan cipher.PubKey) (<-chan RemoteVisor, error)
}

type key struct {
	remoteAddr string
	pk         cipher.PubKey
	sk         cipher.SecKey
}

var clients = make(map[key]*httpClient)

// httpClient implements Client for uptime tracker API.
type httpClient struct {
	client    *httpauth.Client
	localAddr string
	pk        cipher.PubKey
	sk        cipher.SecKey
	wsConn    *websocket.Conn
	addrCh    <-chan RemoteVisor
}

// NewHTTP creates a new client setting a public key to the client to be used for auth.
// When keys are set, the client will sign request before submitting.
// The signature information is transmitted in the header using:
// * SW-Public: The specified public key
// * SW-Nonce:  The nonce for that public key
// * SW-Sig:    The signature of the payload + the nonce
func NewHTTP(remoteAddr string, pk cipher.PubKey, sk cipher.SecKey) (APIClient, error) {
	key := key{
		remoteAddr: remoteAddr,
		pk:         pk,
		sk:         sk,
	}

	// Same clients would have nonce collisions. Client should be reused in this case.
	if client, ok := clients[key]; ok {
		return client, nil
	}

	log.Infof("Creating arclient, key = %v", key)

	httpAuthClient, err := httpauth.NewClient(context.Background(), remoteAddr, pk, sk)
	if err != nil {
		return nil, fmt.Errorf("address resolver httpauth: %w", err)
	}

	client := &httpClient{
		client:    httpAuthClient,
		pk:        pk,
		sk:        sk,
		localAddr: "",
	}

	transport := &http.Transport{
		DialContext: func(_ context.Context, network, remoteAddr string) (conn net.Conn, err error) {
			conn, err = reuseport.Dial(network, client.localAddr, remoteAddr)
			if err == nil && client.localAddr == "" {
				client.localAddr = conn.LocalAddr().String()
			}

			return conn, err
		},
		DisableKeepAlives: false,
	}

	httpAuthClient.SetTransport(transport)

	clients[key] = client

	return client, nil
}

func (c *httpClient) LocalAddr() string {
	return c.localAddr
}

// Get performs a new GET request.
func (c *httpClient) Get(ctx context.Context, path string) (*http.Response, error) {
	addr := c.client.Addr() + path

	req, err := http.NewRequest(http.MethodGet, addr, new(bytes.Buffer))
	if err != nil {
		return nil, err
	}

	return c.client.Do(req.WithContext(ctx))
}

// Post performs a POST request.
func (c *httpClient) Post(ctx context.Context, path string, payload interface{}) (*http.Response, error) {
	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	addr := c.client.Addr() + path

	req, err := http.NewRequest(http.MethodPost, addr, body)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req.WithContext(ctx))
}

// Websocket performs a new websocket request.
func (c *httpClient) Websocket(ctx context.Context, path string) (*websocket.Conn, error) {
	header, err := c.client.Header()
	if err != nil {
		return nil, err
	}

	dialOpts := &websocket.DialOptions{
		HTTPClient: c.client.ReuseClient(),
		HTTPHeader: header,
	}

	addr, err := url.Parse(c.client.Addr())
	if err != nil {
		return nil, err
	}
	switch addr.Scheme {
	case "http":
		addr.Scheme = "ws"
	case "https":
		addr.Scheme = "wss"
	}

	addr.Path = path

	conn, resp, err := websocket.Dial(ctx, addr.String(), dialOpts)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		c.client.IncrementNonce()
	}

	return conn, nil
}

// BindRequest stores bind request values.
type BindRequest struct {
	Port string `json:"port"`
}

// BindSTCPR binds client PK to IP:port on address resolver.
func (c *httpClient) BindSTCPR(ctx context.Context, port string) error {
	req := BindRequest{
		Port: port,
	}

	resp, err := c.Post(ctx, bindSTCPRPath, req)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.WithError(err).Warn("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d, error: %w", resp.StatusCode, extractError(resp.Body))
	}

	return nil
}

// BindSTCPR binds client PK to IP:port on address resolver.
func (c *httpClient) BindSUDPR(ctx context.Context, port string) error {
	req := BindRequest{
		Port: port,
	}

	resp, err := c.Post(ctx, bindSUDPRPath, req)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.WithError(err).Warn("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d, error: %w", resp.StatusCode, extractError(resp.Body))
	}

	return nil
}

// ResolveResponse stores response response values.
type ResolveResponse struct {
	Addr string `json:"addr"`
}

func (c *httpClient) ResolveSTCPR(ctx context.Context, pk cipher.PubKey) (string, error) {
	resp, err := c.Get(ctx, resolveSTCPRPath+pk.String())
	if err != nil {
		return "", err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.WithError(err).Warn("Failed to close response body")
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrNoEntry
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status: %d, error: %w", resp.StatusCode, extractError(resp.Body))
	}

	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var resolveResp ResolveResponse

	if err := json.Unmarshal(rawBody, &resolveResp); err != nil {
		return "", err
	}

	return resolveResp.Addr, nil
}

func (c *httpClient) ResolveSTCPH(ctx context.Context, pk cipher.PubKey) (string, error) {
	resp, err := c.Get(ctx, resolveSTCPHPath+pk.String())
	if err != nil {
		return "", err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.WithError(err).Warn("Failed to close response body")
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrNoEntry
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status: %d, error: %w", resp.StatusCode, extractError(resp.Body))
	}

	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var resolveResp ResolveResponse

	if err := json.Unmarshal(rawBody, &resolveResp); err != nil {
		return "", err
	}

	return resolveResp.Addr, nil
}

func (c *httpClient) ResolveSUDPR(ctx context.Context, pk cipher.PubKey) (string, error) {
	resp, err := c.Get(ctx, resolveSUDPRPath+pk.String())
	if err != nil {
		return "", err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.WithError(err).Warn("Failed to close response body")
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrNoEntry
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status: %d, error: %w", resp.StatusCode, extractError(resp.Body))
	}

	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var resolveResp ResolveResponse

	if err := json.Unmarshal(rawBody, &resolveResp); err != nil {
		return "", err
	}

	return resolveResp.Addr, nil
}

func (c *httpClient) ResolveSUDPH(ctx context.Context, pk cipher.PubKey) (string, error) {
	resp, err := c.Get(ctx, resolveSUDPHPath+pk.String())
	if err != nil {
		return "", err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.WithError(err).Warn("Failed to close response body")
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrNoEntry
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status: %d, error: %w", resp.StatusCode, extractError(resp.Body))
	}

	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var resolveResp ResolveResponse

	if err := json.Unmarshal(rawBody, &resolveResp); err != nil {
		return "", err
	}

	return resolveResp.Addr, nil
}

// RemoteVisor contains public key and address of remote visor.
type RemoteVisor struct {
	PK   cipher.PubKey
	Addr string
}

func (c *httpClient) WS(ctx context.Context, dialCh <-chan cipher.PubKey) (<-chan RemoteVisor, error) {
	if c.addrCh == nil {
		if err := c.initWS(ctx, dialCh); err != nil {
			return nil, err
		}
	}

	return c.addrCh, nil
}

func (c *httpClient) initWS(ctx context.Context, dialCh <-chan cipher.PubKey) error {
	// TODO(nkryuchkov): Ensure this works correctly with closed channels and connections.
	addrCh := make(chan RemoteVisor, addrChSize)

	conn, err := c.Websocket(ctx, wsPath)
	if err != nil {
		return err
	}

	c.wsConn = conn

	go func(conn *websocket.Conn, addrCh chan<- RemoteVisor) {
		defer func() {
			close(addrCh)
		}()

		for {
			kind, rawMsg, err := conn.Read(context.TODO())
			if err != nil {
				log.Errorf("Failed to read WS message: %v", err)
				return
			}

			log.Infof("New WS message of type %v: %v", kind.String(), string(rawMsg))

			var remote RemoteVisor
			if err := json.Unmarshal(rawMsg, &remote); err != nil {
				log.Errorf("Failed to read unmarshal message: %v", err)
				continue
			}

			addrCh <- remote
		}
	}(conn, addrCh)

	go func(conn *websocket.Conn, dialCh <-chan cipher.PubKey) {
		for pk := range dialCh {
			if err := conn.Write(ctx, websocket.MessageText, []byte(pk.String())); err != nil {
				log.Errorf("Failed to write to %v: %v", pk, err)
				return
			}
		}
	}(conn, dialCh)

	c.addrCh = addrCh

	return nil
}

func (c *httpClient) Close() error {
	defer func() {
		c.wsConn = nil
	}()

	return c.wsConn.Close(websocket.StatusNormalClosure, "client closed")
}

// extractError returns the decoded error message from Body.
func extractError(r io.Reader) error {
	var apiError Error

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &apiError); err != nil {
		return errors.New(string(body))
	}

	return errors.New(apiError.Error)
}
