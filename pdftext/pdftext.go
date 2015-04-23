package pdftext

import (
	"bytes"
	"os"

	"regexp"

	"rsc.io/pdf"
)

const SpaceWidthPoints = 0.4
const SpaceWidthMaxPoints = SpaceWidthPoints * 20
const KernEpsilon = 1.45

type Font struct {
	Name string
	Size float64
}

func (f Font) Empty() bool { return f.Name == "" }

func FragFont(frag *pdf.Text) Font {
	return Font{Name: frag.Font, Size: frag.FontSize}
}

var rFiLigatureRepair = regexp.MustCompile(`!([a-z])`)
var rFtLigatureRepair = regexp.MustCompile(`([aeiou])!( |$)`)

// The "fi" ligature gets converted to "!" by rsc.io/pdf. This is an evil hack
func repairPDFText(text string) string {
	return rFtLigatureRepair.ReplaceAllString(
		rFiLigatureRepair.ReplaceAllString(text, "fi$1"),
		"${1}ft$2")
}

type Fragment struct {
	Font
	LX, RX, Y float64
	textBuf   *bytes.Buffer
}

func (t *Fragment) Text() string { return repairPDFText(t.textBuf.String()) }
func (t *Fragment) Empty() bool  { return t.Font.Empty() && t.textBuf.Len() == 0 }

func (t *Fragment) mergeable(frag *pdf.Text) bool {
	if t.Empty() {
		return true
	}
	fragmentAboveMe := frag.Y > t.Y // PDF Y axis has zero at bottom of page.
	if fragmentAboveMe || t.Font != FragFont(frag) {
		return false
	}
	if t.Y == frag.Y {
		xgap := frag.X - t.RX
		if xgap < 0-KernEpsilon {
			return false
		}
		return xgap < SpaceWidthMaxPoints
	}
	// Accept a merge if it's immediately one line below (or closer):
	return t.Y-frag.Y <= frag.FontSize
}

func (t *Fragment) hasSpaceBefore(frag *pdf.Text) bool {
	if t.Empty() {
		return false
	}
	if frag.Y != t.Y {
		return true
	}
	gap := frag.X - t.RX
	return gap >= SpaceWidthPoints && gap <= SpaceWidthMaxPoints
}

func (t *Fragment) merge(frag *pdf.Text) bool {
	if t.mergeable(frag) {
		if t.Empty() {
			t.LX = frag.X
		}
		if t.hasSpaceBefore(frag) {
			t.textBuf.WriteByte(' ')
		}
		t.Font = FragFont(frag)
		t.Y = frag.Y
		t.RX = frag.X + frag.W
		t.textBuf.WriteString(frag.S)
		return true
	}
	return false
}

type Scanner struct {
	page, textIndex int
	currentPage     pdf.Content
	Pages           []pdf.Content
}

// NewFileScanner creates a new PDF scanner from a file. Note that the entire
// PDF is read into memory and the file is immediately closed.
func NewFileScanner(pdffile string) (*Scanner, error) {
	pdfh, err := os.Open(pdffile)
	if err != nil {
		return nil, err
	}
	defer pdfh.Close()

	stat, err := pdfh.Stat()
	if err != nil {
		return nil, err
	}

	r, err := pdf.NewReader(pdfh, stat.Size())
	if err != nil {
		return nil, err
	}
	return NewScanner(r), nil
}

// NewScanner creates a new PDF scanner given a PDF reader. Note that the entire
// PDF is pre-fetched and retained in memory.
func NewScanner(pdfr *pdf.Reader) *Scanner {
	scanner := &Scanner{}
	pageCount := pdfr.NumPage()
	scanner.Pages = make([]pdf.Content, pageCount)
	for i := 1; i <= pageCount; i++ {
		scanner.Pages[i-1] = pdfr.Page(i).Content()
	}
	if pageCount > 0 {
		scanner.currentPage = scanner.Pages[0]
	}
	return scanner
}

func emptyContent(c pdf.Content) bool { return len(c.Text) == 0 && len(c.Rect) == 0 }

func (r *Scanner) nextFragment() *pdf.Text {
	for emptyContent(r.currentPage) || r.textIndex >= len(r.currentPage.Text) {
		r.page++
		if r.page >= len(r.Pages) {
			return nil
		}
		r.textIndex = 0
		r.currentPage = r.Pages[r.page]
	}
	nextFrag := &r.currentPage.Text[r.textIndex]
	r.textIndex++
	return nextFrag
}

func (r *Scanner) rewind() {
	r.textIndex--
	if r.textIndex < 0 {
		r.currentPage = pdf.Content{}
		if r.page > 0 {
			r.page--
			r.currentPage = r.Pages[r.page]
			r.textIndex = len(r.currentPage.Text) - 1
		}
		if r.textIndex < 0 {
			r.textIndex = 0
		}
	}
}

// NextText returns the next continuous fragment of text from the underlying PDF
// data, or nil if all fragments have been read.
func (r *Scanner) NextText() *Fragment {
	var pending = &Fragment{textBuf: &bytes.Buffer{}}
	result := func() *Fragment {
		if pending.Empty() {
			return nil
		}
		return pending
	}
	for {
		frag := r.nextFragment()
		if frag == nil {
			return result()
		}
		if !pending.merge(frag) {
			r.rewind()
			return result()
		}
	}
}
