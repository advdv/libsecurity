// Copyright 2010 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
//	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"code.google.com/p/freetype-go/freetype"
)

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "../../testdata/luxisr.ttf", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 12, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
)


// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

func main() {
	flag.Parse()

        if len(os.Args) != 3 {
          fmt.Fprintf(os.Stderr, "Usage: %s: <text name to read> <image name to create>", os.Args)
          return;
        }
        fileToRead := os.Args[1]
        fileToWrite := os.Args[2]

        // Read input text file        
        text, err := readLines(fileToRead)
        if err != nil {
            log.Fatalf("readLines: %s", err)
        }
	
        // Read the font data.
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	// Initialize the context.
	fg, bg := image.White, image.Black
	//ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	if *wonb {
		fg, bg = image.Black, image.White
	//	ruler = color.RGBA{0x22, 0x22, 0x22, 0xff}
	}
	rgba := image.NewRGBA(image.Rect(0, 0, 640, 480))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(font)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch *hinting {
	default:
		c.SetHinting(freetype.NoHinting)
	case "full":
		c.SetHinting(freetype.FullHinting)
	}

	// Draw the guidelines.
	//for i := 0; i < 200; i++ {
	//	rgba.Set(10, 10+i, ruler)
	//	rgba.Set(10+i, 10, ruler)
	//}

	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFix32(*size)>>8))
	for _, s := range text {
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFix32(*size * *spacing)
	}

	// Save that RGBA image to disk.
	f, err := os.Create(fileToWrite)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote %s OK.", fileToWrite)
}
