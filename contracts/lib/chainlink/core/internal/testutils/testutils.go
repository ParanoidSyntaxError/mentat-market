package testutils

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"math/big"
	mrand "math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"go.uber.org/zap/zaptest/observer"

	// NOTE: To avoid circular dependencies, this package MUST NOT import
	// anything from "github.com/smartcontractkit/chainlink/v2/core"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

const (
	// Password just a password we use everywhere for testing
	Password = "16charlengthp4SsW0rD1!@#_"
)

// FixtureChainID matches the chain always added by fixtures.sql
// It is set to 0 since no real chain ever has this ID and allows a virtual
// "test" chain ID to be used without clashes
var FixtureChainID = big.NewInt(0)

// SimulatedChainID is the chain ID for the go-ethereum simulated backend
var SimulatedChainID = big.NewInt(1337)

// NewAddress return a random new address
func NewAddress() common.Address {
	return common.BytesToAddress(randomBytes(20))
}

// NewPrivateKeyAndAddress returns a new private key and the corresponding address
func NewPrivateKeyAndAddress(t testing.TB) (*ecdsa.PrivateKey, common.Address) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, address
}

// NewRandomPositiveInt64 returns a (non-cryptographically secure) random positive int64
func NewRandomPositiveInt64() int64 {
	id := mrand.Int63()
	return id
}

