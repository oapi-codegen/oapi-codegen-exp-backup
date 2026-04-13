package output

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Verify schema types can be instantiated
func TestSchemaInstantiation(t *testing.T) {
	req := TreePlantingRequest{
		Location:    "north meadow",
		Kind:        "oak",
		CallbackURL: "https://example.com/callback",
	}
	assert.Equal(t, "oak", req.Kind)
	assert.Equal(t, "north meadow", req.Location)
	assert.Equal(t, "https://example.com/callback", req.CallbackURL)

	tree := TreeWithID{
		Location: "north meadow",
		Kind:     "oak",
	}
	_ = tree.ID // UUID field exists

	result := TreePlantingResult{
		Success: true,
	}
	assert.True(t, result.Success)
	_ = result.ID // UUID field exists
}

// Verify callback request body type alias
func TestCallbackRequestBodyAlias(t *testing.T) {
	var body TreePlantedJSONRequestBody = TreePlantingResult{
		Success: true,
	}
	assert.True(t, body.Success)
}

// Verify ServerInterface is implementable (parent operation)
type testServer struct{}

func (s *testServer) PlantTree(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func TestServerInterfaceImplementable(t *testing.T) {
	var si ServerInterface = &testServer{}
	handler := Handler(si)
	require.NotNil(t, handler)
}

// Verify CallbackInitiator can be created
func TestCallbackInitiator(t *testing.T) {
	initiator, err := NewCallbackInitiator()
	require.NoError(t, err)
	require.NotNil(t, initiator)

	var _ CallbackInitiatorInterface = initiator
}

// Verify callback request builder works
func TestCallbackRequestBuilder(t *testing.T) {
	body := TreePlantingResult{
		Success: true,
	}

	req, err := NewTreePlantedCallbackRequest("https://example.com/callback", body)
	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, req.Method)
	assert.Equal(t, "https://example.com/callback", req.URL.String())
}

// Verify CallbackReceiverInterface is implementable
type testReceiver struct{}

func (r *testReceiver) HandleTreePlantedCallback(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestCallbackReceiverInterfaceImplementable(t *testing.T) {
	var ri CallbackReceiverInterface = &testReceiver{}
	handler := TreePlantedCallbackHandler(ri, func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, err.Error(), http.StatusBadRequest)
	})
	require.NotNil(t, handler)
}

// Verify ApplyDefaults methods exist
func TestApplyDefaults(t *testing.T) {
	req := &TreePlantingRequest{}
	req.ApplyDefaults()
	tree := &TreeWithID{}
	tree.ApplyDefaults()
	result := &TreePlantingResult{}
	result.ApplyDefaults()
}
