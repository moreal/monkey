package repl

import (
	"bufio"
	"fmt"
	"github.com/moreal/monkey/lexer"
	"github.com/moreal/monkey/token"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer, err io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(err, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
