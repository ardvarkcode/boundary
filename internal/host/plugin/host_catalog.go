// Package plugin provides a plugin host catalog, and plugin host set resource
// which are used to interact with a host plugin as well as a repository to
// perform CRUDL and custom actions on these resource types.
package plugin

import (
	"context"

	"github.com/hashicorp/boundary/internal/errors"
	"github.com/hashicorp/boundary/internal/host/plugin/store"
	"github.com/hashicorp/boundary/internal/oplog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// A HostCatalog contains plugin host sets. It is owned by
// a scope.
type HostCatalog struct {
	*store.HostCatalog
	tableName string `gorm:"-"`

	Secrets *structpb.Struct `gorm:"-"`
}

// NewHostCatalog creates a new in memory HostCatalog assigned to a scopeId
// and pluginId. Name and description are the only valid options. All other
// options are ignored.
func NewHostCatalog(ctx context.Context, scopeId, pluginId string, opt ...Option) (*HostCatalog, error) {
	const op = "plugin.NewHostCatalog"
	opts := getOpts(opt...)

	attrs, err := proto.Marshal(opts.withAttributes)
	if err != nil {
		return nil, errors.Wrap(ctx, err, op, errors.WithCode(errors.InvalidParameter))
	}

	hc := &HostCatalog{
		HostCatalog: &store.HostCatalog{
			ScopeId:     scopeId,
			PluginId:    pluginId,
			Name:        opts.withName,
			Description: opts.withDescription,
			Attributes:  attrs,
		},
		Secrets: opts.withSecrets,
	}
	return hc, nil
}

func allocHostCatalog() *HostCatalog {
	return &HostCatalog{
		HostCatalog: &store.HostCatalog{},
	}
}

// clone provides a deep copy of the HostCatalog with the exception of the
// secret.  The secret shallow copied.
func (c *HostCatalog) clone() *HostCatalog {
	cp := proto.Clone(c.HostCatalog)
	newSecret := proto.Clone(c.Secrets)

	hc := &HostCatalog{
		HostCatalog: cp.(*store.HostCatalog),
		Secrets:     newSecret.(*structpb.Struct),
	}
	// proto.Clone will convert slices with length and capacity of 0 to nil.
	// Fix this since gorm treats empty slices differently than nil.
	if c.Attributes != nil && len(c.Attributes) == 0 && hc.Attributes == nil {
		hc.Attributes = []byte{}
	}
	return hc
}

// TableName returns the table name for the host catalog.
func (c *HostCatalog) TableName() string {
	if c.tableName != "" {
		return c.tableName
	}
	return "host_plugin_catalog"
}

// SetTableName sets the table name. If the caller attempts to
// set the name to "" the name will be reset to the default name.
func (c *HostCatalog) SetTableName(n string) {
	c.tableName = n
}

func (s *HostCatalog) oplog(op oplog.OpType) oplog.Metadata {
	metadata := oplog.Metadata{
		"resource-public-id": []string{s.PublicId},
		"resource-type":      []string{"plugin-host-catalog"},
		"op-type":            []string{op.String()},
	}
	if s.ScopeId != "" {
		metadata["scope-id"] = []string{s.ScopeId}
	}
	return metadata
}