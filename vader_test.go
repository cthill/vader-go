package vader

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/iseurie/vader-go/absolutes"
)

var tricky = [...]string{
	"Sentiment analysis has never been good.",
	"Sentiment analysis has never been this good!",
	"Most automated sentiment analysis tools are shit.",
	"With VADER, sentiment analysis is the shit!",
	"Other sentiment analysis tools can be quite bad.",
	"On the other hand, VADER is quite bad ass",
	"VADER is such a badass!", // slang with punctuation emphasis
	"Without a doubt, excellent idea.",
	"Roger Dodger is one of the most compelling variations on this theme.",
	"Roger Dodger is at least compelling as a variation on the theme.",
	"Roger Dodger is one of the least compelling variations on this theme.",
	"Not such a badass after all.",        // Capitalized negation with slang
	"Without a doubt, an excellent idea.", // "without {any} doubt" as negation
}

var stockRatings = [...]float64{
	-0.34, 0.56, 0.65, -0.58, 0.80, 0.40, 0.70, 0.29,
	0.23, -0.17, -0.26, 0.70,
}

var L = &absolutes.DefaultLexicon

func TestTrickySentences(t *testing.T) {
	fmt.Println("##. cmpd; neg.; pos.; neu.: input")
	tolerance := 0.10
	for i, s := range tricky {
		st := NewSentiText(s, L)
		S := st.ScoreValence()
		P := S.Polarity
		err := math.Abs(S.Compound - stockRatings[i])
		if err > tolerance {
			t.Errorf("Compound test error lies %f beyond tolerances on #%d: ``%s''.",
				err, i+1, tricky[i])
		} else {
			fmt.Printf(
				"%02d. %02.3f; %02.3f; %02.3f; %02.3f: %s\n",
				i+1, S.Compound, P.Negative, P.Positive, P.Neutral, s)
		}
	}
}

func BenchmarkSTConstruction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		si := tricky[rand.Int()%len(tricky)]
		b.ResetTimer()
		NewSentiText(si, L)
	}
}

func BenchmarkSTScoring(b *testing.B) {
	S := make([]SentiText, len(tricky), len(tricky))
	for i := 0; i < len(tricky); i++ {
		S[i] = NewSentiText(tricky[i], L)
	}

	for i := 0; i < b.N; i++ {
		st := S[rand.Int()%len(tricky)]
		b.ResetTimer()
		st.ScoreValence()
	}
}
