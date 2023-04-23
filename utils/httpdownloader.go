package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/liumingmin/goutils/log"
)

type HttpDownloader struct {
	HttpClient    *http.Client
	Headers       http.Header
	GoroutinesCnt int
	RetryCnt      int
}

func (t *HttpDownloader) getContentLength(ctx context.Context, url string, header http.Header) (int64, error) {
	req, err := t.createRequest(ctx, "HEAD", url, header)
	if err != nil {
		return 0, err
	}

	res, err := t.HttpClient.Do(req)
	if err != nil {
		return 0, err
	}

	headerValue := res.Header["Content-Length"]
	if len(headerValue) > 0 {
		return strconv.ParseInt(headerValue[0], 10, 64)
	}

	return 0, nil
}

func (t *HttpDownloader) Download(ctx context.Context, url string, header http.Header, savePath string) error {
	var wg sync.WaitGroup

	length, err := t.getContentLength(ctx, url, header)
	if err != nil {
		return err
	}

	goroutinesCnt := t.GoroutinesCnt
	if length < int64(goroutinesCnt) {
		goroutinesCnt = int(length)
	}

	lenSub := length / int64(goroutinesCnt) // Bytes for each Go-routine
	diff := length % int64(goroutinesCnt)   // Get the remaining for the last request

	file, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var mutex sync.Mutex

	// Make up a temporary array to hold the data to be written to the file
	for i := 0; i < goroutinesCnt; i++ {
		wg.Add(1)

		min := lenSub * int64(i)   // Min range
		max := lenSub * int64(i+1) // Max range

		if i == goroutinesCnt-1 {
			max += diff // Add the remaining bytes in the last request
		}

		writer := NewOffsetWriter(file, &mutex, min, max)

		go func(min int64, max int64) {
			defer func() {
				recover()
			}()
			defer wg.Done()

			req, err := t.createRequest(ctx, "GET", url, header)
			if err != nil {
				return
			}

			rangeHeader := fmt.Sprintf("bytes=%d-%d", min, max-1) // Add the data for the Range header of the form "bytes=0-100"
			req.Header.Add("Range", rangeHeader)

			for j := 0; j < t.RetryCnt; j++ {
				err := t.downloadBlock(req, writer)
				if err == nil {
					break
				}

				log.Error(ctx, "downloadBlock failed, retry download, url: %v, range: %v, err: %v, retry: %v", url, rangeHeader, err, j)
				writer.ResetOffset()
			}
		}(min, max)
	}
	wg.Wait()

	return nil
}

func (t *HttpDownloader) downloadBlock(req *http.Request, writer *FileOffsetWriter) error {
	resp, err := t.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	_, err = io.Copy(writer, resp.Body)
	return err
}

func (t *HttpDownloader) createRequest(ctx context.Context, method, url string, header http.Header) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	reqHeader := http.Header{}

	if t.Headers != nil {
		for key, value := range t.Headers {
			value2 := make([]string, len(value))
			copy(value2, value)
			reqHeader[key] = value2
		}
	}

	if header != nil {
		for key, value := range header {
			value2 := make([]string, len(value))
			copy(value2, value)
			reqHeader[key] = value2
		}
	}

	req.Header = reqHeader
	return req, nil
}
