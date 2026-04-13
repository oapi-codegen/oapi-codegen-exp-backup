package templates

// Import represents a Go import with optional alias.
type Import struct {
	Path  string
	Alias string // empty if no alias
}

// ServerTemplate defines a template for server generation.
type ServerTemplate struct {
	Name     string   // Template name (e.g., "interface", "handler")
	Imports  []Import // Required imports for this template
	Template string   // Template path in embedded FS
}

// ReceiverTemplate defines a template for receiver (webhook/callback) generation.
type ReceiverTemplate struct {
	Name     string   // Template name (e.g., "receiver")
	Imports  []Import // Required imports for this template
	Template string   // Template path in embedded FS
}

// StdHTTPReceiverTemplates contains receiver templates for StdHTTP servers.
var StdHTTPReceiverTemplates = map[string]ReceiverTemplate{
	"receiver": {
		Name: "receiver",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
		},
		Template: "server/stdhttp/receiver.go.tmpl",
	},
}

// ChiReceiverTemplates contains receiver templates for Chi servers.
var ChiReceiverTemplates = map[string]ReceiverTemplate{
	"receiver": {
		Name: "receiver",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
		},
		Template: "server/chi/receiver.go.tmpl",
	},
}

// EchoReceiverTemplates contains receiver templates for Echo v5 servers.
var EchoReceiverTemplates = map[string]ReceiverTemplate{
	"receiver": {
		Name: "receiver",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
			{Path: "github.com/labstack/echo/v5"},
		},
		Template: "server/echo/receiver.go.tmpl",
	},
}

// EchoV4ReceiverTemplates contains receiver templates for Echo v4 servers.
var EchoV4ReceiverTemplates = map[string]ReceiverTemplate{
	"receiver": {
		Name: "receiver",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
			{Path: "github.com/labstack/echo/v4"},
		},
		Template: "server/echo-v4/receiver.go.tmpl",
	},
}

// GinReceiverTemplates contains receiver templates for Gin servers.
var GinReceiverTemplates = map[string]ReceiverTemplate{
	"receiver": {
		Name: "receiver",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
			{Path: "github.com/gin-gonic/gin"},
		},
		Template: "server/gin/receiver.go.tmpl",
	},
}

// GorillaReceiverTemplates contains receiver templates for Gorilla servers.
var GorillaReceiverTemplates = map[string]ReceiverTemplate{
	"receiver": {
		Name: "receiver",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
		},
		Template: "server/gorilla/receiver.go.tmpl",
	},
}

// FiberReceiverTemplates contains receiver templates for Fiber servers.
var FiberReceiverTemplates = map[string]ReceiverTemplate{
	"receiver": {
		Name: "receiver",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
			{Path: "github.com/gofiber/fiber/v3"},
		},
		Template: "server/fiber/receiver.go.tmpl",
	},
}

// IrisReceiverTemplates contains receiver templates for Iris servers.
var IrisReceiverTemplates = map[string]ReceiverTemplate{
	"receiver": {
		Name: "receiver",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
			{Path: "github.com/kataras/iris/v12"},
		},
		Template: "server/iris/receiver.go.tmpl",
	},
}

// StdHTTPServerTemplates contains templates for StdHTTP server generation.
var StdHTTPServerTemplates = map[string]ServerTemplate{
	"interface": {
		Name: "interface",
		Imports: []Import{
			{Path: "net/http"},
		},
		Template: "server/stdhttp/interface.go.tmpl",
	},
	"handler": {
		Name: "handler",
		Imports: []Import{
			{Path: "net/http"},
		},
		Template: "server/stdhttp/handler.go.tmpl",
	},
	"wrapper": {
		Name: "wrapper",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
		},
		Template: "server/stdhttp/wrapper.go.tmpl",
	},
}

