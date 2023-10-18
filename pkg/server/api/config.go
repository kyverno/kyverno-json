package api

import (
	"github.com/kyverno/kyverno-json/pkg/server/api/scan"
)

type Configuration struct {
	PolicyProvider scan.PolicyProvider
}
