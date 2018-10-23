package absolutes

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func norm(raw string) string {
	return strings.TrimSpace(strings.ToLower(raw))
}

// Lexicon represents a mapping from emoji to their respective descriptions,
// and terms to their respective ratings.
type Lexicon struct {
	Emotes  map[string]string
	Ratings map[string]float64
}

// EmGet retrieves the description of an emoji.
func (L Lexicon) EmGet(e string) string {
	e = norm(e)
	if v, ok := L.Emotes[e]; ok {
		return v
	}
	return ""
}

// Rates checks whether the lexicon has a rating for a given token.
func (L Lexicon) Rates(s string) bool {
	k := norm(s)
	_, ok := L.Ratings[k]
	return ok
}

// Rating retrieves the rating of a token, returning zero where none is found.
func (L Lexicon) Rating(s string) float64 {
	k := norm(s)
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
func LoadLexicon(emoji, lexicon io.Reader) (L *Lexicon, err error) {
	L = new(Lexicon)
	if lexicon != nil {
		sc := bufio.NewScanner(lexicon)
		for sc.Scan() {
			var t string
			var r float64
			fmt.Sscanf(norm(sc.Text()), "%s\t%f", &t, &r)
			L.Ratings[t] = r
		}
		if err := sc.Err(); err != io.EOF {
			return nil, err
		}
	} else {
		L.Ratings = DefaultLexicon.Ratings
	}
	if emoji != nil {
		sc := bufio.NewScanner(emoji)
		for sc.Scan() {
			var emote string
			var desc string
			fmt.Sscanf(norm(sc.Text()), "%s\t%s", &emote, &desc)
			L.Emotes[emote] = desc
		}
		if err := sc.Err(); err != io.EOF {
			return nil, err
		}
	} else {
		L.Emotes = DefaultLexicon.Emotes
	}
	return
}