// NewRandomEVMChainID returns a suitable random chain ID that will not conflict
// with fixtures
func NewRandomEVMChainID() *big.Int {
	id := mrand.Int63n(math.MaxInt32) + 10000
	return big.NewInt(id)
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

// Random32Byte returns a random [32]byte
func Random32Byte() (b [32]byte) {
	copy(b[:], randomBytes(32))
	return b
}

// RandomizeName appends a random UUID to the provided name
func RandomizeName(n string) string {
	id := uuid.New().String()
	return n + id
}

// DefaultWaitTimeout is the default wait timeout. If you have a *testing.T, use WaitTimeout instead.
const DefaultWaitTimeout = 30 * time.Second

// WaitTimeout returns a timeout based on the test's Deadline, if available.
// Especially important to use in parallel tests, as their individual execution
// can get paused for arbitrary amounts of time.
func WaitTimeout(t *testing.T) time.Duration {
	if d, ok := t.Deadline(); ok {
		// 10% buffer for cleanup and scheduling delay
		return time.Until(d) * 9 / 10
	}
	return DefaultWaitTimeout
}

// Context returns a context with the test's deadline, if available.
func Context(tb testing.TB) context.Context {
	return tb.Context()
}

// MustParseURL parses the URL or fails the test
func MustParseURL(t testing.TB, input string) *url.URL {
	u, err := url.Parse(input)
	require.NoError(t, err)
	return u
}

// MustParseBigInt parses a big int value from string or fails the test
func MustParseBigInt(t *testing.T, input string) *big.Int {
	i := new(big.Int)
	_, err := fmt.Sscan(input, i)
	require.NoError(t, err)
	return i
}

// JSONRPCHandler is called with the method and request param(s).
// respResult will be sent immediately. notifyResult is optional, and sent after a short delay.
type JSONRPCHandler func(reqMethod string, reqParams gjson.Result) JSONRPCResponse

type JSONRPCResponse struct {
	Result, Notify string // raw JSON (i.e. quoted strings etc.)

	Error struct {
		Code    int
		Message string
	}
}

type testWSServer struct {
	t       *testing.T
	s       *httptest.Server
	mu      sync.RWMutex
	wsconns []*websocket.Conn
	wg      sync.WaitGroup
}

// NewWSServer starts a websocket server which invokes callback for each message received.
// If chainID is set, then eth_chainId calls will be automatically handled.
func NewWSServer(t *testing.T, chainID *big.Int, callback JSONRPCHandler) (ts *testWSServer) {
	ts = new(testWSServer)
	ts.t = t
	ts.wsconns = make([]*websocket.Conn, 0)
	handler := ts.newWSHandler(chainID, callback)
	ts.s = httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return
}

func (ts *testWSServer) Close() {
	if func() bool {
		ts.mu.Lock()
		defer ts.mu.Unlock()
		if ts.wsconns == nil {
			ts.t.Log("Test WS server already closed")
			return false
		}
		ts.s.CloseClientConnections()
		ts.s.Close()
		for _, ws := range ts.wsconns {
			ws.Close()
		}
		ts.wsconns = nil // nil indicates server closed
		return true
	}() {
		ts.wg.Wait()
	}
}

func (ts *testWSServer) WSURL() *url.URL {
	return WSServerURL(ts.t, ts.s)
}

func (ts *testWSServer) MustWriteBinaryMessageSync(t *testing.T, msg string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	conns := ts.wsconns
	if len(conns) != 1 {
		t.Fatalf("expected 1 conn, got %d", len(conns))
	}
	conn := conns[0]
	err := conn.WriteMessage(websocket.BinaryMessage, []byte(msg))
	require.NoError(t, err)
}

func (ts *testWSServer) newWSHandler(chainID *big.Int, callback JSONRPCHandler) (handler http.HandlerFunc) {
	if callback == nil {
		callback = func(method string, params gjson.Result) (resp JSONRPCResponse) { return }
	}
	t := ts.t
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ts.mu.Lock()
		if ts.wsconns == nil { // closed
			ts.mu.Unlock()
			return
		}
		ts.wg.Add(1)
		defer ts.wg.Done()
		conn, err := upgrader.Upgrade(w, r, nil)
		if !assert.NoError(t, err, "Failed to upgrade WS connection") {
			ts.mu.Unlock()
			return
		}
		defer conn.Close()
		ts.wsconns = append(ts.wsconns, conn)
		ts.mu.Unlock()

		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
					ts.t.Log("Websocket closing")
					return
				}
				ts.t.Logf("Failed to read message: %v", err)
				return
			}
			ts.t.Log("Received message", string(data))
			req := gjson.ParseBytes(data)
			if !req.IsObject() {
				ts.t.Logf("Request must be object: %v", req.Type)
				return
			}
			if e := req.Get("error"); e.Exists() {
				ts.t.Logf("Received jsonrpc error: %v", e)
				continue
			}
			m := req.Get("method")
			if m.Type != gjson.String {
				ts.t.Logf("Method must be string: %v", m.Type)
				return
			}

			var resp JSONRPCResponse
			if chainID != nil && m.String() == "eth_chainId" {
				resp.Result = `"0x` + chainID.Text(16) + `"`
			} else if m.String() == "eth_syncing" {
				resp.Result = "false"
			} else {
				resp = callback(m.String(), req.Get("params"))
			}
			id := req.Get("id")
			var msg string
			if resp.Error.Message != "" {
				msg = fmt.Sprintf(`{"jsonrpc":"2.0","id":%s,"error":{"code":%d,"message":"%s"}}`, id, resp.Error.Code, resp.Error.Message)
			} else {
				msg = fmt.Sprintf(`{"jsonrpc":"2.0","id":%s,"result":%s}`, id, resp.Result)
			}
			ts.t.Logf("Sending message: %v", msg)
			ts.mu.Lock()
			err = conn.WriteMessage(websocket.BinaryMessage, []byte(msg))
			ts.mu.Unlock()
			if err != nil {
				ts.t.Logf("Failed to write message: %v", err)
				return
			}

			if resp.Notify != "" {
				time.Sleep(100 * time.Millisecond)
				msg := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_subscription","params":{"subscription":"0x00","result":%s}}`, resp.Notify)
				ts.t.Log("Sending message", msg)
				ts.mu.Lock()
				err = conn.WriteMessage(websocket.BinaryMessage, []byte(msg))
				ts.mu.Unlock()
				if err != nil {
					ts.t.Logf("Failed to write message: %v", err)
					return
				}
			}
		}
	}
}

