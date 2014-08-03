package httpagent

import (
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"time"
)

const DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36"

type HttpAgent struct {
	UserAgent     string
	Client        http.Client
	ForceEncoding string
}

func Agent() *HttpAgent {
	return &HttpAgent{
		UserAgent: DefaultUserAgent,
		Client:    http.Client{Timeout: time.Second * 30},
	}
}

func (h *HttpAgent) Get(requrl string) (*http.Response, error) {
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

func (h *HttpAgent) GetDoc(requrl string) (*goquery.Document, error) {
	res, err := h.Get(requrl)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromResponse(res)
}

func (h *HttpAgent) Translate(res *http.Response, err error) (*http.Response, error) {
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

func (h *HttpAgent) reencodeBody(body io.ReadCloser, enc string) (io.ReadCloser, error) {
	translator, err := charset.TranslatorFrom(enc)
	if err != nil {
		return body, err
	}
	return TranslatingReader{
		ReadCloser: body,
		reader:     charset.NewTranslatingReader(body, translator),
	}, nil
}
