package util

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/cheggaaa/pb"
	"github.com/hashicorp/go-getter/v2"
	"github.com/rs/zerolog/log"
)

// GetWithTracker gets resources and tracks the progress
func GetWithTracker(req *getter.Request) error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working dir: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	client := &getter.Client{}

	req.Pwd = pwd
	req.ProgressListener = defaultProgressBar

	wg := sync.WaitGroup{}
	wg.Add(1)
	errChan := make(chan error, 2)
	go func() {
		defer wg.Done()
		defer cancel()
		if _, err := client.Get(ctx, req); err != nil {
			errChan <- err
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case sig := <-c:
		signal.Reset(os.Interrupt)
		cancel()
		log.Info().Msgf("signal %s", sig)
		return nil
	case <-ctx.Done():
		wg.Wait()
		return nil
	case err := <-errChan:
		wg.Wait()
		return err
	}
}

// defaultProgressBar is the default instance of a cheggaaa
// progress bar.
var defaultProgressBar getter.ProgressTracker = &ProgressBar{}

// ProgressBar wraps a github.com/cheggaaa/pb.Pool
// in order to display download progress for one or multiple
// downloads.
//
// If two different instance of ProgressBar try to
// display a progress only one will be displayed.
// It is therefore recommended to use DefaultProgressBar
type ProgressBar struct {
	// lock everything below
	lock sync.Mutex

	pool *pb.Pool

	pbs int
}

// ProgressBarConfig configures progressBar
func ProgressBarConfig(bar *pb.ProgressBar, prefix string) {
	bar.SetUnits(pb.U_BYTES)
	bar.Prefix(prefix)
}

// TrackProgress instantiates a new progress bar that will
// display the progress of stream until closed.
// total can be 0.
func (cpb *ProgressBar) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser {
	cpb.lock.Lock()
	defer cpb.lock.Unlock()

	newPb := pb.New64(totalSize)
	newPb.Set64(currentSize)
	ProgressBarConfig(newPb, filepath.Base(src))
	if cpb.pool == nil {
		cpb.pool = pb.NewPool()
		_ = cpb.pool.Start()
	}
	cpb.pool.Add(newPb)
	reader := newPb.NewProxyReader(stream)

	cpb.pbs++
	return &readCloser{
		Reader: reader,
		close: func() error {
			cpb.lock.Lock()
			defer cpb.lock.Unlock()

			newPb.Finish()
			cpb.pbs--
			if cpb.pbs <= 0 {
				_ = cpb.pool.Stop()
				cpb.pool = nil
			}
			return nil
		},
	}
}

type readCloser struct {
	io.Reader
	close func() error
}

func (c *readCloser) Close() error { return c.close() }
