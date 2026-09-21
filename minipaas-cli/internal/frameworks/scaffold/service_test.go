package scaffold

import "testing"

func TestNamespaceNetworks(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		namespace string
		want      string
	}{
		{
			name:      "default namespace keeps names",
			content:   "networks:\n  minipaas_network:\n  minipaas_public:\n",
			namespace: "minipaas",
			want:      "networks:\n  minipaas_network:\n  minipaas_public:\n",
		},
		{
			name:      "custom namespace prefixes networks",
			content:   "networks:\n  minipaas_network:\n  minipaas_public:\n",
			namespace: "acme",
			want:      "networks:\n  acme_network:\n  acme_public:\n",
		},
	}
	for _, tt := range tests {
		if got := namespaceNetworks(tt.content, tt.namespace); got != tt.want {
			t.Fatalf("%s: got %q, want %q", tt.name, got, tt.want)
		}
	}
}
