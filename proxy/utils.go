package proxy

import "io"

func transfer(dest io.WriteCloser, src io.ReadCloser) {
	defer func() {
		if dest != nil {
			dest.Close()
		}
	}()
	defer func() {
		if src != nil {
			src.Close()
		}
	}()
	if dest != nil && src != nil {
		buf := make([]byte, 1, 1)
		io.CopyBuffer(dest, src, buf)
	}
}