// WaitWithTimeout waits for the channel to close (or receive anything) and
// fatals the test if the default wait timeout is exceeded
func WaitWithTimeout(t *testing.T, ch <-chan struct{}, failMsg string) {
	select {
	case <-ch:
	case <-time.After(WaitTimeout(t)):
		t.Fatal(failMsg)
	}
}

// WSServerURL returns a ws:// url for the server
func WSServerURL(t *testing.T, s *httptest.Server) *url.URL {
	u, err := url.Parse(s.URL)
	require.NoError(t, err, "Failed to parse url")
	u.Scheme = "ws"
	return u
}

// IntToHex converts int to geth-compatible hex
func IntToHex(n int) string {
	return hexutil.EncodeBig(big.NewInt(int64(n)))
}

// TestInterval is just a sensible poll interval that gives fast tests without
// risk of spamming
const TestInterval = 100 * time.Millisecond

// AssertEventually calls assert.Eventually with default wait and tick durations.
func AssertEventually(t *testing.T, f func() bool) bool {
	return assert.Eventually(t, f, WaitTimeout(t), TestInterval/2)
}

// RequireEventually calls assert.Eventually with default wait and tick durations.
func RequireEventually(t *testing.T, f func() bool) {
	require.Eventually(t, f, WaitTimeout(t), TestInterval/2)
}

// RequireLogMessage fails the test if emitted logs don't contain the given message
func RequireLogMessage(t *testing.T, observedLogs *observer.ObservedLogs, msg string) {
	for _, l := range observedLogs.All() {
		if strings.Contains(l.Message, msg) {
			return
		}
	}
	t.Log("observed logs", observedLogs.All())
	t.Fatalf("expected observed logs to contain msg %q, but it didn't", msg)
}

// WaitForLogMessage waits until at least one log message containing the
// specified msg is emitted.
// NOTE: This does not "pop" messages so it cannot be used multiple times to
// check for new instances of the same msg. See WaitForLogMessageCount instead.
//
// Get a *observer.ObservedLogs like so:
//
//	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
//	lggr := logger.TestLogger(t, observedZapCore)
func WaitForLogMessage(t *testing.T, observedLogs *observer.ObservedLogs, msg string) (le observer.LoggedEntry) {
	RequireEventually(t, func() bool {
		for _, l := range observedLogs.All() {
			if strings.Contains(l.Message, msg) {
				le = l
				return true
			}
		}
		return false
	})
	return
}

func WaitForLogMessageWithField(t *testing.T, observedLogs *observer.ObservedLogs, msg, field, value string) (le observer.LoggedEntry) {
	RequireEventually(t, func() bool {
		for _, l := range observedLogs.All() {
			if strings.Contains(l.Message, msg) && strings.Contains(l.ContextMap()[field].(string), value) {
				le = l
				return true
			}
		}
		return false
	})
	return
}

// WaitForLogMessageCount waits until at least count log message containing the
// specified msg is emitted
func WaitForLogMessageCount(t *testing.T, observedLogs *observer.ObservedLogs, msg string, count int) {
	RequireEventually(t, func() bool {
		i := 0
		for _, l := range observedLogs.All() {
			if strings.Contains(l.Message, msg) {
				i++
				if i >= count {
					return true
				}
			}
		}
		return false
	})
}

// SkipShortDB skips tb during -short runs, and notes the DB dependency.
func SkipShortDB(tb testing.TB) {
	tests.SkipShort(tb, "DB dependency")
}

func AssertCount(t testing.TB, ds sqlutil.DataSource, tableName string, expected int64) {
	t.Helper()
	ctx := Context(t)
	var count int64
	err := ds.GetContext(ctx, &count, fmt.Sprintf(`SELECT count(*) FROM %s;`, tableName))
	require.NoError(t, err)
	require.Equal(t, expected, count)
}

// Ptr takes pointer of anything
func Ptr[T any](v T) *T {
	return &v
}

func MustDecodeBase64(s string) (b []byte) {
	var err error
	b, err = base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return
}

func MustRandBytes(n int) (b []byte) {
	b = make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return
}
