package chain

import (
	"fmt"
	"strings"
	"unicode"
)

// Token represents a lexical token
type Token struct {
	Type  TokenType
	Value string
	Pos   int
}

// TokenType represents the type of a token
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenPipe
	TokenIdentifier
	TokenColon
	TokenEquals
	TokenString
	TokenError
)

// Step represents a single step in a pipeline
type Step struct {
	Command string            // Command name
	Variant string            // Variant (e.g., "fire" in "gradient:fire")
	Args    map[string]string // key=value arguments
	Text    string            // Text argument (for banner, etc.)
}

// Pipeline represents a parsed DSL pipeline
type Pipeline struct {
	steps []*Step
}

// Lexer tokenizes DSL input
type Lexer struct {
	input string
	pos   int
	curr  rune
}

// NewLexer creates a new lexer
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	if len(input) > 0 {
		l.curr = rune(input[0])
	}
	return l
}

// advance moves to the next character
func (l *Lexer) advance() {
	l.pos++
	if l.pos >= len(l.input) {
		l.curr = 0
	} else {
		l.curr = rune(l.input[l.pos])
	}
}

// peek looks at the next character without advancing
func (l *Lexer) peek() rune {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return rune(l.input[l.pos+1])
}

// skipWhitespace skips whitespace characters
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.curr) {
		l.advance()
	}
}

// NextToken returns the next token in the input
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	pos := l.pos

	if l.curr == 0 {
		return Token{Type: TokenEOF, Pos: pos}
	}

	// Pipe separator
	if l.curr == '|' {
		l.advance()
		return Token{Type: TokenPipe, Value: "|", Pos: pos}
	}

	// Colon separator
	if l.curr == ':' {
		l.advance()
		return Token{Type: TokenColon, Value: ":", Pos: pos}
	}

	// Equals operator
	if l.curr == '=' {
		l.advance()
		return Token{Type: TokenEquals, Value: "=", Pos: pos}
	}

	// Quoted string
	if l.curr == '\'' || l.curr == '"' {
		return l.readString(pos)
	}

	// Identifier, keyword, or number
	if unicode.IsLetter(l.curr) || l.curr == '_' || unicode.IsDigit(l.curr) {
		return l.readIdentifier(pos)
	}

	// Unknown character
	l.advance()
	return Token{Type: TokenError, Value: string(l.curr), Pos: pos}
}

// readString reads a quoted string
func (l *Lexer) readString(pos int) Token {
	quote := l.curr
	l.advance()

	var sb strings.Builder
	for l.curr != 0 && l.curr != quote {
		if l.curr == '\\' {
			l.advance()
			if l.curr == 0 {
				break
			}
			sb.WriteRune(l.curr)
		} else {
			sb.WriteRune(l.curr)
		}
		l.advance()
	}

	if l.curr == quote {
		l.advance()
	}

	return Token{Type: TokenString, Value: sb.String(), Pos: pos}
}

// readIdentifier reads an identifier
func (l *Lexer) readIdentifier(pos int) Token {
	var sb strings.Builder
	for unicode.IsLetter(l.curr) || unicode.IsDigit(l.curr) || l.curr == '_' || l.curr == '-' {
		sb.WriteRune(l.curr)
		l.advance()
	}
	return Token{Type: TokenIdentifier, Value: sb.String(), Pos: pos}
}

// Parser parses DSL tokens into a Pipeline
type Parser struct {
	lexer   *Lexer
	current Token
}

// NewParser creates a new parser
func NewParser(input string) *Parser {
	p := &Parser{lexer: NewLexer(input)}
	p.current = p.lexer.NextToken()
	return p
}

// Parse parses the DSL input and returns a Pipeline
func Parse(input string) (*Pipeline, error) {
	p := NewParser(input)
	return p.parsePipeline()
}

// parsePipeline parses a complete pipeline
func (p *Parser) parsePipeline() (*Pipeline, error) {
	pipeline := &Pipeline{}

	for p.current.Type != TokenEOF {
		step, err := p.parseStep()
		if err != nil {
			return nil, err
		}
		pipeline.steps = append(pipeline.steps, step)

		if p.current.Type == TokenPipe {
			p.current = p.lexer.NextToken() // consume pipe
		} else if p.current.Type != TokenEOF {
			return nil, fmt.Errorf("expected '|' or end of input at position %d, got %q", p.current.Pos, p.current.Value)
		}
	}

	if len(pipeline.steps) == 0 {
		return nil, fmt.Errorf("empty pipeline")
	}

	return pipeline, nil
}

// parseStep parses a single pipeline step
func (p *Parser) parseStep() (*Step, error) {
	if p.current.Type != TokenIdentifier {
		return nil, fmt.Errorf("expected command name at position %d, got %q", p.current.Pos, p.current.Value)
	}

	step := &Step{
		Command: p.current.Value,
		Args:    make(map[string]string),
	}
	p.current = p.lexer.NextToken()

	// Check for variant (colon notation)
	if p.current.Type == TokenColon {
		p.current = p.lexer.NextToken()
		if p.current.Type != TokenIdentifier {
			return nil, fmt.Errorf("expected variant name after ':' at position %d", p.current.Pos)
		}
		step.Variant = p.current.Value
		p.current = p.lexer.NextToken()
	}

	// Parse arguments
	for p.current.Type != TokenPipe && p.current.Type != TokenEOF {
		if p.current.Type == TokenString {
			// Text argument
			step.Text = p.current.Value
			p.current = p.lexer.NextToken()
		} else if p.current.Type == TokenIdentifier {
			// This might be a key=value argument
			key := p.current.Value
			p.current = p.lexer.NextToken()

			// Check if followed by equals
			if p.current.Type == TokenEquals {
				p.current = p.lexer.NextToken()

				if p.current.Type == TokenString || p.current.Type == TokenIdentifier {
					step.Args[key] = p.current.Value
					p.current = p.lexer.NextToken()
				} else {
					return nil, fmt.Errorf("expected value after '=' at position %d", p.current.Pos)
				}
			} else {
				// Identifier without equals - might be start of next step
				// Put it back by stopping here
				break
			}
		} else {
			break
		}
	}

	return step, nil
}

// Steps returns the steps in the pipeline
func (pl *Pipeline) Steps() []*Step {
	return pl.steps
}
