package lexer

type stateFn func(*Lexer) stateFn

func isAlpha(r rune) bool {
	if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
		return true
	}
	return false
}

func isNum(r rune) bool {
	if r >= '0' && r <= '9' {
		return true
	}
	return false
}

func isSpace(r rune) bool {
	if r == ' ' || r == '\t' || r == '\n' {
		return true
	}

	return false
}

func isPunct(r rune) bool {
	if r > 32 && r < 127 && !isAlpha(r) && !isNum(r) {
		return true
	}
	return false
}

func lexFunction(l *Lexer) stateFn {

	for {
		switch r := l.next(); {
		case r == 0:
			return nil

		case isNum(r):
			l.backup()
			return lexNumber(l)

		case isSpace(r):
			l.ignore()
			if r == '\n' {
				l.currRow += 1
				l.currColumn = 1
			}

		case isAlpha(r):
			l.backup()
			return lexKeywordOrId(l)

		default:
			return lexKeyword(l)
		}
	}
}

func lexNumber(l *Lexer) stateFn {

	l.acceptRun(isNum)

	if l.acceptSet(".") {
		l.acceptRun(isNum)
	}

	l.emit(tokenNamesToTokenMeta[Number])

	return lexFunction
}

// lexing non-alpha keywords
func lexKeyword(l *Lexer) stateFn {

	l.acceptRun(isPunct)

	for l.start != l.pos {

		tok, ok := reservedWordsToTokenMeta[l.currWord()]
		if !ok {
			l.backup()
			continue
		}

		l.emit(tok)
		return lexFunction
	}

	l.acceptRun(isPunct)
	l.fail()
	return nil
}

// lexing alpha keywords and variables
func lexKeywordOrId(l *Lexer) stateFn {

	l.acceptRun(isAlpha)

	tok, ok := reservedWordsToTokenMeta[l.currWord()]
	if !ok { // got variable
		if l.pos-l.start > 1 || l.pos-l.start == 0 {
			l.fail()
			return nil
		}

		l.emit(tokenNamesToTokenMeta[Variable])
		return lexFunction
	}

	l.emit(tok)
	return lexFunction
}
