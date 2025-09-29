package metadata

import "context"

type MetadataSource interface {
	GetMetadata(ctx context.Context) (Metadata, error)
}
