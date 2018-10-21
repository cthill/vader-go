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

var b2f = map[bool]float64{true: 1, false: -1}

func sgn(x float64) float64 {
	return b2f[math.Signbit(x)]
}

type Tokens []prose.Token

func (T Tokens) OfRaw(s string) {
	docopt := prose.DocOpt{
		UsingModel:       false,
		WithSegmentation: false,
		WithTagging:      true,
		WithTokenization: true,
	}
	T = prose.NewDocument(s).Tokens()
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

// SentiText is an intermediate struct for evaluating sentiment of a given document.
type SentiText struct {
	// The cartesian product of absolutes.Punctuation the set of tokens in
	// the list, permuting each token by prepending and appending punctuation.
	w2wp map[string]string

	// Remove leading and trailing punctuation, but leave contractions and emoticons
	wes string

	// Keep the raw input
	raw Tokens

	// Lexicon is the structure's lexicon
	Lexicon

	capDiff bool
}

// NewSentiText constructs a new sentiment evaluation model from the given
// parameters.
func NewSentiText(raw string, lexicon Lexicon) (rtn SentiText) {
	rtn.raw.OfRaw(raw)
	rtn.w2wp = make(map[string]string, len(T)*len(absolutes.Punctuation)*2)
	strip := runes.Remove(
		transform.Chain(
			runes.Map(unicode.ToLower),
			runes.Remove(unicode.IsPunct)))
	var wes strings.Builder
	for _, t := range rtn.raw {
		if t.Tag != "SYM" {
			k := transform.String(strip, t.Text)
			for p := range absolutes.Punctuation {
				rtn.w2wp[k+p] = k
				rtn.w2wp[p+k] = k
			}
			wes.WriteString(k)
			wes.WriteByte(' ')
		}
	}
	rtn.wes = wes
	rtn.Lexicon = lexicon
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
func (v Valence) chain(T Tokens, i int) float64 {
	// check the capitalization differential
	capdiff := T.ClDiff()
	scalar := absolutes.Boost(t)
	scalar *= sgn(v)
	if strings.ToUpper(t) == t && capdiff {
		scalar += sgn(v) * absolutes.CFactor
	}
	return scalar
}

// TODO: Implement actual heuristics and checks for sentiment analysis
func (v Valence) Gauge(doc Tokens) {

}
