package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("error: ")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	sm := map[int]string{
		0:   "NULL",
		1:   "SOH",
		2:   "STX",
		3:   "ETX",
		4:   "EOT",
		5:   "EOQ",
		6:   "ACK",
		7:   "bell",
		8:   "backspace",
		9:   "tab",
		10:  "Newline",
		11:  "Vertical Tab",
		12:  "FF",
		13:  "CR",
		14:  "SO",
		15:  "SI",
		16:  "DLE",
		17:  "DC1",
		18:  "DC2",
		19:  "DC3",
		20:  "DC4",
		21:  "NAK",
		22:  "SYN",
		23:  "ETB",
		24:  "CAN",
		25:  "EM",
		26:  "SUB",
		27:  "Escape",
		28:  "FS",
		29:  "GS",
		30:  "RS",
		31:  "US",
		32:  "SPACE",
		33:  "Ex",         // ! exclamation symbol
		34:  "Double",     // " double quote
		35:  "Hash",       // # hash symbol
		36:  "Dollar",     // $ dollar symbol
		37:  "Percent",    // % percent cymbol
		38:  "And",        // & and symbol
		39:  "Quote",      // ' single quote symbol
		40:  "OP",         // ( opening parenthesis symbol
		41:  "CP",         // ( closing parenthesis symbol
		42:  "Star",       // * star symbol
		43:  "Plus",       // + plus symbol
		44:  "Camma",      // , camma symbol
		45:  "Minus",      // - minus symbol
		46:  "Period",     // . period symbol
		47:  "Slash",      // / slash symbol
		58:  "Colon",      // : colon symbol
		59:  "Semicolon",  // ; semicolon symbol
		60:  "LT",         // < less than symbol
		61:  "Equal",      // = eual symbol
		62:  "GT",         // > greater than symbol
		63:  "Q",          // ? question symbol
		64:  "At",         // @ atmark symbol
		91:  "OS",         // [ opening square bracket symbol
		92:  "Backslash",  // \ backslash symbol
		93:  "CS",         // ] closing square bracket symbol
		94:  "Carret",     // ^ carret symbol
		95:  "Underscore", // _ underscore symbol
		96:  "Tick",       // ` backtick symbol
		123: "OC",         // { opening curly brace
		124: "Pipe",       // | pipe symbol
		125: "CC",         // { closing curly brace
		126: "Tilda",      // ~ tilda symbol
		127: "Delete",
	}
	for i := 0; i < 128; i++ {
		aiff := fmt.Sprintf("%03d.aiff", i+1)
		wave := fmt.Sprintf("%03d.wav", i+1)
		text := string(rune(i))

		if s, ok := sm[i]; ok {
			text = s
		}
		if err := exec.Command("say", "-v", "Alex", "-r", "272", "-o", aiff, text).Run(); err != nil {
			return err
		}
		if err := exec.Command("sox", aiff, wave, "pitch", "-240").Run(); err != nil {
			return err
		}

		fmt.Printf("%s\tDONE\n", wave)
	}

	return nil
}
