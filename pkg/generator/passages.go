package generator

import (
	"bytes"
	"io"
	"math/rand"
)

var (
	passages = []string{
		"There is some good in this world, and it’s worth fighting for.",
		"It is only with the heart that one can see rightly; what is essential is invisible to the eye.",
		"I am no bird; and no net ensnares me: I am a free human being with an independent will, which I now exert to leave you.",
		"It was the best of times, it was the worst of times, it was the age of wisdom, it was the age of foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of Darkness, it was the spring of hope, it was the winter of despair.",
		"Beware; for I am fearless, and therefore powerful.",
		"A man, after he has brushed off the dust and chips of his life, will have left only the hard, clean questions: Was it good or was it evil? Have I done well — or ill?",
		"This above all: To thine own self be true, And it must follow, as the night the day, Thou canst not then be false to any man.",
		"Why did you do all this for me?’ he asked. ‘I don’t deserve it. I’ve never done anything for you.’ ‘You have been my friend,’ replied Charlotte. ‘That in itself is a tremendous thing.",
		"It does not do to dwell on dreams and forget to live.",
		"When you play the game of thrones you win or you die.",
		"The world breaks everyone, and afterward, many are strong at the broken places.",
		"But soft! What light through yonder window breaks? It is the east, and Juliet is the sun.",
		"And, when you want something, all the universe conspires in helping you to achieve it.",
		"It is our choices, Harry, that show what we truly are, far more than our abilities.",
		"As Gregor Samsa awoke one morning from uneasy dreams he found himself transformed in his bed into an enormous insect.",
		"There is nothing like looking, if you want to find something. You certainly usually find something, if you look, but it is not always quite the something you were after.",
		"All we can know is that we know nothing. And that’s the height of human wisdom.",
		"Toto, I've a feeling we're not in Kansas anymore.",
		"They may take our lives, but they'll never take our freedom!",
		"Here's Johnny!",
		"It was beauty killed the beast.",
		"Keep your friends close, but your enemies closer.",
		"Life was like a box of chocolates; you never know what you’re gonna get.",
	}
)

func writePassages(desiredSize int64, writer io.Writer) error {
	var currentSize int
	b := new(bytes.Buffer)

	for currentSize < int(desiredSize) {
		p := passages[rand.Intn(len(passages))] + " "
		r := desiredSize - int64(currentSize)
		if len(p) > int(r) {
			p = p[0:r]
		}
		written, err := b.WriteString(p)
		if err != nil {
			return err
		}
		currentSize += written
	}
	io.WriteString(writer, b.String())
	return nil
}
