package graph

import (
	"github.com/erknas/forum/internal/service"
	"github.com/erknas/forum/internal/subscription"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Svc service.Servicer
	Sub subscription.Subscriber
}
