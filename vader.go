package vader

import (
	"bufio"
	"github.com/iseurie/vader-go/absolutes"
	"github.com/jdkato/prose"
	"golang.org/x/text/transform"
	"io"
	"math"
	"strings"
	"unicode"
)

func sgn(x float64) float64 {
	btf := map[bool]float64{true: 1, false: -1}
	return b2f[math.Signbit(x)]
}

func strp(raw string) string {
	strip :=
		transform.Chain(
			norm.NFD,
			runes.Remove(unicode.IsPunct),
			runes.Remove(unicode.In(unicode.Mn)),
			norm.NFC)
	return transform.String(strip, raw)
}

// Tokens represents a tokenized document.
type Tokens []prose.Token

// OfRaw maps a raw string onto a tokenized document.
func (T Tokens) OfRaw(s string, tag bool) {
	docopt := prose.DocOpt{
		UsingModel:       false,
		WithSegmentation: false,
		WithTagging:      tag,
		WithTokenization: true,
	}
	T = prose.NewDocument(s).Tokens()
}

func (T Tokens) Strip() {
	rtn := make([]string, len(T), len(T))
	for i, t := range T {
		T[i].Text = strp(t.Text)
	}
}

func (T Tokens) Stripped() string {
	var S strings.Builder
	for _, t := range T {
		S.WriteByte(' ')
		S.WriteString(strp(t.Text))
	}
	return S.String()
}

// Negated returns whether the receiver contains any negations.
func (T Tokens) Negated() bool {
	for i, t := range T {
		k := strings.ToLower(t.Text)
		switch {
		case absolutes.Negate(k):
		case strings.Contains("n't", k):
		case k == "least" && (i == 0 || T[i-1] == "at"):
			return true
		default:
			return false
		}
	}
}

// ClDiff (Caps-lock differential) checks whether some words in the input are ALL CAPS.
func (T Tokens) ClDiff() bool {
	allcaps := 0
	for _, t := range T {
		if strings.ToUpper(t.Text) == t.Text {
			allcaps++
		}
	}
	caseDifferential := len(T) - allcaps
	return 0 < caseDifferential < len(T)
}

// Valence represents a sentiment score.
type Valence float64

// Norm norms a valence score, providing a sensible point of entry from the
// API surface to the client code.
func (v Valence) Norm() float64 {
	s := float64(v)
	norm := v / math.Sqrt(v*v+15)
	norm = math.Min(norm, 1)
	norm = math.Max(norm, -1)
	return norm
}

// Chain negates or affirms the valence based on the last word.
func (v Valence) SclrSgn(T Tokens, i int) float64 {
	// check the capitalization differential
	capdiff := T.ClDiff()
	scalar := absolutes.Boost(t)
	scalar *= sgn(v)
	if strings.ToUpper(t) == t && capdiff {
		scalar += sgn(v) * absolutes.CFactor
	}
	return scalar
}

func (v Valence) ButCk(prev 

// TODO: Implement actual heuristics and checks for sentiment analysis
func (v Valence) Gauge(doc Tokens) {

}
