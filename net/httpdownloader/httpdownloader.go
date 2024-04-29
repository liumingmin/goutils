package httpdownloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/liumingmin/goutils/container"
	"github.com/liumingmin/goutils/log"
)

type HttpDownloader struct {
	HttpClient   *http.Client
	Headers      http.Header
	ConBlockChan chan struct{}
	BlockSize    int
	RetryCnt     int
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
	var err error

	length, err := t.getContentLength(ctx, url, header)
	if err != nil {
		return err
	}

	blockCount := length / int64(t.BlockSize)

	if length%int64(t.BlockSize) > 0 {
		blockCount += 1
	}

	file, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var mutex sync.Mutex
	var wg sync.WaitGroup

	for i := int64(0); i < blockCount; i++ {
		min := int64(t.BlockSize) * int64(i)   // Min range
		max := int64(t.BlockSize) * int64(i+1) // Max range

		if max > length {
			max = length
		}

		writer := NewOffsetWriter(file, &mutex, min, max)

		select {
		case <-ctx.Done():
			err = context.Canceled
			goto ExitFor
		case t.ConBlockChan <- struct{}{}:
		}

		wg.Add(1)

		go func(min int64, max int64) {
			defer func() {
				recover()
			}()

			defer func() {
				select {
				case <-t.ConBlockChan:
				default:
				}
			}()

			defer wg.Done()

			err = t.downloadBlock(ctx, url, header, min, max, writer)
		}(min, max)
	}
ExitFor:
	wg.Wait()

	return err
}

func (t *HttpDownloader) downloadBlock(ctx context.Context, url string, header http.Header, min, max int64, writer *HttpFileOffsetWriter) error {
	req, err := t.createRequest(ctx, "GET", url, header)
	if err != nil {
		return err
	}

	rangeHeader := fmt.Sprintf("bytes=%d-%d", min, max-1)
	req.Header.Set("Range", rangeHeader)

	buff := container.PoolBuffer4M.Get()
	defer container.PoolBuffer4M.Put(buff)

	for j := 0; j < t.RetryCnt; j++ {
		err := t.downloadToWriter(req, writer, buff)
		if err == nil {
			break
		}

		log.Error(ctx, "downloadBlock failed, retry download, url: %v, range: %v, err: %v, retry: %v", url, rangeHeader, err, j)
		writer.ResetOffset()
	}

	return err
}

func (t *HttpDownloader) downloadToWriter(req *http.Request, writer *HttpFileOffsetWriter, buff []byte) error {
	resp, err := t.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	_, err = io.CopyBuffer(writer, resp.Body, buff)
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
