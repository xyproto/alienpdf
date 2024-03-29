package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xyproto/randomstring"
)

const versionString = "AlienPDF 1.0.1"

func init() {
	randomstring.Seed()
}

var randStringFunc = randomstring.HumanFriendlyString

// Return an alien sentence. May return an empty string.
func sentence(wordCount int) string {
	var stringForHumans strings.Builder
	var maybeCommaCounter int
	for i := 0; i < wordCount; i++ {
		if maybeCommaCounter == 3 && rand.Intn(5) == 0 {
			stringForHumans.WriteString(", ")
			maybeCommaCounter = 0
		} else if i > 0 {
			stringForHumans.WriteString(" ")
		}
		word := randStringFunc(rand.Intn(5) + 2)
		if i == 0 {
			word = strings.Title(word)
		}
		stringForHumans.WriteString(word)
		maybeCommaCounter++
	}
	return stringForHumans.String()
}

// Return an enire letter
func messageFromMothership(sentenceCount int) string {
	var sb strings.Builder
	sb.WriteString(strings.Title(randStringFunc(3)))
	sb.WriteString(" ")
	sb.WriteString(strings.Title(randStringFunc(10)))
	sb.WriteString(",\n\n")
	lastWasNewline := true
	for i := 0; i < sentenceCount; i++ {
		if rand.Intn(5) == 0 {
			sb.WriteString("\n\n")
			lastWasNewline = true
		} else if !lastWasNewline {
			sb.WriteString(" ")
		} else {
			lastWasNewline = false
		}
		length := int(math.Round(math.Log2(float64(rand.Intn(450) + 1))))
		if length == 0 {
			continue
		}
		generatedString := sentence(length)
		if generatedString == "" {
			continue
		}
		sb.WriteString(generatedString)
		if !strings.Contains(generatedString, " ") {
			sb.WriteString("! ")
		} else {
			if rand.Intn(20) == 0 {
				sb.WriteString("? ")
			} else if rand.Intn(100) == 0 {
				sb.WriteString("! ")
			} else {
				sb.WriteString(". ")
			}
		}
		if i == 0 {
			lastWasNewline = false
		}
	}
	sb.WriteString("\n\n")
	sb.WriteString(strings.ToUpper(randStringFunc(rand.Intn(4) + 5)))
	sb.WriteString("!")
	return strings.ReplaceAll(sb.String(), "  ", " ")
}

func place() string {
	return strings.Title(randStringFunc(20))
}

func main() {
	fmt.Println(versionString)

	outputFilenameFlag := flag.String("o", "mothership.pdf", "an output PDF filename")
	englishLikeFlag := flag.Bool("e", false, "use a letter frequency more similar to English")
	randomFlag := flag.Bool("r", false, "use a more random letter frequency")
	flag.Parse()

	filename := *outputFilenameFlag

	if *englishLikeFlag {
		randStringFunc = randomstring.HumanFriendlyEnglishString
	} else if *randomFlag {
		randStringFunc = randomstring.EnglishFrequencyString
	}

	timestamp := time.Now().Format("2006-01-02")

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTopMargin(30)
	topLeftText := "1/1"
	topRightText := timestamp + ", " + place()
	pdf.SetHeaderFunc(func() {
		pdf.SetY(5)
		pdf.SetFont("Helvetica", "", 6)
		pdf.CellFormat(80, 0, topLeftText, "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 0, topRightText, "", 0, "R", false, 0, "")
	})
	pdf.AddPage()
	pdf.SetY(20)
	lines := strings.Split(messageFromMothership(40), "\n")
	pdf.SetFont("Courier", "B", 12)
	pdf.Write(5, lines[0]+"\n")
	pdf.SetFont("Courier", "", 12)
	pdf.Write(5, strings.Join(lines[1:len(lines)-1], "\n"))
	pdf.SetFont("Courier", "B", 12)
	pdf.Write(5, "\n"+lines[len(lines)-1])

	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s already exists\n", filename)
		os.Exit(1)
	}
	fmt.Printf("Writing %s... ", filename)
	if err := pdf.OutputFileAndClose(filename); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Println("done")
}
