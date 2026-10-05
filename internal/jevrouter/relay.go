package jevrouter

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"strings"
	"sync"
)

const OpenRouterAPI = "https://openrouter.ai/api/v1"
const relayBodyLimit = 1 << 20

type Refusal struct {
	Status      int
	LimitSource string
}
type Observation struct {
	Models, Providers []string
	ResponseID        string
	Refusal           *Refusal
}

// Relay owns its loopback listener and retains only response observations.
type Relay struct {
	mu           sync.Mutex
	observations map[string]Observation
	server       *http.Server
	baseURL      string
	done         chan struct{}
}
type relayTokenKey struct{}

func StartRelay(upstream string, client *http.Client) (*Relay, error) {
	target, err := url.Parse(upstream)
	if err != nil || (target.Scheme != "http" && target.Scheme != "https") || target.Host == "" || target.User != nil || target.RawQuery != "" || target.Fragment != "" {
		return nil, errors.New("start router relay: invalid upstream URL")
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("listen for router relay: %w", err)
	}
	r := &Relay{observations: make(map[string]Observation), baseURL: "http://" + listener.Addr().String(), done: make(chan struct{})}
	defaultTransport := http.DefaultTransport.(*http.Transport).Clone()
	defaultTransport.DisableCompression = true
	var transport http.RoundTripper = defaultTransport
	if client != nil && client.Transport != nil {
		transport = client.Transport
	}
	proxy := &httputil.ReverseProxy{
		Transport:     transport,
		FlushInterval: -1,
		Rewrite: func(req *httputil.ProxyRequest) {
			_, rest, _ := strings.Cut(strings.TrimPrefix(req.In.URL.EscapedPath(), "/"), "/")
			escaped := strings.TrimRight(target.EscapedPath(), "/") + "/" + rest
			path, _ := url.PathUnescape(escaped)
			req.Out.URL.Scheme, req.Out.URL.Host = target.Scheme, target.Host
			req.Out.URL.Path, req.Out.URL.RawPath = path, escaped
			req.Out.URL.RawQuery = req.In.URL.RawQuery
			req.Out.Host = target.Host
		},
		ErrorHandler: func(w http.ResponseWriter, _ *http.Request, _ error) { w.WriteHeader(http.StatusBadGateway) },
		ErrorLog:     log.New(io.Discard, "", 0),
	}
	proxy.ModifyResponse = func(response *http.Response) error {
		token, _ := response.Request.Context().Value(relayTokenKey{}).(string)
		var refusal *Refusal
		if response.StatusCode == http.StatusPaymentRequired {
			r.mu.Lock()
			observation, ok := r.observations[token]
			if ok && observation.Refusal == nil {
				refusal = &Refusal{Status: response.StatusCode}
				observation.Refusal = refusal
				r.observations[token] = observation
			}
			r.mu.Unlock()
		}
		mediaType, _, _ := mime.ParseMediaType(response.Header.Get("Content-Type"))
		response.Body = &relayBody{ReadCloser: response.Body, stream: mediaType == "text/event-stream", observe: func(data []byte) { r.observe(token, data, refusal) }}
		return nil
	}
	r.server = &http.Server{ErrorLog: log.New(io.Discard, "", 0), Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, _, _ := strings.Cut(strings.TrimPrefix(req.URL.EscapedPath(), "/"), "/")
		r.mu.Lock()
		_, open := r.observations[token]
		r.mu.Unlock()
		if !open {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		proxy.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), relayTokenKey{}, token)))
	})}
	go func() { defer close(r.done); _ = r.server.Serve(listener) }()
	return r, nil
}

func (r *Relay) Open() (token, baseURL string) {
	token = rand.Text()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.observations[token] = Observation{}
	return token, r.baseURL + "/" + token
}

func (r *Relay) Take(token string) Observation {
	r.mu.Lock()
	defer r.mu.Unlock()
	observation, ok := r.observations[token]
	if ok {
		r.observations[token] = Observation{}
	}
	// Refusal may still be enriched by a response finishing concurrently.
	if observation.Refusal != nil {
		copy := *observation.Refusal
		observation.Refusal = &copy
	}
	return observation
}

func (r *Relay) Release(token string) (open int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.observations, token)
	return len(r.observations)
}

func (r *Relay) Close(ctx context.Context) error {
	err := r.server.Shutdown(ctx)
	if err != nil {
		_ = r.server.Close()
	}
	<-r.done
	r.mu.Lock()
	clear(r.observations)
	r.mu.Unlock()
	if err != nil {
		return fmt.Errorf("close router relay: %w", err)
	}
	return nil
}

func (r *Relay) observe(token string, data []byte, refusal *Refusal) {
	var fields struct {
		ID       string `json:"id"`
		Model    string `json:"model"`
		Provider string `json:"provider"`
	}
	if json.Unmarshal(data, &fields) != nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	observation, ok := r.observations[token]
	if !ok {
		return
	}
	if fields.Model != "" && !slices.Contains(observation.Models, fields.Model) {
		observation.Models = append(observation.Models, fields.Model)
	}
	if fields.Provider != "" && !slices.Contains(observation.Providers, fields.Provider) {
		observation.Providers = append(observation.Providers, fields.Provider)
	}
	if fields.ID != "" {
		observation.ResponseID = fields.ID
	}
	if refusal != nil && observation.Refusal == refusal {
		var credit struct {
			Error struct {
				Metadata struct {
					LimitSource string `json:"limit_source"`
				} `json:"metadata"`
			} `json:"error"`
		}
		if json.Unmarshal(data, &credit) == nil && credit.Error.Metadata.LimitSource != "" {
			refusal.LimitSource = credit.Error.Metadata.LimitSource
		}
	}
	r.observations[token] = observation
}

// relayBody copies only for parsing; Read always returns the upstream bytes.
// Partial or oversized JSON is discarded without affecting forwarding.
type relayBody struct {
	io.ReadCloser
	stream    bool
	buffer    []byte
	oversized bool
	observe   func([]byte)
}

func (body *relayBody) Read(p []byte) (int, error) {
	n, err := body.ReadCloser.Read(p)
	data := p[:n]
	if body.stream {
		for len(data) > 0 {
			line, rest, newline := bytes.Cut(data, []byte{'\n'})
			body.collect(line)
			if newline {
				body.finishLine()
				data = rest
			} else {
				break
			}
		}
	} else {
		body.collect(data)
	}
	if err != nil {
		if errors.Is(err, io.EOF) {
			if body.stream {
				body.finishLine()
			} else if !body.oversized {
				body.observe(body.buffer)
			}
		}
		body.buffer = nil
	}
	return n, err
}
func (body *relayBody) collect(data []byte) {
	if body.oversized {
		return
	}
	// The limit bounds a whole body, or one line of a stream; finishLine resets
	// it for the next line.
	if len(body.buffer)+len(data) > relayBodyLimit {
		body.buffer = nil
		body.oversized = true
		return
	}
	body.buffer = append(body.buffer, data...)
}
func (body *relayBody) finishLine() {
	if !body.oversized && bytes.HasPrefix(body.buffer, []byte("data: ")) {
		body.observe(bytes.TrimSuffix(body.buffer[6:], []byte{'\r'}))
	}
	body.buffer = nil
	body.oversized = false
}
func (body *relayBody) Close() error {
	body.buffer = nil
	body.observe = nil
	return body.ReadCloser.Close()
}