// ChiServerTemplates contains templates for Chi server generation.
var ChiServerTemplates = map[string]ServerTemplate{
	"interface": {
		Name: "interface",
		Imports: []Import{
			{Path: "net/http"},
		},
		Template: "server/chi/interface.go.tmpl",
	},
	"handler": {
		Name: "handler",
		Imports: []Import{
			{Path: "net/http"},
			{Path: "github.com/go-chi/chi/v5"},
		},
		Template: "server/chi/handler.go.tmpl",
	},
	"wrapper": {
		Name: "wrapper",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
			{Path: "github.com/go-chi/chi/v5"},
		},
		Template: "server/chi/wrapper.go.tmpl",
	},
}

// EchoServerTemplates contains templates for Echo v5 server generation.
var EchoServerTemplates = map[string]ServerTemplate{
	"interface": {
		Name: "interface",
		Imports: []Import{
			{Path: "net/http"},
			{Path: "github.com/labstack/echo/v5"},
		},
		Template: "server/echo/interface.go.tmpl",
	},
	"handler": {
		Name: "handler",
		Imports: []Import{
			{Path: "github.com/labstack/echo/v5"},
		},
		Template: "server/echo/handler.go.tmpl",
	},
	"wrapper": {
		Name: "wrapper",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
			{Path: "github.com/labstack/echo/v5"},
		},
		Template: "server/echo/wrapper.go.tmpl",
	},
}

// EchoV4ServerTemplates contains templates for Echo v4 server generation.
var EchoV4ServerTemplates = map[string]ServerTemplate{
	"interface": {
		Name: "interface",
		Imports: []Import{
			{Path: "net/http"},
			{Path: "github.com/labstack/echo/v4"},
		},
		Template: "server/echo-v4/interface.go.tmpl",
	},
	"handler": {
		Name: "handler",
		Imports: []Import{
			{Path: "github.com/labstack/echo/v4"},
		},
		Template: "server/echo-v4/handler.go.tmpl",
	},
	"wrapper": {
		Name: "wrapper",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
			{Path: "github.com/labstack/echo/v4"},
		},
		Template: "server/echo-v4/wrapper.go.tmpl",
	},
}

// GinServerTemplates contains templates for Gin server generation.
var GinServerTemplates = map[string]ServerTemplate{
	"interface": {
		Name: "interface",
		Imports: []Import{
			{Path: "net/http"},
			{Path: "github.com/gin-gonic/gin"},
		},
		Template: "server/gin/interface.go.tmpl",
	},
	"handler": {
		Name: "handler",
		Imports: []Import{
			{Path: "github.com/gin-gonic/gin"},
		},
		Template: "server/gin/handler.go.tmpl",
	},
	"wrapper": {
		Name: "wrapper",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
			{Path: "github.com/gin-gonic/gin"},
		},
		Template: "server/gin/wrapper.go.tmpl",
	},
}

// GorillaServerTemplates contains templates for Gorilla server generation.
var GorillaServerTemplates = map[string]ServerTemplate{
	"interface": {
		Name: "interface",
		Imports: []Import{
			{Path: "net/http"},
		},
		Template: "server/gorilla/interface.go.tmpl",
	},
	"handler": {
		Name: "handler",
		Imports: []Import{
			{Path: "net/http"},
			{Path: "github.com/gorilla/mux"},
		},
		Template: "server/gorilla/handler.go.tmpl",
	},
	"wrapper": {
		Name: "wrapper",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
			{Path: "github.com/gorilla/mux"},
		},
		Template: "server/gorilla/wrapper.go.tmpl",
	},
}

// FiberServerTemplates contains templates for Fiber server generation.
var FiberServerTemplates = map[string]ServerTemplate{
	"interface": {
		Name: "interface",
		Imports: []Import{
			{Path: "github.com/gofiber/fiber/v3"},
		},
		Template: "server/fiber/interface.go.tmpl",
	},
	"handler": {
		Name: "handler",
		Imports: []Import{
			{Path: "github.com/gofiber/fiber/v3"},
		},
		Template: "server/fiber/handler.go.tmpl",
	},
	"wrapper": {
		Name: "wrapper",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
			{Path: "github.com/gofiber/fiber/v3"},
		},
		Template: "server/fiber/wrapper.go.tmpl",
	},
}

