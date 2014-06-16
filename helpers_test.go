package bluemonday

import (
	"testing"
)

func TestRegexpVars(t *testing.T) {
	// CellAlign
	if !CellAlign.MatchString("CENTER") {
		t.Error("CellAlign did not match: CENTER")
	}
	if !CellAlign.MatchString("justIFY") {
		t.Error("CellAlign did not match: justIFY")
	}
	if !CellAlign.MatchString("left") {
		t.Error("CellAlign did not match: left")
	}
	if !CellAlign.MatchString("right") {
		t.Error("CellAlign did not match: right")
	}
	if !CellAlign.MatchString("char") {
		t.Error("CellAlign did not match: char")
	}
	if CellAlign.MatchString("char char") {
		t.Error("CellAlign matched: char char")
	}

	// CellVerticalAlign
	if !CellVerticalAlign.MatchString("BASELINE") {
		t.Error("CellVerticalAlign did not match: BASELINE")
	}
	if !CellVerticalAlign.MatchString("boTtOM") {
		t.Error("CellVerticalAlign did not match: boTtOM")
	}
	if !CellVerticalAlign.MatchString("middle") {
		t.Error("CellVerticalAlign did not match: middle")
	}
	if !CellVerticalAlign.MatchString("top") {
		t.Error("CellVerticalAlign did not match: top")
	}
	if CellVerticalAlign.MatchString("top top") {
		t.Error("CellVerticalAlign matched: top top")
	}

	// Direction
	if !Direction.MatchString("RTL") {
		t.Error("Direction did not match: RTL")
	}
	if !Direction.MatchString("ltr") {
		t.Error("Direction did not match: ltr")
	}
	if Direction.MatchString("") {
		t.Error("Direction matched: ")
	}
	if Direction.MatchString("rtltr") {
		t.Error("Direction matched: rtltr")
	}
	if Direction.MatchString("rtl rtl") {
		t.Error("Direction matched: rtl rtl")
	}

	// ImageAlign
	if !ImageAlign.MatchString("LEFT") {
		t.Error("CellAlign did not match: LEFT")
	}
	if !ImageAlign.MatchString("right") {
		t.Error("CellAlign did not match: right")
	}
	if !ImageAlign.MatchString("tOP") {
		t.Error("CellAlign did not match: tOP")
	}
	if !ImageAlign.MatchString("texttop") {
		t.Error("CellAlign did not match: texttop")
	}
	if !ImageAlign.MatchString("middle") {
		t.Error("CellAlign did not match: middle")
	}
	if !ImageAlign.MatchString("absmiddle") {
		t.Error("CellAlign did not match: absmiddle")
	}
	if !ImageAlign.MatchString("baseline") {
		t.Error("CellAlign did not match: baseline")
	}
	if !ImageAlign.MatchString("bottom") {
		t.Error("CellAlign did not match: bottom")
	}
	if !ImageAlign.MatchString("absbottom") {
		t.Error("CellAlign did not match: absbottom")
	}
	if ImageAlign.MatchString("left right") {
		t.Error("CellAlign matched: left right")
	}
	if ImageAlign.MatchString("left left") {
		t.Error("CellAlign matched: left left")
	}
	if ImageAlign.MatchString("char") {
		t.Error("CellAlign matched: char")
	}
	if ImageAlign.MatchString("char") {
		t.Error("CellAlign matched: char")
	}

	// Integer
	if !Integer.MatchString("123") {
		t.Error("Integer did not match: 123")
	}
	if !Integer.MatchString("0") {
		t.Error("Integer did not match: 0")
	}
	if Integer.MatchString("-1") {
		t.Error("Integer matched: -1")
	}
	if Integer.MatchString("1abc") {
		t.Error("Integer matched: 1abc")
	}

	// ISO8601
	if !ISO8601.MatchString("2014") {
		t.Error("ISO8601 did not match: 2014")
	}
	if !ISO8601.MatchString("2014-02") {
		t.Error("ISO8601 did not match: 2014-02")
	}
	if !ISO8601.MatchString("2014-02-28") {
		t.Error("ISO8601 did not match: 2014-02-28")
	}
	if !ISO8601.MatchString("2014-02-28T23:59") {
		t.Error("ISO8601 did not match: 2014-02-28T23:59")
	}
	if !ISO8601.MatchString("2014-02-28T23:59-05:00") {
		t.Error("ISO8601 did not match: 2014-02-28T23:59-05:00")
	}
	if !ISO8601.MatchString("2014-02-28T23:59-05:00") {
		t.Error("ISO8601 did not match: 2014-02-28T23:59+01:00")
	}
	if !ISO8601.MatchString("2014-02-28T23:59:59") {
		t.Error("ISO8601 did not match: 2014-02-28T23:59:59")
	}
	if !ISO8601.MatchString("2014-02-28T23:59:59-05:00") {
		t.Error("ISO8601 did not match: 2014-02-28T23:59:59-05:00")
	}
	if !ISO8601.MatchString("2014-02-28T23:59:59-05:00") {
		t.Error("ISO8601 did not match: 2014-02-28T23:59:59+01:00")
	}
	if ISO8601.MatchString("201-") {
		t.Error("ISO8601 matched: 201-")
	}
	if ISO8601.MatchString("2014-0") {
		t.Error("ISO8601 matched: 2014-0")
	}
	if ISO8601.MatchString("2014-02-28-") {
		t.Error("ISO8601 matched: 2014-02-28-")
	}
	if ISO8601.MatchString("20-02-28T23:59-05:00") {
		t.Error("ISO8601 matched: 20-02-28T23:59+01:00")
	}
	if ISO8601.MatchString(" 2014-02-28T23:59:59") {
		t.Error("ISO8601 matched:  2014-02-28T23:59:59")
	}

	// ListType
	if !ListType.MatchString("CIRCLE") {
		t.Error("ListType did not match: CIRCLE")
	}
	if !ListType.MatchString("disc") {
		t.Error("ListType did not match: disc")
	}
	if !ListType.MatchString("square") {
		t.Error("ListType did not match: square")
	}
	if !ListType.MatchString("a") {
		t.Error("ListType did not match: a")
	}
	if !ListType.MatchString("A") {
		t.Error("ListType did not match: A")
	}
	if !ListType.MatchString("i") {
		t.Error("ListType did not match: i")
	}
	if !ListType.MatchString("I") {
		t.Error("ListType did not match: I")
	}
	if !ListType.MatchString("1") {
		t.Error("ListType did not match: 1")
	}
	if ListType.MatchString("circle circle") {
		t.Error("ListType matched: circle circle")
	}
	if ListType.MatchString("aa") {
		t.Error("ListType matched: aa")
	}

	// SpaceSeparatedTokens
	if !SpaceSeparatedTokens.MatchString("nofollow") {
		t.Error("SpaceSeparatedTokens did not match: nofollow")
	}
	if !SpaceSeparatedTokens.MatchString("nofollow person") {
		t.Error("SpaceSeparatedTokens did not match: nofollow person")
	}
	if !SpaceSeparatedTokens.MatchString("header") {
		t.Error("SpaceSeparatedTokens did not match: header")
	}
	if !SpaceSeparatedTokens.MatchString("bläh") {
		t.Error("SpaceSeparatedTokens did not match: bläh")
	}
	if !SpaceSeparatedTokens.MatchString("blah bläh") {
		t.Error("SpaceSeparatedTokens did not match: blah bläh")
	}
	if SpaceSeparatedTokens.MatchString("bläh blah ☃") {
		t.Error("SpaceSeparatedTokens matched: bläh blah ☃")
	}
	if SpaceSeparatedTokens.MatchString("header javascript:alert(1)") {
		t.Error("SpaceSeparatedTokens matched: header javascript:alert(1)")
	}
	if SpaceSeparatedTokens.MatchString("header &gt;") {
		t.Error("SpaceSeparatedTokens matched: header &gt;")
	}

	// Number
	if !Number.MatchString("0") {
		t.Error("Number did not match: 0")
	}
	if !Number.MatchString("1") {
		t.Error("Number did not match: 1")
	}
	if !Number.MatchString("+1") {
		t.Error("Number did not match: +1")
	}
	if !Number.MatchString("-1") {
		t.Error("Number did not match: -1")
	}
	if !Number.MatchString("1.1") {
		t.Error("Number did not match: 1.1")
	}
	if !Number.MatchString("1.2e3") {
		t.Error("Number did not match: 1.2e3")
	}
	if !Number.MatchString("7E-10") {
		t.Error("Number did not match: 7E-10")
	}
	if Number.MatchString("e7.13") {
		t.Error("Number matched: e7.13")
	}
	if Number.MatchString(`7E`) {
		t.Error(`Number matched: 7E`)
	}

	// NumberOrPercent
	if !NumberOrPercent.MatchString("0") {
		t.Error("NumberOrPercent did not match: 0")
	}
	if !NumberOrPercent.MatchString("1") {
		t.Error("NumberOrPercent did not match: 1")
	}
	if !NumberOrPercent.MatchString(`1%`) {
		t.Error(`NumberOrPercent did not match: 1(percent)`)
	}
	if NumberOrPercent.MatchString("-1") {
		t.Error("NumberOrPercent matched: -1")
	}
	if NumberOrPercent.MatchString("1.1") {
		t.Error("NumberOrPercent matched: 1.1")
	}
	if NumberOrPercent.MatchString("1.2e3") {
		t.Error("NumberOrPercent matched: 1.2e3")
	}
	if NumberOrPercent.MatchString("7E-10") {
		t.Error("NumberOrPercent matched: 7E-10")
	}
	if NumberOrPercent.MatchString("e7.13") {
		t.Error("NumberOrPercent matched: e7.13")
	}
	if NumberOrPercent.MatchString(`7E`) {
		t.Error(`NumberOrPercent matched: 7E`)
	}

	// Paragraph
	if !Paragraph.MatchString("hello world") {
		t.Error("Paragraph did not match: hello world")
	}
	if !Paragraph.MatchString("blah bläh blah") {
		t.Error("Paragraph did not match: blah bläh blah")
	}
	if Paragraph.MatchString("bläh blah ☃") {
		t.Error("Paragraph matched: bläh blah ☃")
	}
	if Paragraph.MatchString("javascript:alert(1)") {
		t.Error("Paragraph matched: javascript:alert(1)")
	}
}
