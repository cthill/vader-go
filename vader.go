package vader

import (
	"fmt"
	"github.com/iseurie/vader-go/absolutes"
	"math"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

const debug = false

func sgn(x float64) float64 {
	if x < 0 {
		return -1
	}
	return 1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func strip(raw string) string {
	rtn, _, _ := transform.String(runes.Remove(runes.In(unicode.Punct)), raw)
	return rtn
}

func negatep(t string) bool {
	return absolutes.Negate(t) || strings.Contains(t, "n't")
}

func negated(raw []string) bool {
	for i, t := range raw {
		if negatep(t) || (i > 0 && t == "least" && raw[i-1] != "at") {
			return true
		}
	}
	return false
}

func backfill(src []string, i int, out ...*string) {
	for d := len(out); d > 0; d-- {
		j := i - d
		if j < 0 || j > len(src)-1 {
			break
		}
		if out[d-1] != nil {
			*out[d-1] = src[j]
		}
	}
}

func fillback(src []string, i int, out ...*string) {
	for d := 0; d < len(out)-1; d++ {
		j := i + d
		if j > len(src)-1 {
			break
		}
		if out[d] != nil {
			*out[d] = src[j]
		}
	}
}

// SentiText is an intermediate struct used to evaluate sentiment.
type SentiText struct {
	raw, wes, wesl []string
	acdiff         bool
	L              *absolutes.Lexicon
}

// NewSentiText constructs a new document from which to measure sentiment.
func NewSentiText(raw string, L *absolutes.Lexicon) (new SentiText) {
	new.L = L
	for i, r := range raw {
		if desc, ok := L.Emotes[string(r)]; ok {
			raw = fmt.Sprintf("%s %s %s", raw[:i-1], desc, raw[i:])
		}
	}
	new.raw = strings.Fields(strip(raw))
	new.wes = new.mkwes(new.raw)
	new.wesl = new.mkwesl(new.raw)
	clc := 0
	for _, t := range new.raw {
		if strings.ToUpper(t) == t {
			clc++
		}
	}
	cdiff := len(new.raw) - clc
	new.acdiff = 0 < cdiff && cdiff < len(new.raw)
	return new
}

func (ST SentiText) scalarIncDec(t string, valence float64) float64 {
	tl := strings.ToLower(t)
	scalar := sgn(valence) * absolutes.Boost(tl)
	if ST.acdiff && strings.ToUpper(t) == t {
		scalar += sgn(valence) * absolutes.CScalar
	}
	return scalar
}

func (ST SentiText) negationCk(valence float64, i int) float64 {
	nst := func(a, b string) bool {
		return a == "never" && (b == "so" || b == "this")
	}
	// wdt := func(a, b string) bool {
	// 	return a == "without" && b == "doubt"
	// }
	var a, b, c string
	backfill(ST.wesl, i, nil, &a, &b, &c)
	if negated([]string{a, b, c}) {
		valence *= absolutes.NScalar
	}
	if nst(a, b) || nst(b, c) || nst(a, c) {
		return valence * 5 / 4
	}
	// else if wdt(a, b) || wdt(a, c) {
	// 	return valence
	// }
	return valence
}

func (ST SentiText) leastCk(valence float64, i int) float64 {
	var a, b string
	backfill(ST.wesl, i, nil, &a, &b)
	if !ST.L.Rates(a) && a == "least" &&
		b != "very" && b != "at" {
		valence *= absolutes.NScalar
	}
	return valence
}

func (ST SentiText) butCk(sentiments []float64) {
	for bi, t := range ST.wesl {
		if t == "but" {
			for si := range sentiments {
				if si < bi {
					sentiments[si] *= 0.5
				} else if si > bi {
					sentiments[si] *= 1.5
				}
			}
		}
	}
}

func (ST SentiText) specialIdiomsCk(valence float64, i int) float64 {
	var a, b, c, d string
	backfill(ST.wesl, i, &a, &b, &c, &d)
	trigrams := []string{
		strings.Join([]string{c, b, a}, " "),
		strings.Join([]string{d, c, b}, " "),
	}
	bigrams := []string{
		strings.Join([]string{b, a}, " "),
		strings.Join([]string{c, b}, " "),
		strings.Join([]string{d, c}, " "),
	}
	switch {
	case len(ST.raw)-1 > i+1:
		trigrams = append(trigrams, strings.Join(ST.wesl[i:i+3], " "))
		fallthrough
	case len(ST.raw)-1 > i:
		bigrams = append(bigrams, strings.Join(ST.wesl[i:i+2], " "))
	}
	for _, k := range bigrams {
		valence = absolutes.SpecialCaseIdioms[k]
		valence += absolutes.Boost(k)
	}
	for _, k := range trigrams {
		valence = absolutes.SpecialCaseIdioms[k]
	}
	return valence
}

func (ST SentiText) sentimentValence(valence float64, i int, sentiments []float64) {
	t := ST.wes[i]
	defer func() {
		sentiments = append(sentiments, valence)
	}()
	if !L.Rates(t) {
		return
	}
	if strings.ToUpper(t) == t && ST.acdiff {
		valence += sgn(valence) * absolutes.CScalar
	}
	for j := 0; j < 2; j++ {
		s := ST.scalarIncDec(ST.wes[i-(j+1)], valence)
		// distance damping
		s *= (1 + (float64(j-2) / 20))
		valence += s
		valence = ST.negationCk(valence, i)
		if j == 2 {
			valence = ST.specialIdiomsCk(valence, i)
		}
	}
	valence = ST.leastCk(valence, i)
}

func (ST SentiText) mkwes(raw []string) []string {
	W := make([]string, len(raw), len(raw))
	p2w := make(map[string]string, len(W))
	for _, t := range raw {
		if len(t) < 2 {
			continue
		}
		k := strip(t)
		if _, ok := p2w[k+p2w[absolutes.Punctuation[0]]]; !ok {
			for _, p := range absolutes.Punctuation {
				p2w[k+p] = t
				p2w[p+k] = t
			}
		}
		W = append(W, p2w[k])
	}
	return W
}

func (ST SentiText) mkwesl(raw []string) []string {
	T := make([]string, len(raw), len(raw))
	for i := range T {
		T[i] = strings.ToLower(T[i])
	}
	return T
}

func (ST SentiText) punctEmph() float64 {
	epc := 0
	qmc := 0
	for _, t := range ST.wes {
		for _, c := range t {
			if c == '?' && qmc < 4 {
				qmc++
			} else if c == '!' && epc < 4 {
				epc++
			} else {
				break
			}
		}
	}
	pEmph := float64(epc) * 0.292
	if qmc <= 3 {
		pEmph += float64(qmc) * 0.18
	} else {
		pEmph += 0.96
	}
	return pEmph
}

// Polarities attempts to gauges the text's valence sentiments.
func (ST SentiText) Polarities() []float64 {
	rtn := make([]float64, len(ST.raw), len(ST.raw))
	for i, t := range ST.wesl {
		var valence float64
		if absolutes.IsBoosted(t) ||
			(i < len(ST.wesl)-1 && strings.ToLower(t) == "kind" &&
				strings.ToLower(ST.wesl[i+1]) == "of") {
			rtn = append(rtn, valence)
			continue
		}
		ST.sentimentValence(valence, i, rtn)
	}
	ST.butCk(rtn)
	return rtn
}

// Polarity represents the result of sifting the valence sentiments read from a
// text.
type Polarity struct {
	Positive, Negative, Neutral float64
}

// Sentiment includes the compound score in the sentiment assessment.
type Sentiment struct {
	Polarity
	Compound float64
}

func (ST SentiText) sift(S []float64) (P Polarity) {
	for _, x := range S {
		if x > 0 {
			P.Positive += (x + 1)
		} else if x < 0 {
			P.Negative -= (x - 1)
		} else {
			P.Neutral++
		}
	}
	pEmph := ST.punctEmph()
	d := P.Positive - math.Abs(P.Negative)
	if d > 0 {
		P.Positive += pEmph
	} else if d < 0 {
		P.Negative -= pEmph
	}
	total := P.Positive + math.Abs(P.Negative) + P.Neutral
	P.Positive /= total
	P.Negative /= total
	P.Neutral /= total
	return
}

// ScoreValence returns Sentiment of the text.
func (ST SentiText) ScoreValence() (S Sentiment) {
	var Σ float64
	V := ST.Polarities()
	for _, x := range V {
		Σ += x
	}

	S.Polarity = ST.sift(V)
	// Normalize
	S.Compound = Σ / math.Sqrt((Σ*Σ)+15)
	S.Compound = math.Max(Σ, -1)
	S.Compound = math.Min(Σ, 1)
	return
}
