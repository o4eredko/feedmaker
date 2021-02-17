package broadcaster_test

import (
	"time"

	"go-feedmaker/entity"
)

var (
	defaultMessage    = []byte("default message")
	defaultGeneration = &entity.Generation{
		ID:            "0xDEADBEEF",
		Type:          "degeneration",
		Progress:      42,
		DataFetched:   true,
		FilesUploaded: 13,
		StartTime:     time.Now().UTC(),
		EndTime:       time.Time{},
	}
)
