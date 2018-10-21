package absolutes

import "github.com/iseurie/vader-go"

// DefaultLexicon is the default lexicon, indexing emoji against their ratings and
// all known terms against their sentiment ratings, respectively.
var DefaultLexicon = Lexicon{
	terms:  terms,
	emotes: emoji,
}
