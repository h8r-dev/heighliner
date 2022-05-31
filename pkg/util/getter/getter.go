package getter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/cavaliergopher/grab/v3"
)

// Request of getting resources
type Request struct {
	Src  string
	Dst  string
	Name string
}

// NewRequest returns a request
func NewRequest(src, dst, name string) *Request {
	return &Request{
		Src:  src,
		Dst:  dst,
		Name: name,
	}
}

// Get the resources
func Get(w io.Writer, req *Request) error {
	client := grab.NewClient()

	src := req.Src
	dst := req.Dst
	name := req.Name
	r, err := grab.NewRequest(filepath.Join(dst, name), src)
	if err != nil {
		return fmt.Errorf("bad request: %w", err)
	}

	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	resp := client.Do(r)

	t := time.NewTicker(3 * time.Second)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Fprintf(w, "  downloaded %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size(),
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	return nil
}
