// Copyright (c) 2017 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package virtcontainers

import (
	"fmt"
)

// This is the no proxy implementation of the proxy interface. This
// is a generic implementation for any case (basically any agent),
// where no actual proxy is needed. This happens when the combination
// of the VM and the agent can handle multiple connections without
// additional component to handle the multiplexing. Both the runtime
// and the shim will connect to the agent through the VM, bypassing
// the proxy model.
// That's why this implementation is very generic, and all it does
// is to provide both shim and runtime the correct URL to connect
// directly to the VM.
type noProxy struct {
}

// start is noProxy start implementation for proxy interface.
func (p *noProxy) start(params proxyParams) (int, string, error) {
	if params.logger == nil {
		return -1, "", fmt.Errorf("proxy logger is not set")
	}

	params.logger.Debug("No proxy started because of no-proxy implementation")

	if params.agentURL == "" {
		return -1, "", fmt.Errorf("AgentURL cannot be empty")
	}

	return params.hid, params.agentURL, nil
}

// stop is noProxy stop implementation for proxy interface.
func (p *noProxy) stop(pid int) error {
	return nil
}

// The noproxy doesn't need to watch the vm console, thus return false always.
func (p *noProxy) consoleWatched() bool {
	return false
}
