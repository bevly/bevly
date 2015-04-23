package httpagent

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"github.com/PuerkitoBio/goquery"
)

const DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36"

type Agent struct {
	UserAgent     string
	Client        http.Client
	ForceEncoding string
}

func Win1252Agent() *Agent {
	agent := New()
	agent.ForceEncoding = "windows-1252"
	return agent
}

func New() *Agent {
	return &Agent{
		UserAgent: DefaultUserAgent,
		Client:    http.Client{Timeout: time.Second * 30},
	}
}

func (h *Agent) Get(requrl string) (*http.Response, error) {
	parsedUrl, err := url.Parse(requrl)
	if err != nil {
		return nil, err
	}
	req := http.Request{
		URL: parsedUrl,
		Header: http.Header{
			"User-Agent": {h.UserAgent},
		},
	}
	return h.Translate(h.Client.Do(&req))
}

// GetFile downloads requrl to destFile, creating any parent directories of
// destFile if necessary. On successful download, returns the number of bytes
// written to destFile.
func (h *Agent) GetFile(requrl, destFile string) (int64, error) {
	res, err := h.Get(requrl)
	if err != nil {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
		return 0, err
	}
	defer res.Body.Close()

	parentDir := filepath.Dir(destFile)
	if err = os.MkdirAll(parentDir, os.ModeDir|0755); err != nil {
		return 0, err
	}

	file, err := os.Create(destFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return io.Copy(file, res.Body)
}

func (h *Agent) GetDoc(requrl string) (*goquery.Document, error) {
	res, err := h.Get(requrl)
	if err != nil {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
		return nil, err
	}
	return goquery.NewDocumentFromResponse(res)
}

func (h *Agent) Translate(res *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return res, err
	}

	if h.ForceEncoding != "" {
		res.Body, err = h.reencodeBody(res.Body, h.ForceEncoding)
		if err != nil {
			return res, err
		}
	}

	if res.StatusCode >= 400 {
		err = errors.New(fmt.Sprintf("http request error: %s", res.StatusCode))
		return res, err
	}

	return res, nil
}

type TranslatingReader struct {
	io.ReadCloser
	reader io.Reader
}

func (t TranslatingReader) Read(p []byte) (int, error) {
	return t.reader.Read(p)
}

func (h *Agent) reencodeBody(body io.ReadCloser, enc string) (io.ReadCloser, error) {
	translator, err := charset.TranslatorFrom(enc)
	if err != nil {
		return body, err
	}
	return TranslatingReader{
		ReadCloser: body,
		reader:     charset.NewTranslatingReader(body, translator),
	}, nil
}
