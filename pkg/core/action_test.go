package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAction_ToString(t *testing.T) {
	tests := []struct {
		action   Action
		expected string
	}{
		{Allow, "Allow"},
		{Deny, "Deny"},
	}
	for _, tt := range tests {
		t.Run(tt.action.String(), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.action.String())
		})
	}

}

func TestAction_UnmarshalText(t *testing.T) {
	tests := []struct {
		s         string
		expected  Action
		expectErr bool
	}{
		{"Allow", Allow, false},
		{"Deny", Deny, false},
		{"allow", Allow, false},
		{"deny", Deny, false},
		{"ALLOW", Allow, false},
		{"DENY", Deny, false},
		{"someOtherString", Unknown, true},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			var a Action
			err := a.UnmarshalText([]byte(tt.s))

			assert.Equal(t, tt.expected, a)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