// IrisServerTemplates contains templates for Iris server generation.
var IrisServerTemplates = map[string]ServerTemplate{
	"interface": {
		Name: "interface",
		Imports: []Import{
			{Path: "net/http"},
			{Path: "github.com/kataras/iris/v12"},
		},
		Template: "server/iris/interface.go.tmpl",
	},
	"handler": {
		Name: "handler",
		Imports: []Import{
			{Path: "github.com/kataras/iris/v12"},
		},
		Template: "server/iris/handler.go.tmpl",
	},
	"wrapper": {
		Name: "wrapper",
		Imports: []Import{
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "net/http"},
			{Path: "net/url"},
			{Path: "github.com/kataras/iris/v12"},
		},
		Template: "server/iris/wrapper.go.tmpl",
	},
}

// SharedServerTemplates contains templates shared across all server implementations.
var SharedServerTemplates = map[string]ServerTemplate{
	"errors": {
		Name: "errors",
		Imports: []Import{
			{Path: "fmt"},
		},
		Template: "server/errors.go.tmpl",
	},
	"param_types": {
		Name: "param_types",
		Imports: []Import{},
		Template: "server/param_types.go.tmpl",
	},
}

// InitiatorTemplate defines a template for initiator (webhook/callback sender) generation.
type InitiatorTemplate struct {
	Name     string   // Template name (e.g., "initiator_base")
	Imports  []Import // Required imports for this template
	Template string   // Template path in embedded FS
}

// InitiatorTemplates contains the base template for initiator generation.
var InitiatorTemplates = map[string]InitiatorTemplate{
	"initiator_base": {
		Name: "initiator_base",
		Imports: []Import{
			{Path: "context"},
			{Path: "net/http"},
		},
		Template: "initiator/base.go.tmpl",
	},
}

// ClientTemplate defines a template for client generation.
type ClientTemplate struct {
	Name     string   // Template name (e.g., "base")
	Imports  []Import // Required imports for this template
	Template string   // Template path in embedded FS
}

// ClientTemplates contains the base template for client generation.
var ClientTemplates = map[string]ClientTemplate{
	"base": {
		Name: "base",
		Imports: []Import{
			{Path: "context"},
			{Path: "net/http"},
			{Path: "net/url"},
			{Path: "strings"},
		},
		Template: "client/base.go.tmpl",
	},
}

// SenderTemplate defines a template shared between client and initiator generation.
type SenderTemplate struct {
	Name     string   // Template name (e.g., "sender_interface")
	Imports  []Import // Required imports for this template
	Template string   // Template path in embedded FS
}

// SenderTemplates contains templates shared between client and initiator generators.
var SenderTemplates = map[string]SenderTemplate{
	"sender_interface": {
		Name: "sender_interface",
		Imports: []Import{
			{Path: "context"},
			{Path: "io"},
			{Path: "net/http"},
		},
		Template: "sender/interface.go.tmpl",
	},
	"sender_methods": {
		Name: "sender_methods",
		Imports: []Import{
			{Path: "context"},
			{Path: "io"},
			{Path: "net/http"},
		},
		Template: "sender/methods.go.tmpl",
	},
	"sender_request_builders": {
		Name: "sender_request_builders",
		Imports: []Import{
			{Path: "bytes"},
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "io"},
			{Path: "net/http"},
			{Path: "net/url"},
			{Path: "strings"},
		},
		Template: "sender/request_builders.go.tmpl",
	},
	"sender_simple": {
		Name: "sender_simple",
		Imports: []Import{
			{Path: "context"},
			{Path: "encoding/json"},
			{Path: "fmt"},
			{Path: "io"},
			{Path: "net/http"},
		},
		Template: "sender/simple.go.tmpl",
	},
}
