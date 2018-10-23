package vader

import (
	"github.com/iseurie/vader-go/absolutes"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"unicode"
	"unicode/utf8"
)

// Lexicon represents a mapping from emoji to their respective descriptions,
// and terms to their respective ratings.
type Lexicon struct {
	Emotes  map[rune]string
	Ratings map[string]float64
}

// RateTokens will tstrip a text down to its constituent tokens
// and rate the result.
func (L Lexicon) RateTokens(raw string) []float64 {
	T := tstrip(raw)
	rtn := make([]float64, len(raw), len(raw))
	for _, t := range tstrip(raw) {
		r := []rune(t)
		if len(r) == 1 {
			rtn = append(rtn, L.EmRate(r[0]))
		} else {
			rtn = append(rtn, L.RatingGet(t))
		}
	}
	return rtn
}

// EmGet retrieves the description of an emoji.
func (L Lexicon) EmGet(e rune) string {
	if v, ok := L.Emotes[e]; ok {
		return strings.Join(tstrip(t), " ")
	}
	return ""
}

// EmRate retrieves the rating of an emoji by evaluating its description
// against the lexicon.
func (L Lexicon) EmRate(e rune) float64 {
	var rtn float64
	// the descriptions of emoji are generally naive enough to be simply summed,
	// rarely containing any complex language
	for _, t := range tstrip(L.EmGet(e)) {
		rtn += L.RatingGet(t)
	}
	return rtn
}

// RatingGet retrieves the rating of a token.
func (L Lexicon) RatingGet(s string) float64 {
	k := transform.String(s, normalize)
	if v, ok := L.Ratings[k]; ok {
		return v
	}
	return 0
}

// LoadLexicon loads a lexicon into memory. Where either emoji or dict are nil,
// the default resources are restored and used in their stead. Streams provided
// to LoadLexicon must be tab-delimited with one row per line. The emoji
// lexicon, if provided, must map UTF-8 emojis to their descriptions, and the
// lexicon must map English terms to their sentiment ratings as expressed in
// scientific notation.
func LoadLexicon(emoji, lexicon io.Reader) (Lexicon, error) {
	rtn := Lexicon{}
	if lexicon != nil {
		sc := bufio.NewScanner(lexicon)
		for sc.Scan() {
			var t string
			var r float64
			fmt.Fscanf(strings.Trim(sc.Text()), "%s\t%f", &t, &r)
			Lexicon.Ratings[t] = r
		}
		if err := sc.Err(); err != io.EOF {
			return err
		}
	} else {
		Lexicon.Ratings = absolutes.DefaultLexicon.Ratings
	}
	if emoji != nil {
		sc := bufio.NewScanner(emoji)
		for sc.Scan() {
			var emote rune
			var desc string
			fmt.Fscanf(strings.Trim(sc.Text()), "%c\t%s", &emote, &desc)
			Lexicon.Emotes[emote] = desc
		}
		if err := sc.Err(); err != io.EOF {
			return err
		}
	} else {
		Lexicon.Emotes = absolutes.DefaultLexicon.Emotes
	}
	return rtn, nil
}
