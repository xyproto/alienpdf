package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xyproto/randomstring"
)

func init() {
	randomstring.Seed()
}

// Return an alien sentence. May return an empty string.
func sentence(wordCount int) string {
	var stringForHumans strings.Builder
	for i := 0; i < wordCount; i++ {
		if i > 0 {
			stringForHumans.WriteString(" ")
		}
		word := randomstring.HumanFriendlyString(rand.Intn(5) + 2)
		if i == 0 {
			word = strings.Title(word)
		}
		stringForHumans.WriteString(word)
	}
	return stringForHumans.String()
}

// Return an enire letter
func messageFromMothership(sentenceCount int) string {
	var stringsForHumans strings.Builder
	stringsForHumans.WriteString(strings.Title(randomstring.HumanFriendlyString(3)))
	stringsForHumans.WriteString(" ")
	stringsForHumans.WriteString(strings.Title(randomstring.HumanFriendlyString(10)))
	stringsForHumans.WriteString(",\n\n")
	lastWasNewline := false
	for i := 0; i < sentenceCount; i++ {
		if i > 0 {
			if rand.Intn(5) == 0 {
				stringsForHumans.WriteString("\n\n")
				lastWasNewline = true
			} else if !lastWasNewline {
				stringsForHumans.WriteString(" ")
			} else {
				lastWasNewline = false
			}
		}
		length := int(math.Round(math.Log2(float64(rand.Intn(450) + 1))))
		if length == 0 {
			continue
		}
		generatedString := sentence(length)
		if generatedString == "" {
			continue
		}
		stringsForHumans.WriteString(generatedString)
		if !strings.Contains(generatedString, " ") {
			stringsForHumans.WriteString("!")
		} else {
			if rand.Intn(40) == 0 {
				stringsForHumans.WriteString("?")
			} else if rand.Intn(100) == 0 {
				stringsForHumans.WriteString("!")
			} else {
				stringsForHumans.WriteString(".")
			}
		}
	}
	stringsForHumans.WriteString("\n\n")
	stringsForHumans.WriteString(strings.ToUpper(randomstring.HumanFriendlyString(rand.Intn(4) + 5)))
	stringsForHumans.WriteString("!")
	return stringsForHumans.String()
}

func place() string {
	return strings.Title(randomstring.HumanFriendlyString(20))
}

func main() {
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

	filename := "mothership.pdf"
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s already exists!\n", filename)
		os.Exit(1)
	}
	fmt.Printf("Writing %s... ", filename)
	if err := pdf.OutputFileAndClose(filename); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Println("done.")
}
