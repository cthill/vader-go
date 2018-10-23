package vader

import (
	"strings"
	"unicode"
)

// SentiText is an intermediate struct for evaluating sentiment of a given document.
type SentiText struct {
	// The cartesian product of absolutes.Punctuation the set of tokens in
	// the list, permuting each token by prepending and appending punctuation.
	w2wp map[string]string

	// Remove leading and trailing punctuation, but leave contractions and emoticons
	wes Tokens

	raw Tokens

	// Lexicon is the structure's lexicon
	Lexicon

	clDiff bool
}

// NewSentiText constructs a new sentiment evaluation model from the given
// parameters.
func NewSentiText(raw string, lexicon Lexicon) (rtn SentiText) {
	rtn.raw.OfRaw(raw, true)
	rtn.w2wp = make(map[string]string, len(T)*len(absolutes.Punctuation)*2)
	var wes strings.Builder
	for _, t := range rtn.raw {
		if t.Tag != "SYM" { // preserve emotes
			k := t.Text
			a := 0
			b := len(k)-1
			for ; unicode.IsPunct(k[a]); a++ {}
			for ; unicode.IsPunct(k[b]); b-- {}
			k = k[a+1:b]
			k := transform.String(strip, t.Text)
			for p := range absolutes.Punctuation {
				rtn.w2wp[k+p] = k
				rtn.w2wp[p+k] = k
			}
			wes.WriteString(k)
			wes.WriteByte(' ')
		}
	}
	rtn.wes = wes.String()
	rtn.Lexicon = lexicon
	rtn.clDiff = rtn.raw.ClDiff()
}

func (S SentiText) PolarityScores() []float64 {
	rtn := make([]float64, len(S.wes))
	valence := 0
	for _, x := range S.wes {
		k := dncase(s)
		switch {
		case absolutes.IsBoosted(k):
		case i < len(S.wes-1) && k == "kind" && strings.ToLower(S.wes[i+1]) == "of":
			rtn = append(rtn, valence)
			continue
		default:
			valence := S.clDiff
	}
}

func (S SentiText) SentimentValence(valence float64, item string)
