package session

import (
	"context"
	"github.com/hashicorp/boundary/internal/db"
	"github.com/hashicorp/boundary/internal/errors"
	"github.com/hashicorp/boundary/internal/kms"
)

// Repository is the session database repository
type ConnectionRepository struct {
	reader db.Reader
	writer db.Writer
	kms    *kms.Kms

	// defaultLimit provides a default for limiting the number of results returned from the repo
	defaultLimit int
}

// NewRepository creates a new session Connection Repository. Supports the options: WithLimit
// which sets a default limit on results returned by repo operations.
func NewConnectionRepository(r db.Reader, w db.Writer, kms *kms.Kms, opt ...Option) (*ConnectionRepository, error) {
	const op = "sessionConnection.NewRepository"
	if r == nil {
		return nil, errors.NewDeprecated(errors.InvalidParameter, op, "nil reader")
	}
	if w == nil {
		return nil, errors.NewDeprecated(errors.InvalidParameter, op, "nil writer")
	}
	if kms == nil {
		return nil, errors.NewDeprecated(errors.InvalidParameter, op, "nil kms")
	}
	opts := getOpts(opt...)
	if opts.withLimit == 0 {
		// zero signals the boundary defaults should be used.
		opts.withLimit = db.DefaultLimit
	}
	return &ConnectionRepository{
		reader:       r,
		writer:       w,
		kms:          kms,
		defaultLimit: opts.withLimit,
	}, nil
}

// list will return a listing of resources and honor the WithLimit option or the
// repo defaultLimit.  Supports WithOrder option.
func (r *ConnectionRepository) list(ctx context.Context, resources interface{}, where string, args []interface{}, opt ...Option) error {
	const op = "session.(ConnectionRepository).list"
	opts := getOpts(opt...)
	limit := r.defaultLimit
	var dbOpts []db.Option
	if opts.withLimit != 0 {
		// non-zero signals an override of the default limit for the repo.
		limit = opts.withLimit
	}
	dbOpts = append(dbOpts, db.WithLimit(limit))
	switch opts.withOrderByCreateTime {
	case db.AscendingOrderBy:
		dbOpts = append(dbOpts, db.WithOrder("create_time asc"))
	case db.DescendingOrderBy:
		dbOpts = append(dbOpts, db.WithOrder("create_time"))
	}
	if err := r.reader.SearchWhere(ctx, resources, where, args, dbOpts...); err != nil {
		return errors.Wrap(ctx, err, op)
	}
	return nil
}
