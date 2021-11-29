// Code generated by "make api"; DO NOT EDIT.
package sessions

type Connection struct {
	ClientTcpAddress   string `json:"client_tcp_address,omitempty"`
	ClientTcpPort      uint32 `json:"client_tcp_port,omitempty"`
	EndpointTcpAddress string `json:"endpoint_tcp_address,omitempty"`
	EndpointTcpPort    uint32 `json:"endpoint_tcp_port,omitempty"`
	BytesUp            uint64 `json:"bytes_up,omitempty"`
	BytesDown          uint64 `json:"bytes_down,omitempty"`
	ClosedReason       string `json:"closed_reason,omitempty"`
}