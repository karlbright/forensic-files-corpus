package forensicfilescorpus

import (
	"errors"
	"math"
	"math/rand"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wargarblgarbl/libgosubs/srt"
)

// MinimumLineLength used to determine the minimum length for a subtitle line in order to be used.
// This is here initially to avoid weird issues with some Forensic Files subtitles that were
// breaking across lines for a name, particular for a police officers title. Such as "Lt."
const MinimumLineLength = 8

// RemoveFromSubtitleRegexp matches things we do not want to have as part of our subtitles for
// various reasons. This covers, in the same order as the regexp:
// - Single source change line, such as "DIANNE M. ANDERSON:"
// - Dialogue source change, such as "Narrator: They ran away" or "Skip Palenik: I'm a genius!"
// - Conversation subtitles that can appear next to each other on the same screen.
// - Actions and descriptive audio subtitles such as "[sirens]" and "[theme music]"
var RemoveFromSubtitleRegexp = regexp.MustCompile(`^(.+:$|.+: |-|\[.+\])`)

// IgnoreSubtitleRegexp matches subtitles that are not formatted correctly. This is primarily
// used to avoid subtitles that exist from Youtube subtitles, and other sources that are unknown
// to me. Some subtitles from the subtitles I have used were in ALL CAPS, and some of them
// came from youtube and contained HTML, like "<font color="#CCCCC">Foo</Foo>".
var IgnoreSubtitleRegexp = regexp.MustCompile(`^(>> Narrator:|Narrator:|^[^a-z]+$|<\/?.+?>)`)

// StartToken matches against lines that can be used to begin a sentence. We do this by checking for
// a capital letter, or a number. We also ensure that the line does not end with a quotation mark
// to avoid an edge case where some titles and dialogue was being seen as a sentence in itself.
var StartToken = regexp.MustCompile(`^[A-Z0-9].+[^"]$`)

// EndToken matches against lines that we feel comfortable in ending a sentence in. These are the
// common characters that will end a sentence.
// - question mark
// - excalamation mark
// - a full stop
// - a double quotation mark.
// Although we do not accept this within `StartToken` at the end of a line, we allow it here as
// we can be confident that a subtitle will not split a quote from someone over a line with the
// double quotation mark being by itself at the end of a previous line. We hope.
var EndToken = regexp.MustCompile(`(\?|!|\.|")$`)

// StripAll is a convenience method to strip relevant subtitle sentences from a number of
// subtitle files. See `Strip` for more information on how the subtitles are stripped.
func StripAll(paths []string) (all []string) {
	for _, path := range paths {
		sentences, err := Strip(path)

		if err != nil {
			continue
		}

		all = append(all, sentences...)
	}

	return all
}

// Strip will remove any sentences from a subtitle with replacements for things such as conversations,
// dialogue target changes, descriptive audio lines, etc. We also make sure that the subtitle we are
// stripping does not contain any ignored subtitles. See `IgnoreSubtitleRegexp` for more information
// on what can cause a subtitle file to be ignored. In the case of a subtitle file encountering
// a subtitle that matches the ignoring rules, then the whole subtitle is ignored.
func Strip(path string) (sentences []string, err error) {
	target, err := filepath.Abs(path)

	if err != nil {
		return sentences, errors.New("unable to retrieve absolute path for target")
	}

	subtitles, err := srt.ParseSrt(target)

	if err != nil {
		return sentences, errors.New("error parsing subtitle file")
	}

	var lines []string

	for _, subtitle := range subtitles.Subtitle.Content {
		subtitle := strings.Join(subtitle.Line, " ")
		subtitle = RemoveFromSubtitleRegexp.ReplaceAllString(subtitle, "")

		if IgnoreSubtitleRegexp.MatchString(subtitle) {
			return sentences, errors.New("ignored subtitle file")
		}

		if subtitle != "" && len(subtitle) > MinimumLineLength {
			lines = append(lines, subtitle)
		}
	}

	for index, line := range lines {
		if StartToken.MatchString(line) {
			if EndToken.MatchString(line) {
				sentences = append(sentences, line)
			} else {
				currentIndex := index
				currentLine := line

				for {
					currentIndex = currentIndex + 1
					if currentIndex > len(lines)-1 {
						break
					}

					currentLine = strings.Join([]string{currentLine, lines[currentIndex]}, " ")

					if EndToken.MatchString(currentLine) {
						sentences = append(sentences, currentLine)
						break
					}
				}
			}
		}
	}

	return sentences, nil
}

// Pick a random sentence from a collection of sentences provided as the first argument to the
// method. A minimum and maximum chracter length for the sentence can be used to filter the
// sentences. Passing a negative value to either one of these will use sensible defaults.
// Uses the default source for randomisation, it is advised that you seed the default source
// using something like `rand.Seed(time.Now().UnixNano()` in order to ensure you are picking
// random values each time, rather than using the same deterministic seed.
func Pick(sentences []string, min, max int) (string, error) {
	if len(sentences) == 0 {
		return "", errors.New("unable to pick from empty sentences slice")
	}

	n := rand.Intn(len(sentences))

	if min < 0 {
		min = 0
	}

	if max < 0 {
		max = math.MaxInt32
	}

	if min > max {
		return "", errors.New("min value must be smaller than max")
	}

	if max < MinimumLineLength {
		return "", errors.New("max value must be larger than the minimum sentence length")
	}

	if max == -1 {
		return sentences[n], nil
	}

	var filtered []string
	for _, sentence := range sentences {
		if len(sentence) > min && len(sentence) < max {
			filtered = append(filtered, sentence)
		}
	}

	if len(filtered) == 0 {
		return "", errors.New("no candidates with given min and max values")
	}

	n = rand.Intn(len(filtered))
	return filtered[n], nil
}

// Generate a random paragraph which can consist of one or many random sentences. This uses the
// `Pick` method to pick a random sentence of a given length. A minimum and maximum character length
// for the final paragraph can be provided. Passing a negative value to either one of these will
// use sensible defaults. As this uses `Pick` to determine which random sentence is used, it should
// be noted that you should seed the default source for randomisation. See `Pick` for more details.
func Generate(sentences []string, min, max int) (out string, err error) {
	for {
		sentence, err := Pick(sentences, -1, max-len(out))

		if err != nil {
			return out, err
		}

		out = out + sentence

		if len(out) > min {
			break
		}

		out = out + " "
	}

	return out, nil
}
