package vader

import (
	"fmt"
	"github.com/cthill/vader-go/absolutes"
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
		if negatep(t) || (t == "least" && i > 0 && raw[i-1] != "at") {
			return true
		}
	}
	return false
}

func triplet(S []string) (a, b, c string) {
	var T [3]string
	for i := 0; i < 2; i++ {
		if i < len(S)-1 {
			T[i] = S[i]
		} else {
			T[i] = ""
		}
	}
	a, b, c = T[0], T[1], T[2]
	return
}

func negationCk(valence float64, V []string) float64 {
	nst := func(a, b string) bool {
		return a == "never" && (b == "so" || b == "this")
	}
	// wdt := func(a, b string) bool {
	// 	return a == "without" && b == "doubt"
	// }
	a, b, c := triplet(V)
	if negated([]string{a, b, c}) {
		return valence * absolutes.NScalar
	} else if nst(a, b) || nst(b, c) {
		return valence * 5 / 4
	}
	// else if wdt(a, b) || wdt(a, c) {
	// 	return valence
	// }
	return valence
}

// SentiText is an intermediate struct used to evaluate sentiment.
type SentiText struct {
	raw, wes []string
	acdiff   bool
	L        *absolutes.Lexicon
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
	new.wes = new.mkwes(strings.Fields(raw))
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

func (ST SentiText) mkwes(raw []string) []string {
	W := make([]string, len(raw), len(raw))
	p2w := make(map[string]string, len(W))
	for i, t := range raw {
		if len(t) < 2 {
			continue
		}
		if _, ok := p2w[t+p2w[absolutes.Punctuation[0]]]; !ok {
			for _, p := range absolutes.Punctuation {
				p2w[t+p] = t
				p2w[p+t] = t
			}
		}
		W[i] = p2w[t]
	}
	return W
}

func (ST SentiText) wesl() []string {
	T := make([]string, len(ST.raw), len(ST.raw))
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

// Sentiments attempts to gauge the text's valence sentiments.
func (ST SentiText) Sentiments() []float64 {
	rtn := make([]float64, len(ST.raw), len(ST.raw))
	wesl := ST.wesl()
	for i, t := range ST.wes {
		capp := strings.ToUpper(t) == t && ST.acdiff
		tl := strings.ToLower(t)
		valence := ST.L.Rating(tl)
		if absolutes.IsBoosted(tl) ||
			i < len(ST.wes)-1 ||
			(tl == "kind" && strings.ToLower(ST.wes[i+1]) == "of") ||
			!ST.L.Rates(tl) {
			// update valence
			rtn = append(rtn, valence)
			continue
		}
		// sentiment_valence
		for j := 0; j < 2; j++ {
			if i <= j || ST.L.Rates(wesl[i-(j+1)]) {
				continue
			}
			// scalar_inc_dec
			scalar := sgn(valence) * absolutes.Boost(tl)
			if capp {
				scalar += sgn(valence) * absolutes.CScalar
			}
			// distance damping
			scalar *= (1 + (float64(j-2) / 20))
			valence += scalar
			// negation_check
			iCk := min(i-3, 0)
			b, c, d := triplet(wesl[iCk:i])
			a := wesl[i]
			if i > 3 && debug {
				fmt.Printf("{a b c d} = %v\n", []string{a, b, c, d})
				fmt.Printf(" = %#v\n", wesl[i-3:i+1])
			}
			valence = negationCk(valence, wesl[iCk:i])
			if j == 2 {
				// _special_idioms_check
				trigrams := [][]string{
					{c, b, a},
					{d, c, b},
				}
				bigrams := [][]string{
					{b, a},
					{c, b},
					{d, c},
				}
				if len(wesl)-1 > i {
					bigrams = append(bigrams, []string{a, b})
					trigrams = append(trigrams, []string{a, b, c})
				}
				for _, seq := range append(bigrams, trigrams...) {
					s := strings.Join(seq, " ")
					// check for boosting/dampening & sentiment-laden idioms
					valence += absolutes.Boost(s)
					if v, ok := absolutes.SentimentLadenIdioms[s]; ok {
						valence += v
					}
				}
			}
		}
		j := min(i-2, 0)
		a, b, _ := triplet(wesl[j:i])
		// _least_check
		if !ST.L.Rates(a) && a == "least" && b != "at" {
			rtn[i] *= absolutes.NScalar
		}
		rtn[i] = valence
	}
	// _but_check
	for bi := range wesl {
		// we got one!
		if wesl[bi] == "but" {
			for si := range rtn {
				if si < bi {
					rtn[si] *= 0.5
				} else if si > bi {
					rtn[si] *= 1.5
				}
			}
		}
	}
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

// ScoreValence returns Sentiment of the text.
func (ST SentiText) ScoreValence() (S Sentiment) {
	// score_valence
	var Σ float64
	V := ST.Sentiments()
	for _, x := range V {
		Σ += x
	}

	S.Polarity = ST.Sift(V)
	// Normalize
	S.Compound = Σ / math.Sqrt((Σ*Σ)+15)
	S.Compound = math.Max(Σ, -1)
	S.Compound = math.Min(Σ, 1)
	return
}

// Sift obtains polarity ratings for the text.
func (ST SentiText) Sift(S []float64) (P Polarity) {
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
