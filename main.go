package main

import (
	"flag"
	"fmt"
	"justify/asciiart"
	"os"
)

func main() {
	align := flag.String("align", "left", "Alignment type: left, center, right, justify")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
		fmt.Println("Example: go run . --align=right something standard")
		os.Exit(0)
	}

	if !isValidAlign(*align) {
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
		fmt.Println("Example: go run . --align=right something standard")
		os.Exit(0)
	}

	bannerFile := asciiart.GetBannerFileFromArgs(args)

	fileContent, err := asciiart.ReadBannerFile(bannerFile)
	if err != nil {
		fmt.Println("Error: Missing file")
		os.Exit(0)
	}

	lines := asciiart.SplitLines(fileContent, bannerFile)
	if len(lines) != 856 {
		fmt.Println("corrupt file")
		return
	}

	argument := args[0]

	termSize, err := asciiart.GetTerminalSize()
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		return
	}

	asciiart.PrintASCIIArt(lines, argument, *align, termSize.Width)
}

func isValidAlign(align string) bool {
	return align == "left" || align == "center" || align == "right" || align == "justify"
}
