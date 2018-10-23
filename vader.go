package vader

import (
	"github.com/iseurie/vader-go/absolutes"
	"math"
	"strings"
)

func sgn(x float64) float64 {
	if x < 0 {
		return -1
	}
	return 1
}

func negate(t string) bool {
	switch {
	case absolutes.Negate(t):
		fallthrough
	case strings.Contains(t, "n't"):
		fallthrough
	case t == "least" && i > 0 && raw[i-1] != "at":
		return true
	}
	return false
}

func negated(raw []string) bool {
	for i, t := range raw {
		if negate(t) {
			return true
		}
	}
	return false
}

func triplet(S []string) (a, b, c string) {
	var T [3]string
	for i := 0; i < 2 && i < len(S); i++ {
		T[i] = S[i]
	}
	a, b, c = rtn[0], rtn[1], rtn[2]
	return
}

func negationCk(float64 valence, V []string) float64 {
	nst := func(a, b string) bool {
		return a == "never" && (b == "so" || b == "this")
	}
	a, b, c := triplet(V)
	switch {
	case negate(a), negate(b), negate(c):
		return valence * absolutes.NScalar
	case nst(a, b):
		return valence * 5 / 4
	case nst(b, c):
		return valence * 5 / 4
	}
}

// SentiText is an intermediate struct used to evaluate sentiment.
type SentiText struct {
	raw, wes []string
	acdiff   bool
}

// NewSentiText constructs a new document from which to measure sentiment.
func NewSentiText(raw string, L Lexicon) (new SentiText) {
	for i, r := range raw {
		if desc, ok := L.Emotes[r]; ok {
			raw = fmt.Sprintf("%s %s %s", raw[:i-1], desc, raw[i:])
		}
	}
	new.raw = strings.Fields(raw)
	new.wes = new.wes(new.raw)
	clc := 0
	for _, t = range T {
		if strings.ToUpper(t) == t {
			clc++
		}
	}
	cdiff := len(T) - clc
	new.acdiff = 0 < cdiff && cdiff < len(T)
	return new
}

func (ST SentiText) wes(raw []string) []string {
	W := make([]string, len(raw), len(raw))
	p2w := make(map[string]string, len(W))
	for i, t := range raw {
		if len < 2 {
			continue
		}
		if _, ok := p2w[t+p2w[p[0]]]; !ok {
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
	copy(T, ST.raw)
	for i := range T {
		T[i] = strings.ToLower(T[i])
	}
	return T
}

// Sentiments attempts to gauge the text's valence sentiments.
func (ST SentiText) Sentiments(L Lexicon) []float64 {
	rtn := make([]float64, len(ST.raw), len(ST.raw))
	for i, t := range ST.wes {
		tl := strings.Lower(t)
		valence := 0
		if absolutes.IsBoosted(k) ||
			i < len(ST.wes)-1 ||
			tl == "kind" && wes[i+1] == "of" {
			rtn = append(rtn, valence)
		}
		// update valence
		capp := strings.ToUpper(t) == t && ST.acdiff
		if capp {
			valence += sgn(valence) * absolutes.CScalar
		}
		wesl := ST.wesl()
		for j := 0; j < 2; j++ {
			if i <= j || L.Rates(wesl[i-j+1]) {
				continue
			}
			// scalar_inc_dec
			scalar := sgn(valence) * absolutes.Boost(tl)
			if capp {
				scalar += sgn(valence) * absolutes.CScalar
			}
			scalar *= (1 + (float64(j-2) / 20))
			valence += scalar
			// negation_check
			k := i - 3
			if iCk < 0 {
				iCk = 0
			}
			b, c, d := triplet(wesl[iCk:i])
			a := wesl[i]
			negationCk(wes[iCk:i])
			if j == 2 {
				// _special_idioms_check
				trigrams := [][]string{
					{c, b, a},
					{d, c, b},
					{a, b, c},
				}
				bigrams := [][]string{
					{a, b},
					{b, a},
					{c, b},
				}
				for _, seq := range append(bigrams, trigrams) {
					s := strings.Join(" ", seq)
					// check for boosting/dampening & sentiment-laden idioms
					valence += absolutes.Boost(s)
					if v, ok := absolutes.SentimentLadenIdioms[s]; ok {
						valence += v
					}
				}
			}
		}
		// _least_check
		if i > 1 && !L.Rates(b) && b == "least" {
			if c != "at" && c != "very" {
				valence *= absolutes.NScalar
			}
		}
		rtn[i] = valence
	}
	// _but_check
	for bi := range rtn {
		// we got one!
		if wesl[bi] == "but" {
			for si := range rtn {
				if si < bi {
					rtn[i] = 0.5 * rtn[i]
				} else if si > bi {
					rtn[i] = 1.5 * rtn[i]
				}
			}
		}
	}
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
	Σ := 0
	for _, x := range rtn {
		Σ += x
	}
	epc := 0
	qmc := 0
	for _, t := range ST.raw {
		for _, c := range t {
			if epc == 4 && qmc == 4 {
				break
			}
			switch c {
			case '?' && qmc < 4:
				qmc++
			case '!' && epc < 4:
				epc++
			}
		}
	}
	pEmph := float64(epc) * 0.292
	if qmc <= 3 {
		pEmph += float64(qmc) * 0.18
	} else {
		pEmph += 0.96
	}
	S.Polarity = ST.Sift()
	S.Compound = Σ / math.Sqrt((Σ*Σ)+15)
	S.Compound = math.Max(score, -1)
	S.Compound = math.Min(score, 1)
	return
}

// Sift obtains polarity ratings for the text.
func (ST SentiText) Sift(P Polarity) {
	S := ST.Sentiments()
	for _, x := range S {
		if x > 0 {
			P.Positive += x
		} else if x < 0 {
			P.Negative += x
		} else {
			P.Neutral++
		}
	}
	if P.Positive > math.Abs(P.Negative) {
		P.Positive += pEmph
	} else {
		P.Negative -= pEmph
	}
	total := P.Positive + math.Abs(P.Negative) + P.Neutral
	P.Positive /= total
	P.Negative /= total
	P.Neutral /= total
	return
}
