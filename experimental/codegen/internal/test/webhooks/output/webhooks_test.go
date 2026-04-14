package output

import (
	"net/http"
	"testing"

	"github.com/oapi-codegen/oapi-codegen-exp/experimental/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Verify schema types can be instantiated
func TestSchemaInstantiation(t *testing.T) {
	reg := WebhookRegistration{URL: "https://example.com/hook"}
	assert.Equal(t, "https://example.com/hook", reg.URL)

	resp := WebhookRegistrationResponse{}
	_ = resp.ID // UUID field exists

	p := Person{Name: "Alice"}
	assert.Equal(t, "Alice", p.Name)
}

// Verify enum constants exist
func TestWebhookKindEnum(t *testing.T) {
	assert.Equal(t, PostAPIWebhookKindParameter("enterEvent"), EnterEvent)
	assert.Equal(t, PostAPIWebhookKindParameter("exitEvent"), ExitEvent)
}

// Verify type aliases for webhook request bodies
func TestWebhookRequestBodyAliases(t *testing.T) {
	var enterBody EnterEventJSONRequestBody = Person{Name: "Bob"}
	var exitBody ExitEventJSONRequestBody = Person{Name: "Carol"}
	assert.Equal(t, "Bob", enterBody.Name)
	assert.Equal(t, "Carol", exitBody.Name)
}

// Verify ServerInterface is implementable (server for registration endpoints)
type testServer struct{}

func (s *testServer) DeregisterWebhook(w http.ResponseWriter, r *http.Request, id types.UUID) {
	w.WriteHeader(http.StatusNoContent)
}
func (s *testServer) RegisterWebhook(w http.ResponseWriter, r *http.Request, kind string) {
	w.WriteHeader(http.StatusCreated)
}

func TestServerInterfaceImplementable(t *testing.T) {
	var si ServerInterface = &testServer{}
	_ = Handler(si)
}

// Verify WebhookInitiator can be created and request builders work
func TestWebhookInitiator(t *testing.T) {
	initiator, err := NewWebhookInitiator()
	require.NoError(t, err)
	require.NotNil(t, initiator)

	// Verify the initiator satisfies the interface
	var _ WebhookInitiatorInterface = initiator
}

func TestWebhookRequestBuilders(t *testing.T) {
	body := Person{Name: "Dave"}

	req, err := NewEnterEventWebhookRequest("https://example.com/enter", body)
	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, req.Method)
	assert.Equal(t, "https://example.com/enter", req.URL.String())

	req, err = NewExitEventWebhookRequest("https://example.com/exit", body)
	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, req.Method)
	assert.Equal(t, "https://example.com/exit", req.URL.String())
}

// Verify WebhookReceiverInterface is implementable
type testReceiver struct{}

func (r *testReceiver) HandleEnterEventWebhook(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (r *testReceiver) HandleExitEventWebhook(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestWebhookReceiverInterfaceImplementable(t *testing.T) {
	var ri WebhookReceiverInterface = &testReceiver{}
	handler := EnterEventWebhookHandler(ri, func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, err.Error(), http.StatusBadRequest)
	})
	require.NotNil(t, handler)
}

// Verify ApplyDefaults methods exist
func TestApplyDefaults(t *testing.T) {
	reg := &WebhookRegistration{}
	reg.ApplyDefaults()
	resp := &WebhookRegistrationResponse{}
	resp.ApplyDefaults()
	p := &Person{}
	p.ApplyDefaults()
}
