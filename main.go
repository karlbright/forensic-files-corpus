package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/wargarblgarbl/libgosubs/srt"
)

const minimumLineLength = 8

func main() {
	if len(os.Args) < 1 {
		log.Fatal("Missing source or destination for sentences")
	}

	ss := flag.Bool("strip", false, "Strip sentences from subtitles and write to disk")
	gg := flag.Bool("generate", false, "Generate medium sized paragraph with 1 or more randomly picked sentences")
	flag.Parse()

	if *ss {
		f, err := os.OpenFile("sentences.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		sentences, _ := strip()

		for _, sentence := range sentences {
			f.WriteString(sentence + "\n")
		}
	} else {
		var sentences []string

		f, err := os.Open("sentences.txt")
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		sc := bufio.NewScanner(f)
		for sc.Scan() {
			sentences = append(sentences, sc.Text())
		}

		if *gg {
			fmt.Println(strings.TrimSpace(generate(sentences)))
		} else {
			fmt.Println(pick(sentences))
		}

	}
}

func strip() ([]string, []string) {
	var lines []string

	files, err := filepath.Glob("in/*.srt")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		var slines []string

		targetPath, _ := filepath.Abs(file)
		subtitles, _ := srt.ParseSrt(targetPath)

		re := regexp.MustCompile(`^(.+: |-|\[.+\])`)
		ig := regexp.MustCompile(`^(>> Narrator:|Narrator:|^[^a-z]+$|<\/?.+?>)`)
		skip := false

		for _, sub := range subtitles.Subtitle.Content {
			line := strings.Join(sub.Line, " ")
			line = re.ReplaceAllString(line, "")

			if ig.MatchString(line) || skip == true {
				skip = true
				break
			}

			if line != "" && len(line) > minimumLineLength {
				slines = append(slines, line)
			}
		}

		if skip == false {
			lines = append(lines, slines...)
		}
	}

	var sentences []string
	st := regexp.MustCompile(`^[A-Z0-9].+[^"]$`)
	en := regexp.MustCompile(`(\?|!|\.|")$`)

	for i, l := range lines {
		if st.MatchString(l) {
			if en.MatchString(l) {
				sentences = append(sentences, l)
			} else {
				ti := i
				tl := l

				for {
					ti = ti + 1

					if ti > len(lines)-1 {
						break
					}

					tl = strings.Join([]string{tl, lines[ti]}, " ")

					if en.MatchString(tl) {
						sentences = append(sentences, tl)
						break
					}
				}
			}
		}
	}

	return sentences, lines
}

func pick(sentences []string) string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	return sentences[r.Intn(len(sentences)-1)]
}

func generate(sentences []string) string {
	var str string

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for {
		nstr := strings.Join([]string{str, sentences[r.Intn(len(sentences)-1)]}, " ")

		if len(nstr) > 280 {
			continue
		}

		str = strings.Replace(nstr, "--", "-", -1)

		if len(str) > 220 {
			break
		}
	}

	return str
}
