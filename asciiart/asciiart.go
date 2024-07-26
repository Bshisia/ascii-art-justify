package asciiart

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

type TerminalSize struct {
	Width  int
	Height int
}

func GetTerminalSize() (*TerminalSize, error) {
	var dimensions [4]uint16
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&dimensions)),
		0, 0, 0)
	if err != 0 {
		return nil, err
	}
	return &TerminalSize{
		Width:  int(dimensions[1]),
		Height: int(dimensions[0]),
	}, nil
}

func ReadBannerFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := fileInfo.Size()
	switch filename {
	case "standard.txt":
		if fileSize != 6623 {
			return "", fmt.Errorf("invalid file size for standard.txt: expected 6623, got %d", fileSize)
		}
	case "shadow.txt":
		if fileSize != 7463 {
			return "", fmt.Errorf("invalid file size for shadow.txt: expected 7463, got %d", fileSize)
		}
	case "thinkertoy.txt":
		if fileSize != 5558 {
			return "", fmt.Errorf("invalid file size for thinkertoy.txt: expected 5558, got %d", fileSize)
		}
	}
	buffer := make([]byte, fileSize)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}

func GetBannerFileFromArgs(args []string) string {
	if len(args) == 2 {
		switch args[1] {
		case "shadow":
			return "shadow.txt"
		case "thinkertoy":
			return "thinkertoy.txt"
		case "standard":
			return "standard.txt"
		default:
			fmt.Println(`Unidentified file, Usage: go run . "text" bannerfile`)
			os.Exit(0)
		}
	}
	if len(args) != 2 {
		fmt.Println("Usage: go run . [STRING] [BANNER]")
		os.Exit(0)
	}
	return "standard.txt"
}

func SplitLines(content string, filename string) []string {
	if filename == "thinkertoy.txt" {
		return strings.Split(content, "\r\n")
	}
	return strings.Split(content, "\n")
}

func PrintASCIIArt(lines []string, arguments string, align string, width int) {
	argument := strings.Split(arguments, "\n")
	for _, arg := range argument {
		for _, chr := range arg {
			if chr < 32 || chr > 126 {
				fmt.Println("Error: Non ASCII/printable characters found")
				os.Exit(0)
			}
		}
		if arg == "" {
			fmt.Println()
		} else {
			for i := 0; i < 8; i++ {
				var lineBuilder strings.Builder
				for _, value := range arg {
					start := int(value-32)*9 + 1
					lineBuilder.WriteString(lines[start+i])
				}
				fmt.Println(applyAlignment(lineBuilder.String(), align, width))
			}
		}
	}
}

func applyAlignment(line, align string, width int) string {
	switch align {
	case "center":
		padding := (width - len(line)) / 2
		return strings.Repeat(" ", padding) + line
	case "right":
		padding := width - len(line)
		return strings.Repeat(" ", padding) + line
	case "justify":
		words := strings.Fields(line)
		if len(words) == 1 {
			return line
		}
		spaces := width - len(strings.Join(words, ""))
		spaceBetween := spaces / (len(words) - 1)
		extraSpaces := spaces % (len(words) - 1)
		var result strings.Builder
		for i, word := range words {
			if i > 0 {
				result.WriteString(strings.Repeat(" ", spaceBetween))
				if extraSpaces > 0 {
					result.WriteString(" ")
					extraSpaces--
				}
			}
			result.WriteString(word)
		}
		return result.String()
	case "left":
		fallthrough
	default:
		return line
	}
}
