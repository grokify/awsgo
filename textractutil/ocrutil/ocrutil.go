package ocrutil

import (
	"context"
	"errors"
	"image"
	"os"
	"time"
)

var ErrNoInputData = errors.New("no input data")

type OCRProcessor interface {
	OCRSync(ctx context.Context, b []byte) (TextResults, error)
	OCRSyncImageWriteFile(ctx context.Context, outfilename string, perm os.FileMode, imgs []image.Image) (TextResults, error)
	OCRSyncImageLocationsWriteFile(ctx context.Context, outfilename string, perm os.FileMode, imglocations []string) (TextResults, error)
}

type TextResults struct {
	OCRService   string
	OCRDateTime  time.Time
	Lines        []string
	WordsPrinted []string
}
