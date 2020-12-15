package phonenumbers

import(
    "testing"
)

func assertEquals(t *testing.T, expected string, got string) {
    if expected != got {
        t.Errorf("%s != %s", expected, got)
    }
}

func assertEqualsInt(t *testing.T, expected int, got int) {
    if expected != got {
        t.Errorf("%d != %d", expected, got)
    }
}
  
func TestInvalidRegion(t *testing.T) {
    formatter := NewAsYouTypeFormatter("ZZ")
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+4", formatter.InputDigit('4'))
    assertEquals(t, "+48 ", formatter.InputDigit('8'))
    assertEquals(t, "+48 8", formatter.InputDigit('8'))
    assertEquals(t, "+48 88", formatter.InputDigit('8'))
    assertEquals(t, "+48 88 1", formatter.InputDigit('1'))
    assertEquals(t, "+48 88 12", formatter.InputDigit('2'))
    assertEquals(t, "+48 88 123", formatter.InputDigit('3'))
    assertEquals(t, "+48 88 123 1", formatter.InputDigit('1'))
    assertEquals(t, "+48 88 123 12", formatter.InputDigit('2'))

    formatter.Clear()
    assertEquals(t, "6", formatter.InputDigit('6'))
    assertEquals(t, "65", formatter.InputDigit('5'))
    assertEquals(t, "650", formatter.InputDigit('0'))
    assertEquals(t, "6502", formatter.InputDigit('2'))
    assertEquals(t, "65025", formatter.InputDigit('5'))
    assertEquals(t, "650253", formatter.InputDigit('3'))
}

func TestInvalidPlusSign(t *testing.T) {
    formatter := NewAsYouTypeFormatter("ZZ")
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+4", formatter.InputDigit('4'))
    assertEquals(t, "+48 ", formatter.InputDigit('8'))
    assertEquals(t, "+48 8", formatter.InputDigit('8'))
    assertEquals(t, "+48 88", formatter.InputDigit('8'))
    assertEquals(t, "+48 88 1", formatter.InputDigit('1'))
    assertEquals(t, "+48 88 12", formatter.InputDigit('2'))
    assertEquals(t, "+48 88 123", formatter.InputDigit('3'))
    assertEquals(t, "+48 88 123 1", formatter.InputDigit('1'))
    // A plus sign can only appear at the beginning of the number otherwise, no formatting is
    // applied.
    assertEquals(t, "+48881231+", formatter.InputDigit('+'))
    assertEquals(t, "+48881231+2", formatter.InputDigit('2'))
}

func TestTooLongNumberMatchingMultipleLeadingDigits(t *testing.T) {
    // See https://github.com/google/libphonenumber/issues/36
    // The bug occurred last time for countries which have two formatting rules with exactly the
    // same leading digits pattern but differ in length.
    formatter := NewAsYouTypeFormatter("ZZ")
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+8", formatter.InputDigit('8'))
    assertEquals(t, "+81 ", formatter.InputDigit('1'))
    assertEquals(t, "+81 9", formatter.InputDigit('9'))
    assertEquals(t, "+81 90", formatter.InputDigit('0'))
    assertEquals(t, "+81 90 1", formatter.InputDigit('1'))
    assertEquals(t, "+81 90 12", formatter.InputDigit('2'))
    assertEquals(t, "+81 90 123", formatter.InputDigit('3'))
    assertEquals(t, "+81 90 1234", formatter.InputDigit('4'))
    assertEquals(t, "+81 90 1234 5", formatter.InputDigit('5'))
    assertEquals(t, "+81 90 1234 56", formatter.InputDigit('6'))
    assertEquals(t, "+81 90 1234 567", formatter.InputDigit('7'))
    assertEquals(t, "+81 90 1234 5678", formatter.InputDigit('8'))
    assertEquals(t, "+81 90 12 345 6789", formatter.InputDigit('9'))
    assertEquals(t, "+81901234567890", formatter.InputDigit('0'))
    assertEquals(t, "+819012345678901", formatter.InputDigit('1'))
}

func TestCountryWithSpaceInNationalPrefixFormattingRule(t *testing.T) {
    formatter := NewAsYouTypeFormatter("BY")
    assertEquals(t, "8", formatter.InputDigit('8'))
    assertEquals(t, "88", formatter.InputDigit('8'))
    assertEquals(t, "881", formatter.InputDigit('1'))
    assertEquals(t, "8 819", formatter.InputDigit('9'))
    assertEquals(t, "8 8190", formatter.InputDigit('0'))
    // The formatting rule for 5 digit numbers states that no space should be present after the
    // national prefix.
    assertEquals(t, "881 901", formatter.InputDigit('1'))
    assertEquals(t, "8 819 012", formatter.InputDigit('2'))
    // Too long, no formatting rule applies.
    assertEquals(t, "88190123", formatter.InputDigit('3'))
}

func TestCountryWithSpaceInNationalPrefixFormattingRuleAndLongNdd(t *testing.T) {
    formatter := NewAsYouTypeFormatter("BY")
    assertEquals(t, "9", formatter.InputDigit('9'))
    assertEquals(t, "99", formatter.InputDigit('9'))
    assertEquals(t, "999", formatter.InputDigit('9'))
    assertEquals(t, "9999", formatter.InputDigit('9'))
    assertEquals(t, "99999 ", formatter.InputDigit('9'))
    assertEquals(t, "99999 1", formatter.InputDigit('1'))
    assertEquals(t, "99999 12", formatter.InputDigit('2'))
    assertEquals(t, "99999 123", formatter.InputDigit('3'))
    assertEquals(t, "99999 1234", formatter.InputDigit('4'))
    assertEquals(t, "99999 12 345", formatter.InputDigit('5'))
}

func TestAYTFUS(t *testing.T) {
    formatter := NewAsYouTypeFormatter("US")
    assertEquals(t, "6", formatter.InputDigit('6'))
    assertEquals(t, "65", formatter.InputDigit('5'))
    assertEquals(t, "650", formatter.InputDigit('0'))
    assertEquals(t, "650 2", formatter.InputDigit('2'))
    assertEquals(t, "650 25", formatter.InputDigit('5'))
    assertEquals(t, "650 253", formatter.InputDigit('3'))
    // Note this is how a US local number (without area code) should be formatted.
    assertEquals(t, "650 2532", formatter.InputDigit('2'))
    assertEquals(t, "650 253 22", formatter.InputDigit('2'))
    assertEquals(t, "650 253 222", formatter.InputDigit('2'))
    assertEquals(t, "650 253 2222", formatter.InputDigit('2'))

    formatter.Clear()
    assertEquals(t, "1", formatter.InputDigit('1'))
    assertEquals(t, "16", formatter.InputDigit('6'))
    assertEquals(t, "1 65", formatter.InputDigit('5'))
    assertEquals(t, "1 650", formatter.InputDigit('0'))
    assertEquals(t, "1 650 2", formatter.InputDigit('2'))
    assertEquals(t, "1 650 25", formatter.InputDigit('5'))
    assertEquals(t, "1 650 253", formatter.InputDigit('3'))
    assertEquals(t, "1 650 253 2", formatter.InputDigit('2'))
    assertEquals(t, "1 650 253 22", formatter.InputDigit('2'))
    assertEquals(t, "1 650 253 222", formatter.InputDigit('2'))
    assertEquals(t, "1 650 253 2222", formatter.InputDigit('2'))

    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "01", formatter.InputDigit('1'))
    assertEquals(t, "011 ", formatter.InputDigit('1'))
    assertEquals(t, "011 4", formatter.InputDigit('4'))
    assertEquals(t, "011 44 ", formatter.InputDigit('4'))
    assertEquals(t, "011 44 6", formatter.InputDigit('6'))
    assertEquals(t, "011 44 61", formatter.InputDigit('1'))
    assertEquals(t, "011 44 6 12", formatter.InputDigit('2'))
    assertEquals(t, "011 44 6 123", formatter.InputDigit('3'))
    assertEquals(t, "011 44 6 123 1", formatter.InputDigit('1'))
    assertEquals(t, "011 44 6 123 12", formatter.InputDigit('2'))
    assertEquals(t, "011 44 6 123 123", formatter.InputDigit('3'))
    assertEquals(t, "011 44 6 123 123 1", formatter.InputDigit('1'))
    assertEquals(t, "011 44 6 123 123 12", formatter.InputDigit('2'))
    assertEquals(t, "011 44 6 123 123 123", formatter.InputDigit('3'))

    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "01", formatter.InputDigit('1'))
    assertEquals(t, "011 ", formatter.InputDigit('1'))
    assertEquals(t, "011 5", formatter.InputDigit('5'))
    assertEquals(t, "011 54 ", formatter.InputDigit('4'))
    assertEquals(t, "011 54 9", formatter.InputDigit('9'))
    assertEquals(t, "011 54 91", formatter.InputDigit('1'))
    assertEquals(t, "011 54 9 11", formatter.InputDigit('1'))
    assertEquals(t, "011 54 9 11 2", formatter.InputDigit('2'))
    assertEquals(t, "011 54 9 11 23", formatter.InputDigit('3'))
    assertEquals(t, "011 54 9 11 231", formatter.InputDigit('1'))
    assertEquals(t, "011 54 9 11 2312", formatter.InputDigit('2'))
    assertEquals(t, "011 54 9 11 2312 1", formatter.InputDigit('1'))
    assertEquals(t, "011 54 9 11 2312 12", formatter.InputDigit('2'))
    assertEquals(t, "011 54 9 11 2312 123", formatter.InputDigit('3'))
    assertEquals(t, "011 54 9 11 2312 1234", formatter.InputDigit('4'))

    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "01", formatter.InputDigit('1'))
    assertEquals(t, "011 ", formatter.InputDigit('1'))
    assertEquals(t, "011 2", formatter.InputDigit('2'))
    assertEquals(t, "011 24", formatter.InputDigit('4'))
    assertEquals(t, "011 244 ", formatter.InputDigit('4'))
    assertEquals(t, "011 244 2", formatter.InputDigit('2'))
    assertEquals(t, "011 244 28", formatter.InputDigit('8'))
    assertEquals(t, "011 244 280", formatter.InputDigit('0'))
    assertEquals(t, "011 244 280 0", formatter.InputDigit('0'))
    assertEquals(t, "011 244 280 00", formatter.InputDigit('0'))
    assertEquals(t, "011 244 280 000", formatter.InputDigit('0'))
    assertEquals(t, "011 244 280 000 0", formatter.InputDigit('0'))
    assertEquals(t, "011 244 280 000 00", formatter.InputDigit('0'))
    assertEquals(t, "011 244 280 000 000", formatter.InputDigit('0'))

    formatter.Clear()
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+4", formatter.InputDigit('4'))
    assertEquals(t, "+48 ", formatter.InputDigit('8'))
    assertEquals(t, "+48 8", formatter.InputDigit('8'))
    assertEquals(t, "+48 88", formatter.InputDigit('8'))
    assertEquals(t, "+48 88 1", formatter.InputDigit('1'))
    assertEquals(t, "+48 88 12", formatter.InputDigit('2'))
    assertEquals(t, "+48 88 123", formatter.InputDigit('3'))
    assertEquals(t, "+48 88 123 1", formatter.InputDigit('1'))
    assertEquals(t, "+48 88 123 12", formatter.InputDigit('2'))
    assertEquals(t, "+48 88 123 12 1", formatter.InputDigit('1'))
    assertEquals(t, "+48 88 123 12 12", formatter.InputDigit('2'))
}

func TestAYTFUSFullWidthCharacters(t *testing.T) {
    formatter := NewAsYouTypeFormatter("US")
    assertEquals(t, "\uFF16", formatter.InputDigit('\uFF16'))
    assertEquals(t, "\uFF16\uFF15", formatter.InputDigit('\uFF15'))
    assertEquals(t, "650", formatter.InputDigit('\uFF10'))
    assertEquals(t, "650 2", formatter.InputDigit('\uFF12'))
    assertEquals(t, "650 25", formatter.InputDigit('\uFF15'))
    assertEquals(t, "650 253", formatter.InputDigit('\uFF13'))
    assertEquals(t, "650 2532", formatter.InputDigit('\uFF12'))
    assertEquals(t, "650 253 22", formatter.InputDigit('\uFF12'))
    assertEquals(t, "650 253 222", formatter.InputDigit('\uFF12'))
    assertEquals(t, "650 253 2222", formatter.InputDigit('\uFF12'))
}

func TestAYTFUSMobileShortCode(t *testing.T) {
    formatter := NewAsYouTypeFormatter("US")
    assertEquals(t, "*", formatter.InputDigit('*'))
    assertEquals(t, "*1", formatter.InputDigit('1'))
    assertEquals(t, "*12", formatter.InputDigit('2'))
    assertEquals(t, "*121", formatter.InputDigit('1'))
    assertEquals(t, "*121#", formatter.InputDigit('#'))
}

func TestAYTFUSVanityNumber(t *testing.T) {
    formatter := NewAsYouTypeFormatter("US")
    assertEquals(t, "8", formatter.InputDigit('8'))
    assertEquals(t, "80", formatter.InputDigit('0'))
    assertEquals(t, "800", formatter.InputDigit('0'))
    assertEquals(t, "800 ", formatter.InputDigit(' '))
    assertEquals(t, "800 M", formatter.InputDigit('M'))
    assertEquals(t, "800 MY", formatter.InputDigit('Y'))
    assertEquals(t, "800 MY ", formatter.InputDigit(' '))
    assertEquals(t, "800 MY A", formatter.InputDigit('A'))
    assertEquals(t, "800 MY AP", formatter.InputDigit('P'))
    assertEquals(t, "800 MY APP", formatter.InputDigit('P'))
    assertEquals(t, "800 MY APPL", formatter.InputDigit('L'))
    assertEquals(t, "800 MY APPLE", formatter.InputDigit('E'))
}

func TestAYTFAndRememberPositionUS(t *testing.T) {
    formatter := NewAsYouTypeFormatter("US")
    assertEquals(t, "1", formatter.InputDigitAndRememberPosition('1'))
    assertEqualsInt(t, 1, formatter.GetRememberedPosition())
    assertEquals(t, "16", formatter.InputDigit('6'))
    assertEquals(t, "1 65", formatter.InputDigit('5'))
    assertEqualsInt(t, 1, formatter.GetRememberedPosition())
    assertEquals(t, "1 650", formatter.InputDigitAndRememberPosition('0'))
    assertEqualsInt(t, 5, formatter.GetRememberedPosition())
    assertEquals(t, "1 650 2", formatter.InputDigit('2'))
    assertEquals(t, "1 650 25", formatter.InputDigit('5'))
    // Note the remembered position for digit "0" changes from 4 to 5, because a space is now
    // inserted in the front.
    assertEqualsInt(t, 5, formatter.GetRememberedPosition())
    assertEquals(t, "1 650 253", formatter.InputDigit('3'))
    assertEquals(t, "1 650 253 2", formatter.InputDigit('2'))
    assertEquals(t, "1 650 253 22", formatter.InputDigit('2'))
    assertEqualsInt(t, 5, formatter.GetRememberedPosition())
    assertEquals(t, "1 650 253 222", formatter.InputDigitAndRememberPosition('2'))
    assertEqualsInt(t, 13, formatter.GetRememberedPosition())
    assertEquals(t, "1 650 253 2222", formatter.InputDigit('2'))
    assertEqualsInt(t, 13, formatter.GetRememberedPosition())
    assertEquals(t, "165025322222", formatter.InputDigit('2'))
    assertEqualsInt(t, 10, formatter.GetRememberedPosition())
    assertEquals(t, "1650253222222", formatter.InputDigit('2'))
    assertEqualsInt(t, 10, formatter.GetRememberedPosition())

    formatter.Clear()
    assertEquals(t, "1", formatter.InputDigit('1'))
    assertEquals(t, "16", formatter.InputDigitAndRememberPosition('6'))
    assertEqualsInt(t, 2, formatter.GetRememberedPosition())
    assertEquals(t, "1 65", formatter.InputDigit('5'))
    assertEquals(t, "1 650", formatter.InputDigit('0'))
    assertEqualsInt(t, 3, formatter.GetRememberedPosition())
    assertEquals(t, "1 650 2", formatter.InputDigit('2'))
    assertEquals(t, "1 650 25", formatter.InputDigit('5'))
    assertEqualsInt(t, 3, formatter.GetRememberedPosition())
    assertEquals(t, "1 650 253", formatter.InputDigit('3'))
    assertEquals(t, "1 650 253 2", formatter.InputDigit('2'))
    assertEquals(t, "1 650 253 22", formatter.InputDigit('2'))
    assertEqualsInt(t, 3, formatter.GetRememberedPosition())
    assertEquals(t, "1 650 253 222", formatter.InputDigit('2'))
    assertEquals(t, "1 650 253 2222", formatter.InputDigit('2'))
    assertEquals(t, "165025322222", formatter.InputDigit('2'))
    assertEqualsInt(t, 2, formatter.GetRememberedPosition())
    assertEquals(t, "1650253222222", formatter.InputDigit('2'))
    assertEqualsInt(t, 2, formatter.GetRememberedPosition())

    formatter.Clear()
    assertEquals(t, "6", formatter.InputDigit('6'))
    assertEquals(t, "65", formatter.InputDigit('5'))
    assertEquals(t, "650", formatter.InputDigit('0'))
    assertEquals(t, "650 2", formatter.InputDigit('2'))
    assertEquals(t, "650 25", formatter.InputDigit('5'))
    assertEquals(t, "650 253", formatter.InputDigit('3'))
    assertEquals(t, "650 2532", formatter.InputDigitAndRememberPosition('2'))
    assertEqualsInt(t, 8, formatter.GetRememberedPosition())
    assertEquals(t, "650 253 22", formatter.InputDigit('2'))
    assertEqualsInt(t, 9, formatter.GetRememberedPosition())
    assertEquals(t, "650 253 222", formatter.InputDigit('2'))
    // No more formatting when semicolon is entered.
    assertEquals(t, "650253222", formatter.InputDigit(';'))
    assertEqualsInt(t, 7, formatter.GetRememberedPosition())
    assertEquals(t, "6502532222", formatter.InputDigit('2'))

    formatter.Clear()
    assertEquals(t, "6", formatter.InputDigit('6'))
    assertEquals(t, "65", formatter.InputDigit('5'))
    assertEquals(t, "650", formatter.InputDigit('0'))
    // No more formatting when users choose to do their own formatting.
    assertEquals(t, "650-", formatter.InputDigit('-'))
    assertEquals(t, "650-2", formatter.InputDigitAndRememberPosition('2'))
    assertEqualsInt(t, 5, formatter.GetRememberedPosition())
    assertEquals(t, "650-25", formatter.InputDigit('5'))
    assertEqualsInt(t, 5, formatter.GetRememberedPosition())
    assertEquals(t, "650-253", formatter.InputDigit('3'))
    assertEqualsInt(t, 5, formatter.GetRememberedPosition())
    assertEquals(t, "650-253-", formatter.InputDigit('-'))
    assertEquals(t, "650-253-2", formatter.InputDigit('2'))
    assertEquals(t, "650-253-22", formatter.InputDigit('2'))
    assertEquals(t, "650-253-222", formatter.InputDigit('2'))
    assertEquals(t, "650-253-2222", formatter.InputDigit('2'))

    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "01", formatter.InputDigit('1'))
    assertEquals(t, "011 ", formatter.InputDigit('1'))
    assertEquals(t, "011 4", formatter.InputDigitAndRememberPosition('4'))
    assertEquals(t, "011 48 ", formatter.InputDigit('8'))
    assertEqualsInt(t, 5, formatter.GetRememberedPosition())
    assertEquals(t, "011 48 8", formatter.InputDigit('8'))
    assertEqualsInt(t, 5, formatter.GetRememberedPosition())
    assertEquals(t, "011 48 88", formatter.InputDigit('8'))
    assertEquals(t, "011 48 88 1", formatter.InputDigit('1'))
    assertEquals(t, "011 48 88 12", formatter.InputDigit('2'))
    assertEqualsInt(t, 5, formatter.GetRememberedPosition())
    assertEquals(t, "011 48 88 123", formatter.InputDigit('3'))
    assertEquals(t, "011 48 88 123 1", formatter.InputDigit('1'))
    assertEquals(t, "011 48 88 123 12", formatter.InputDigit('2'))
    assertEquals(t, "011 48 88 123 12 1", formatter.InputDigit('1'))
    assertEquals(t, "011 48 88 123 12 12", formatter.InputDigit('2'))

    formatter.Clear()
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+1", formatter.InputDigit('1'))
    assertEquals(t, "+1 6", formatter.InputDigitAndRememberPosition('6'))
    assertEquals(t, "+1 65", formatter.InputDigit('5'))
    assertEquals(t, "+1 650", formatter.InputDigit('0'))
    assertEqualsInt(t, 4, formatter.GetRememberedPosition())
    assertEquals(t, "+1 650 2", formatter.InputDigit('2'))
    assertEqualsInt(t, 4, formatter.GetRememberedPosition())
    assertEquals(t, "+1 650 25", formatter.InputDigit('5'))
    assertEquals(t, "+1 650 253", formatter.InputDigitAndRememberPosition('3'))
    assertEquals(t, "+1 650 253 2", formatter.InputDigit('2'))
    assertEquals(t, "+1 650 253 22", formatter.InputDigit('2'))
    assertEquals(t, "+1 650 253 222", formatter.InputDigit('2'))
    assertEqualsInt(t, 10, formatter.GetRememberedPosition())

    formatter.Clear()
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+1", formatter.InputDigit('1'))
    assertEquals(t, "+1 6", formatter.InputDigitAndRememberPosition('6'))
    assertEquals(t, "+1 65", formatter.InputDigit('5'))
    assertEquals(t, "+1 650", formatter.InputDigit('0'))
    assertEqualsInt(t, 4, formatter.GetRememberedPosition())
    assertEquals(t, "+1 650 2", formatter.InputDigit('2'))
    assertEqualsInt(t, 4, formatter.GetRememberedPosition())
    assertEquals(t, "+1 650 25", formatter.InputDigit('5'))
    assertEquals(t, "+1 650 253", formatter.InputDigit('3'))
    assertEquals(t, "+1 650 253 2", formatter.InputDigit('2'))
    assertEquals(t, "+1 650 253 22", formatter.InputDigit('2'))
    assertEquals(t, "+1 650 253 222", formatter.InputDigit('2'))
    assertEquals(t, "+1650253222", formatter.InputDigit(';'))
    assertEqualsInt(t, 3, formatter.GetRememberedPosition())
}

func TestAYTFGBFixedLine(t *testing.T) {
    formatter := NewAsYouTypeFormatter("GB")
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "02", formatter.InputDigit('2'))
    assertEquals(t, "020", formatter.InputDigit('0'))
    assertEquals(t, "020 7", formatter.InputDigitAndRememberPosition('7'))
    assertEqualsInt(t, 5, formatter.GetRememberedPosition())
    assertEquals(t, "020 70", formatter.InputDigit('0'))
    assertEquals(t, "020 703", formatter.InputDigit('3'))
    assertEqualsInt(t, 5, formatter.GetRememberedPosition())
    assertEquals(t, "020 7031", formatter.InputDigit('1'))
    assertEquals(t, "020 7031 3", formatter.InputDigit('3'))
    assertEquals(t, "020 7031 30", formatter.InputDigit('0'))
    assertEquals(t, "020 7031 300", formatter.InputDigit('0'))
    assertEquals(t, "020 7031 3000", formatter.InputDigit('0'))
}

func TestAYTFGBTollFree(t *testing.T) {
    formatter := NewAsYouTypeFormatter("GB")
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "08", formatter.InputDigit('8'))
    assertEquals(t, "080", formatter.InputDigit('0'))
    assertEquals(t, "080 7", formatter.InputDigit('7'))
    assertEquals(t, "080 70", formatter.InputDigit('0'))
    assertEquals(t, "080 703", formatter.InputDigit('3'))
    assertEquals(t, "080 7031", formatter.InputDigit('1'))
    assertEquals(t, "080 7031 3", formatter.InputDigit('3'))
    assertEquals(t, "080 7031 30", formatter.InputDigit('0'))
    assertEquals(t, "080 7031 300", formatter.InputDigit('0'))
    assertEquals(t, "080 7031 3000", formatter.InputDigit('0'))
}

func TestAYTFGBPremiumRate(t *testing.T) {
    formatter := NewAsYouTypeFormatter("GB")
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "09", formatter.InputDigit('9'))
    assertEquals(t, "090", formatter.InputDigit('0'))
    assertEquals(t, "090 7", formatter.InputDigit('7'))
    assertEquals(t, "090 70", formatter.InputDigit('0'))
    assertEquals(t, "090 703", formatter.InputDigit('3'))
    assertEquals(t, "090 7031", formatter.InputDigit('1'))
    assertEquals(t, "090 7031 3", formatter.InputDigit('3'))
    assertEquals(t, "090 7031 30", formatter.InputDigit('0'))
    assertEquals(t, "090 7031 300", formatter.InputDigit('0'))
    assertEquals(t, "090 7031 3000", formatter.InputDigit('0'))
}

func TestAYTFNZMobile(t *testing.T) {
    formatter := NewAsYouTypeFormatter("NZ")
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "02", formatter.InputDigit('2'))
    assertEquals(t, "021", formatter.InputDigit('1'))
    assertEquals(t, "02-11", formatter.InputDigit('1'))
    assertEquals(t, "02-112", formatter.InputDigit('2'))
    // Note the unittest is using fake metadata which might produce non-ideal results.
    assertEquals(t, "02-112 3", formatter.InputDigit('3'))
    assertEquals(t, "02-112 34", formatter.InputDigit('4'))
    assertEquals(t, "02-112 345", formatter.InputDigit('5'))
    assertEquals(t, "02-112 3456", formatter.InputDigit('6'))
}

func TestAYTFDE(t *testing.T) {
    formatter := NewAsYouTypeFormatter("DE")
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "03", formatter.InputDigit('3'))
    assertEquals(t, "030", formatter.InputDigit('0'))
    assertEquals(t, "030/1", formatter.InputDigit('1'))
    assertEquals(t, "030/12", formatter.InputDigit('2'))
    assertEquals(t, "030/123", formatter.InputDigit('3'))
    assertEquals(t, "030/1234", formatter.InputDigit('4'))

    // 04134 1234
    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "04", formatter.InputDigit('4'))
    assertEquals(t, "041", formatter.InputDigit('1'))
    assertEquals(t, "041 3", formatter.InputDigit('3'))
    assertEquals(t, "041 34", formatter.InputDigit('4'))
    assertEquals(t, "04134 1", formatter.InputDigit('1'))
    assertEquals(t, "04134 12", formatter.InputDigit('2'))
    assertEquals(t, "04134 123", formatter.InputDigit('3'))
    assertEquals(t, "04134 1234", formatter.InputDigit('4'))

    // 08021 2345
    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "08", formatter.InputDigit('8'))
    assertEquals(t, "080", formatter.InputDigit('0'))
    assertEquals(t, "080 2", formatter.InputDigit('2'))
    assertEquals(t, "080 21", formatter.InputDigit('1'))
    assertEquals(t, "08021 2", formatter.InputDigit('2'))
    assertEquals(t, "08021 23", formatter.InputDigit('3'))
    assertEquals(t, "08021 234", formatter.InputDigit('4'))
    assertEquals(t, "08021 2345", formatter.InputDigit('5'))

    // 00 1 650 253 2250
    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "00", formatter.InputDigit('0'))
    assertEquals(t, "00 1 ", formatter.InputDigit('1'))
    assertEquals(t, "00 1 6", formatter.InputDigit('6'))
    assertEquals(t, "00 1 65", formatter.InputDigit('5'))
    assertEquals(t, "00 1 650", formatter.InputDigit('0'))
    assertEquals(t, "00 1 650 2", formatter.InputDigit('2'))
    assertEquals(t, "00 1 650 25", formatter.InputDigit('5'))
    assertEquals(t, "00 1 650 253", formatter.InputDigit('3'))
    assertEquals(t, "00 1 650 253 2", formatter.InputDigit('2'))
    assertEquals(t, "00 1 650 253 22", formatter.InputDigit('2'))
    assertEquals(t, "00 1 650 253 222", formatter.InputDigit('2'))
    assertEquals(t, "00 1 650 253 2222", formatter.InputDigit('2'))
}

func TestAYTFAR(t *testing.T) {
    formatter := NewAsYouTypeFormatter("AR")
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "01", formatter.InputDigit('1'))
    assertEquals(t, "011", formatter.InputDigit('1'))
    assertEquals(t, "011 7", formatter.InputDigit('7'))
    assertEquals(t, "011 70", formatter.InputDigit('0'))
    assertEquals(t, "011 703", formatter.InputDigit('3'))
    assertEquals(t, "011 7031", formatter.InputDigit('1'))
    assertEquals(t, "011 7031-3", formatter.InputDigit('3'))
    assertEquals(t, "011 7031-30", formatter.InputDigit('0'))
    assertEquals(t, "011 7031-300", formatter.InputDigit('0'))
    assertEquals(t, "011 7031-3000", formatter.InputDigit('0'))
}

func TestAYTFARMobile(t *testing.T) {
    formatter := NewAsYouTypeFormatter("AR")
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+5", formatter.InputDigit('5'))
    assertEquals(t, "+54 ", formatter.InputDigit('4'))
    assertEquals(t, "+54 9", formatter.InputDigit('9'))
    assertEquals(t, "+54 91", formatter.InputDigit('1'))
    assertEquals(t, "+54 9 11", formatter.InputDigit('1'))
    assertEquals(t, "+54 9 11 2", formatter.InputDigit('2'))
    assertEquals(t, "+54 9 11 23", formatter.InputDigit('3'))
    assertEquals(t, "+54 9 11 231", formatter.InputDigit('1'))
    assertEquals(t, "+54 9 11 2312", formatter.InputDigit('2'))
    assertEquals(t, "+54 9 11 2312 1", formatter.InputDigit('1'))
    assertEquals(t, "+54 9 11 2312 12", formatter.InputDigit('2'))
    assertEquals(t, "+54 9 11 2312 123", formatter.InputDigit('3'))
    assertEquals(t, "+54 9 11 2312 1234", formatter.InputDigit('4'))
}

func TestAYTFKR(t *testing.T) {
    // +82 51 234 5678
    formatter := NewAsYouTypeFormatter("KR")
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+8", formatter.InputDigit('8'))
    assertEquals(t, "+82 ", formatter.InputDigit('2'))
    assertEquals(t, "+82 5", formatter.InputDigit('5'))
    assertEquals(t, "+82 51", formatter.InputDigit('1'))
    assertEquals(t, "+82 51-2", formatter.InputDigit('2'))
    assertEquals(t, "+82 51-23", formatter.InputDigit('3'))
    assertEquals(t, "+82 51-234", formatter.InputDigit('4'))
    assertEquals(t, "+82 51-234-5", formatter.InputDigit('5'))
    assertEquals(t, "+82 51-234-56", formatter.InputDigit('6'))
    assertEquals(t, "+82 51-234-567", formatter.InputDigit('7'))
    assertEquals(t, "+82 51-234-5678", formatter.InputDigit('8'))

    // +82 2 531 5678
    formatter.Clear()
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+8", formatter.InputDigit('8'))
    assertEquals(t, "+82 ", formatter.InputDigit('2'))
    assertEquals(t, "+82 2", formatter.InputDigit('2'))
    assertEquals(t, "+82 25", formatter.InputDigit('5'))
    assertEquals(t, "+82 2-53", formatter.InputDigit('3'))
    assertEquals(t, "+82 2-531", formatter.InputDigit('1'))
    assertEquals(t, "+82 2-531-5", formatter.InputDigit('5'))
    assertEquals(t, "+82 2-531-56", formatter.InputDigit('6'))
    assertEquals(t, "+82 2-531-567", formatter.InputDigit('7'))
    assertEquals(t, "+82 2-531-5678", formatter.InputDigit('8'))

    // +82 2 3665 5678
    formatter.Clear()
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+8", formatter.InputDigit('8'))
    assertEquals(t, "+82 ", formatter.InputDigit('2'))
    assertEquals(t, "+82 2", formatter.InputDigit('2'))
    assertEquals(t, "+82 23", formatter.InputDigit('3'))
    assertEquals(t, "+82 2-36", formatter.InputDigit('6'))
    assertEquals(t, "+82 2-366", formatter.InputDigit('6'))
    assertEquals(t, "+82 2-3665", formatter.InputDigit('5'))
    assertEquals(t, "+82 2-3665-5", formatter.InputDigit('5'))
    assertEquals(t, "+82 2-3665-56", formatter.InputDigit('6'))
    assertEquals(t, "+82 2-3665-567", formatter.InputDigit('7'))
    assertEquals(t, "+82 2-3665-5678", formatter.InputDigit('8'))

    // 02-114
    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "02", formatter.InputDigit('2'))
    assertEquals(t, "021", formatter.InputDigit('1'))
    assertEquals(t, "02-11", formatter.InputDigit('1'))
    assertEquals(t, "02-114", formatter.InputDigit('4'))

    // 02-1300
    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "02", formatter.InputDigit('2'))
    assertEquals(t, "021", formatter.InputDigit('1'))
    assertEquals(t, "02-13", formatter.InputDigit('3'))
    assertEquals(t, "02-130", formatter.InputDigit('0'))
    assertEquals(t, "02-1300", formatter.InputDigit('0'))

    // 011-456-7890
    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "01", formatter.InputDigit('1'))
    assertEquals(t, "011", formatter.InputDigit('1'))
    assertEquals(t, "011-4", formatter.InputDigit('4'))
    assertEquals(t, "011-45", formatter.InputDigit('5'))
    assertEquals(t, "011-456", formatter.InputDigit('6'))
    assertEquals(t, "011-456-7", formatter.InputDigit('7'))
    assertEquals(t, "011-456-78", formatter.InputDigit('8'))
    assertEquals(t, "011-456-789", formatter.InputDigit('9'))
    assertEquals(t, "011-456-7890", formatter.InputDigit('0'))

    // 011-9876-7890
    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "01", formatter.InputDigit('1'))
    assertEquals(t, "011", formatter.InputDigit('1'))
    assertEquals(t, "011-9", formatter.InputDigit('9'))
    assertEquals(t, "011-98", formatter.InputDigit('8'))
    assertEquals(t, "011-987", formatter.InputDigit('7'))
    assertEquals(t, "011-9876", formatter.InputDigit('6'))
    assertEquals(t, "011-9876-7", formatter.InputDigit('7'))
    assertEquals(t, "011-9876-78", formatter.InputDigit('8'))
    assertEquals(t, "011-9876-789", formatter.InputDigit('9'))
    assertEquals(t, "011-9876-7890", formatter.InputDigit('0'))
}

func TestAYTF_MX(t *testing.T) {
    formatter := NewAsYouTypeFormatter("MX")

    // +52 800 123 4567
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+5", formatter.InputDigit('5'))
    assertEquals(t, "+52 ", formatter.InputDigit('2'))
    assertEquals(t, "+52 8", formatter.InputDigit('8'))
    assertEquals(t, "+52 80", formatter.InputDigit('0'))
    assertEquals(t, "+52 800", formatter.InputDigit('0'))
    assertEquals(t, "+52 800 1", formatter.InputDigit('1'))
    assertEquals(t, "+52 800 12", formatter.InputDigit('2'))
    assertEquals(t, "+52 800 123", formatter.InputDigit('3'))
    assertEquals(t, "+52 800 123 4", formatter.InputDigit('4'))
    assertEquals(t, "+52 800 123 45", formatter.InputDigit('5'))
    assertEquals(t, "+52 800 123 456", formatter.InputDigit('6'))
    assertEquals(t, "+52 800 123 4567", formatter.InputDigit('7'))

    // +52 55 1234 5678
    formatter.Clear()
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+5", formatter.InputDigit('5'))
    assertEquals(t, "+52 ", formatter.InputDigit('2'))
    assertEquals(t, "+52 5", formatter.InputDigit('5'))
    assertEquals(t, "+52 55", formatter.InputDigit('5'))
    assertEquals(t, "+52 55 1", formatter.InputDigit('1'))
    assertEquals(t, "+52 55 12", formatter.InputDigit('2'))
    assertEquals(t, "+52 55 123", formatter.InputDigit('3'))
    assertEquals(t, "+52 55 1234", formatter.InputDigit('4'))
    assertEquals(t, "+52 55 1234 5", formatter.InputDigit('5'))
    assertEquals(t, "+52 55 1234 56", formatter.InputDigit('6'))
    assertEquals(t, "+52 55 1234 567", formatter.InputDigit('7'))
    assertEquals(t, "+52 55 1234 5678", formatter.InputDigit('8'))

    // +52 212 345 6789
    formatter.Clear()
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+5", formatter.InputDigit('5'))
    assertEquals(t, "+52 ", formatter.InputDigit('2'))
    assertEquals(t, "+52 2", formatter.InputDigit('2'))
    assertEquals(t, "+52 21", formatter.InputDigit('1'))
    assertEquals(t, "+52 212", formatter.InputDigit('2'))
    assertEquals(t, "+52 212 3", formatter.InputDigit('3'))
    assertEquals(t, "+52 212 34", formatter.InputDigit('4'))
    assertEquals(t, "+52 212 345", formatter.InputDigit('5'))
    assertEquals(t, "+52 212 345 6", formatter.InputDigit('6'))
    assertEquals(t, "+52 212 345 67", formatter.InputDigit('7'))
    assertEquals(t, "+52 212 345 678", formatter.InputDigit('8'))
    assertEquals(t, "+52 212 345 6789", formatter.InputDigit('9'))

    // +52 1 55 1234 5678
    formatter.Clear()
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+5", formatter.InputDigit('5'))
    assertEquals(t, "+52 ", formatter.InputDigit('2'))
    assertEquals(t, "+52 1", formatter.InputDigit('1'))
    assertEquals(t, "+52 15", formatter.InputDigit('5'))
    assertEquals(t, "+52 1 55", formatter.InputDigit('5'))
    assertEquals(t, "+52 1 55 1", formatter.InputDigit('1'))
    assertEquals(t, "+52 1 55 12", formatter.InputDigit('2'))
    assertEquals(t, "+52 1 55 123", formatter.InputDigit('3'))
    assertEquals(t, "+52 1 55 1234", formatter.InputDigit('4'))
    assertEquals(t, "+52 1 55 1234 5", formatter.InputDigit('5'))
    assertEquals(t, "+52 1 55 1234 56", formatter.InputDigit('6'))
    assertEquals(t, "+52 1 55 1234 567", formatter.InputDigit('7'))
    assertEquals(t, "+52 1 55 1234 5678", formatter.InputDigit('8'))

    // +52 1 541 234 5678
    formatter.Clear()
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+5", formatter.InputDigit('5'))
    assertEquals(t, "+52 ", formatter.InputDigit('2'))
    assertEquals(t, "+52 1", formatter.InputDigit('1'))
    assertEquals(t, "+52 15", formatter.InputDigit('5'))
    assertEquals(t, "+52 1 54", formatter.InputDigit('4'))
    assertEquals(t, "+52 1 541", formatter.InputDigit('1'))
    assertEquals(t, "+52 1 541 2", formatter.InputDigit('2'))
    assertEquals(t, "+52 1 541 23", formatter.InputDigit('3'))
    assertEquals(t, "+52 1 541 234", formatter.InputDigit('4'))
    assertEquals(t, "+52 1 541 234 5", formatter.InputDigit('5'))
    assertEquals(t, "+52 1 541 234 56", formatter.InputDigit('6'))
    assertEquals(t, "+52 1 541 234 567", formatter.InputDigit('7'))
    assertEquals(t, "+52 1 541 234 5678", formatter.InputDigit('8'))
}

func TestAYTF_International_Toll_Free(t *testing.T) {
    formatter := NewAsYouTypeFormatter("US")
    // +800 1234 5678
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+8", formatter.InputDigit('8'))
    assertEquals(t, "+80", formatter.InputDigit('0'))
    assertEquals(t, "+800 ", formatter.InputDigit('0'))
    assertEquals(t, "+800 1", formatter.InputDigit('1'))
    assertEquals(t, "+800 12", formatter.InputDigit('2'))
    assertEquals(t, "+800 123", formatter.InputDigit('3'))
    assertEquals(t, "+800 1234", formatter.InputDigit('4'))
    assertEquals(t, "+800 1234 5", formatter.InputDigit('5'))
    assertEquals(t, "+800 1234 56", formatter.InputDigit('6'))
    assertEquals(t, "+800 1234 567", formatter.InputDigit('7'))
    assertEquals(t, "+800 1234 5678", formatter.InputDigit('8'))
    assertEquals(t, "+800123456789", formatter.InputDigit('9'))
}

func TestAYTFMultipleLeadingDigitPatterns(t *testing.T) {
    // +81 50 2345 6789
    formatter := NewAsYouTypeFormatter("JP")
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+8", formatter.InputDigit('8'))
    assertEquals(t, "+81 ", formatter.InputDigit('1'))
    assertEquals(t, "+81 5", formatter.InputDigit('5'))
    assertEquals(t, "+81 50", formatter.InputDigit('0'))
    assertEquals(t, "+81 50 2", formatter.InputDigit('2'))
    assertEquals(t, "+81 50 23", formatter.InputDigit('3'))
    assertEquals(t, "+81 50 234", formatter.InputDigit('4'))
    assertEquals(t, "+81 50 2345", formatter.InputDigit('5'))
    assertEquals(t, "+81 50 2345 6", formatter.InputDigit('6'))
    assertEquals(t, "+81 50 2345 67", formatter.InputDigit('7'))
    assertEquals(t, "+81 50 2345 678", formatter.InputDigit('8'))
    assertEquals(t, "+81 50 2345 6789", formatter.InputDigit('9'))

    // +81 222 12 5678
    formatter.Clear()
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+8", formatter.InputDigit('8'))
    assertEquals(t, "+81 ", formatter.InputDigit('1'))
    assertEquals(t, "+81 2", formatter.InputDigit('2'))
    assertEquals(t, "+81 22", formatter.InputDigit('2'))
    assertEquals(t, "+81 22 2", formatter.InputDigit('2'))
    assertEquals(t, "+81 22 21", formatter.InputDigit('1'))
    assertEquals(t, "+81 2221 2", formatter.InputDigit('2'))
    assertEquals(t, "+81 222 12 5", formatter.InputDigit('5'))
    assertEquals(t, "+81 222 12 56", formatter.InputDigit('6'))
    assertEquals(t, "+81 222 12 567", formatter.InputDigit('7'))
    assertEquals(t, "+81 222 12 5678", formatter.InputDigit('8'))

    // 011113
    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "01", formatter.InputDigit('1'))
    assertEquals(t, "011", formatter.InputDigit('1'))
    assertEquals(t, "011 1", formatter.InputDigit('1'))
    assertEquals(t, "011 11", formatter.InputDigit('1'))
    assertEquals(t, "011113", formatter.InputDigit('3'))

    // +81 3332 2 5678
    formatter.Clear()
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+8", formatter.InputDigit('8'))
    assertEquals(t, "+81 ", formatter.InputDigit('1'))
    assertEquals(t, "+81 3", formatter.InputDigit('3'))
    assertEquals(t, "+81 33", formatter.InputDigit('3'))
    assertEquals(t, "+81 33 3", formatter.InputDigit('3'))
    assertEquals(t, "+81 3332", formatter.InputDigit('2'))
    assertEquals(t, "+81 3332 2", formatter.InputDigit('2'))
    assertEquals(t, "+81 3332 2 5", formatter.InputDigit('5'))
    assertEquals(t, "+81 3332 2 56", formatter.InputDigit('6'))
    assertEquals(t, "+81 3332 2 567", formatter.InputDigit('7'))
    assertEquals(t, "+81 3332 2 5678", formatter.InputDigit('8'))
}

func TestAYTFLongIDD_AU(t *testing.T) {
    formatter := NewAsYouTypeFormatter("AU")
    // 0011 1 650 253 2250
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "00", formatter.InputDigit('0'))
    assertEquals(t, "001", formatter.InputDigit('1'))
    assertEquals(t, "0011", formatter.InputDigit('1'))
    assertEquals(t, "0011 1 ", formatter.InputDigit('1'))
    assertEquals(t, "0011 1 6", formatter.InputDigit('6'))
    assertEquals(t, "0011 1 65", formatter.InputDigit('5'))
    assertEquals(t, "0011 1 650", formatter.InputDigit('0'))
    assertEquals(t, "0011 1 650 2", formatter.InputDigit('2'))
    assertEquals(t, "0011 1 650 25", formatter.InputDigit('5'))
    assertEquals(t, "0011 1 650 253", formatter.InputDigit('3'))
    assertEquals(t, "0011 1 650 253 2", formatter.InputDigit('2'))
    assertEquals(t, "0011 1 650 253 22", formatter.InputDigit('2'))
    assertEquals(t, "0011 1 650 253 222", formatter.InputDigit('2'))
    assertEquals(t, "0011 1 650 253 2222", formatter.InputDigit('2'))

    // 0011 81 3332 2 5678
    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "00", formatter.InputDigit('0'))
    assertEquals(t, "001", formatter.InputDigit('1'))
    assertEquals(t, "0011", formatter.InputDigit('1'))
    assertEquals(t, "00118", formatter.InputDigit('8'))
    assertEquals(t, "0011 81 ", formatter.InputDigit('1'))
    assertEquals(t, "0011 81 3", formatter.InputDigit('3'))
    assertEquals(t, "0011 81 33", formatter.InputDigit('3'))
    assertEquals(t, "0011 81 33 3", formatter.InputDigit('3'))
    assertEquals(t, "0011 81 3332", formatter.InputDigit('2'))
    assertEquals(t, "0011 81 3332 2", formatter.InputDigit('2'))
    assertEquals(t, "0011 81 3332 2 5", formatter.InputDigit('5'))
    assertEquals(t, "0011 81 3332 2 56", formatter.InputDigit('6'))
    assertEquals(t, "0011 81 3332 2 567", formatter.InputDigit('7'))
    assertEquals(t, "0011 81 3332 2 5678", formatter.InputDigit('8'))

    // 0011 244 250 253 222
    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "00", formatter.InputDigit('0'))
    assertEquals(t, "001", formatter.InputDigit('1'))
    assertEquals(t, "0011", formatter.InputDigit('1'))
    assertEquals(t, "00112", formatter.InputDigit('2'))
    assertEquals(t, "001124", formatter.InputDigit('4'))
    assertEquals(t, "0011 244 ", formatter.InputDigit('4'))
    assertEquals(t, "0011 244 2", formatter.InputDigit('2'))
    assertEquals(t, "0011 244 25", formatter.InputDigit('5'))
    assertEquals(t, "0011 244 250", formatter.InputDigit('0'))
    assertEquals(t, "0011 244 250 2", formatter.InputDigit('2'))
    assertEquals(t, "0011 244 250 25", formatter.InputDigit('5'))
    assertEquals(t, "0011 244 250 253", formatter.InputDigit('3'))
    assertEquals(t, "0011 244 250 253 2", formatter.InputDigit('2'))
    assertEquals(t, "0011 244 250 253 22", formatter.InputDigit('2'))
    assertEquals(t, "0011 244 250 253 222", formatter.InputDigit('2'))
}

func TestAYTFLongIDD_KR(t *testing.T) {
    formatter := NewAsYouTypeFormatter("KR")
    // 00300 1 650 253 2222
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "00", formatter.InputDigit('0'))
    assertEquals(t, "003", formatter.InputDigit('3'))
    assertEquals(t, "0030", formatter.InputDigit('0'))
    assertEquals(t, "00300", formatter.InputDigit('0'))
    assertEquals(t, "00300 1 ", formatter.InputDigit('1'))
    assertEquals(t, "00300 1 6", formatter.InputDigit('6'))
    assertEquals(t, "00300 1 65", formatter.InputDigit('5'))
    assertEquals(t, "00300 1 650", formatter.InputDigit('0'))
    assertEquals(t, "00300 1 650 2", formatter.InputDigit('2'))
    assertEquals(t, "00300 1 650 25", formatter.InputDigit('5'))
    assertEquals(t, "00300 1 650 253", formatter.InputDigit('3'))
    assertEquals(t, "00300 1 650 253 2", formatter.InputDigit('2'))
    assertEquals(t, "00300 1 650 253 22", formatter.InputDigit('2'))
    assertEquals(t, "00300 1 650 253 222", formatter.InputDigit('2'))
    assertEquals(t, "00300 1 650 253 2222", formatter.InputDigit('2'))
}

func TestAYTFLongNDD_KR(t *testing.T) {
    formatter := NewAsYouTypeFormatter("KR")
    // 08811-9876-7890
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "08", formatter.InputDigit('8'))
    assertEquals(t, "088", formatter.InputDigit('8'))
    assertEquals(t, "0881", formatter.InputDigit('1'))
    assertEquals(t, "08811", formatter.InputDigit('1'))
    assertEquals(t, "08811-9", formatter.InputDigit('9'))
    assertEquals(t, "08811-98", formatter.InputDigit('8'))
    assertEquals(t, "08811-987", formatter.InputDigit('7'))
    assertEquals(t, "08811-9876", formatter.InputDigit('6'))
    assertEquals(t, "08811-9876-7", formatter.InputDigit('7'))
    assertEquals(t, "08811-9876-78", formatter.InputDigit('8'))
    assertEquals(t, "08811-9876-789", formatter.InputDigit('9'))
    assertEquals(t, "08811-9876-7890", formatter.InputDigit('0'))

    // 08500 11-9876-7890
    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "08", formatter.InputDigit('8'))
    assertEquals(t, "085", formatter.InputDigit('5'))
    assertEquals(t, "0850", formatter.InputDigit('0'))
    assertEquals(t, "08500 ", formatter.InputDigit('0'))
    assertEquals(t, "08500 1", formatter.InputDigit('1'))
    assertEquals(t, "08500 11", formatter.InputDigit('1'))
    assertEquals(t, "08500 11-9", formatter.InputDigit('9'))
    assertEquals(t, "08500 11-98", formatter.InputDigit('8'))
    assertEquals(t, "08500 11-987", formatter.InputDigit('7'))
    assertEquals(t, "08500 11-9876", formatter.InputDigit('6'))
    assertEquals(t, "08500 11-9876-7", formatter.InputDigit('7'))
    assertEquals(t, "08500 11-9876-78", formatter.InputDigit('8'))
    assertEquals(t, "08500 11-9876-789", formatter.InputDigit('9'))
    assertEquals(t, "08500 11-9876-7890", formatter.InputDigit('0'))
}

func TestAYTFLongNDD_SG(t *testing.T) {
    formatter := NewAsYouTypeFormatter("SG")
    // 777777 9876 7890
    assertEquals(t, "7", formatter.InputDigit('7'))
    assertEquals(t, "77", formatter.InputDigit('7'))
    assertEquals(t, "777", formatter.InputDigit('7'))
    assertEquals(t, "7777", formatter.InputDigit('7'))
    assertEquals(t, "77777", formatter.InputDigit('7'))
    assertEquals(t, "777777 ", formatter.InputDigit('7'))
    assertEquals(t, "777777 9", formatter.InputDigit('9'))
    assertEquals(t, "777777 98", formatter.InputDigit('8'))
    assertEquals(t, "777777 987", formatter.InputDigit('7'))
    assertEquals(t, "777777 9876", formatter.InputDigit('6'))
    assertEquals(t, "777777 9876 7", formatter.InputDigit('7'))
    assertEquals(t, "777777 9876 78", formatter.InputDigit('8'))
    assertEquals(t, "777777 9876 789", formatter.InputDigit('9'))
    assertEquals(t, "777777 9876 7890", formatter.InputDigit('0'))
}

func TestAYTFShortNumberFormattingFix_AU(t *testing.T) {
    // For Australia, the national prefix is not optional when formatting.
    formatter := NewAsYouTypeFormatter("AU")

    // 1234567890 - For leading digit 1, the national prefix formatting rule has first group only.
    assertEquals(t, "1", formatter.InputDigit('1'))
    assertEquals(t, "12", formatter.InputDigit('2'))
    assertEquals(t, "123", formatter.InputDigit('3'))
    assertEquals(t, "1234", formatter.InputDigit('4'))
    assertEquals(t, "1234 5", formatter.InputDigit('5'))
    assertEquals(t, "1234 56", formatter.InputDigit('6'))
    assertEquals(t, "1234 567", formatter.InputDigit('7'))
    assertEquals(t, "1234 567 8", formatter.InputDigit('8'))
    assertEquals(t, "1234 567 89", formatter.InputDigit('9'))
    assertEquals(t, "1234 567 890", formatter.InputDigit('0'))

    // +61 1234 567 890 - Test the same number, but with the country code.
    formatter.Clear()
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+6", formatter.InputDigit('6'))
    assertEquals(t, "+61 ", formatter.InputDigit('1'))
    assertEquals(t, "+61 1", formatter.InputDigit('1'))
    assertEquals(t, "+61 12", formatter.InputDigit('2'))
    assertEquals(t, "+61 123", formatter.InputDigit('3'))
    assertEquals(t, "+61 1234", formatter.InputDigit('4'))
    assertEquals(t, "+61 1234 5", formatter.InputDigit('5'))
    assertEquals(t, "+61 1234 56", formatter.InputDigit('6'))
    assertEquals(t, "+61 1234 567", formatter.InputDigit('7'))
    assertEquals(t, "+61 1234 567 8", formatter.InputDigit('8'))
    assertEquals(t, "+61 1234 567 89", formatter.InputDigit('9'))
    assertEquals(t, "+61 1234 567 890", formatter.InputDigit('0'))

    // 212345678 - For leading digit 2, the national prefix formatting rule puts the national prefix
    // before the first group.
    formatter.Clear()
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "02", formatter.InputDigit('2'))
    assertEquals(t, "021", formatter.InputDigit('1'))
    assertEquals(t, "02 12", formatter.InputDigit('2'))
    assertEquals(t, "02 123", formatter.InputDigit('3'))
    assertEquals(t, "02 1234", formatter.InputDigit('4'))
    assertEquals(t, "02 1234 5", formatter.InputDigit('5'))
    assertEquals(t, "02 1234 56", formatter.InputDigit('6'))
    assertEquals(t, "02 1234 567", formatter.InputDigit('7'))
    assertEquals(t, "02 1234 5678", formatter.InputDigit('8'))

    // 212345678 - Test the same number, but without the leading 0.
    formatter.Clear()
    assertEquals(t, "2", formatter.InputDigit('2'))
    assertEquals(t, "21", formatter.InputDigit('1'))
    assertEquals(t, "212", formatter.InputDigit('2'))
    assertEquals(t, "2123", formatter.InputDigit('3'))
    assertEquals(t, "21234", formatter.InputDigit('4'))
    assertEquals(t, "212345", formatter.InputDigit('5'))
    assertEquals(t, "2123456", formatter.InputDigit('6'))
    assertEquals(t, "21234567", formatter.InputDigit('7'))
    assertEquals(t, "212345678", formatter.InputDigit('8'))

    // +61 2 1234 5678 - Test the same number, but with the country code.
    formatter.Clear()
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+6", formatter.InputDigit('6'))
    assertEquals(t, "+61 ", formatter.InputDigit('1'))
    assertEquals(t, "+61 2", formatter.InputDigit('2'))
    assertEquals(t, "+61 21", formatter.InputDigit('1'))
    assertEquals(t, "+61 2 12", formatter.InputDigit('2'))
    assertEquals(t, "+61 2 123", formatter.InputDigit('3'))
    assertEquals(t, "+61 2 1234", formatter.InputDigit('4'))
    assertEquals(t, "+61 2 1234 5", formatter.InputDigit('5'))
    assertEquals(t, "+61 2 1234 56", formatter.InputDigit('6'))
    assertEquals(t, "+61 2 1234 567", formatter.InputDigit('7'))
    assertEquals(t, "+61 2 1234 5678", formatter.InputDigit('8'))
}

func TestAYTFShortNumberFormattingFix_KR(t *testing.T) {
    // For Korea, the national prefix is not optional when formatting, and the national prefix
    // formatting rule doesn't consist of only the first group.
    formatter := NewAsYouTypeFormatter("KR")

    // 111
    assertEquals(t, "1", formatter.InputDigit('1'))
    assertEquals(t, "11", formatter.InputDigit('1'))
    assertEquals(t, "111", formatter.InputDigit('1'))

    // 114
    formatter.Clear()
    assertEquals(t, "1", formatter.InputDigit('1'))
    assertEquals(t, "11", formatter.InputDigit('1'))
    assertEquals(t, "114", formatter.InputDigit('4'))

    // 13121234 - Test a mobile number without the national prefix. Even though it is not an
    // emergency number, it should be formatted as a block.
    formatter.Clear()
    assertEquals(t, "1", formatter.InputDigit('1'))
    assertEquals(t, "13", formatter.InputDigit('3'))
    assertEquals(t, "131", formatter.InputDigit('1'))
    assertEquals(t, "1312", formatter.InputDigit('2'))
    assertEquals(t, "13121", formatter.InputDigit('1'))
    assertEquals(t, "131212", formatter.InputDigit('2'))
    assertEquals(t, "1312123", formatter.InputDigit('3'))
    assertEquals(t, "13121234", formatter.InputDigit('4'))

    // +82 131-2-1234 - Test the same number, but with the country code.
    formatter.Clear()
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+8", formatter.InputDigit('8'))
    assertEquals(t, "+82 ", formatter.InputDigit('2'))
    assertEquals(t, "+82 1", formatter.InputDigit('1'))
    assertEquals(t, "+82 13", formatter.InputDigit('3'))
    assertEquals(t, "+82 131", formatter.InputDigit('1'))
    assertEquals(t, "+82 131-2", formatter.InputDigit('2'))
    assertEquals(t, "+82 131-2-1", formatter.InputDigit('1'))
    assertEquals(t, "+82 131-2-12", formatter.InputDigit('2'))
    assertEquals(t, "+82 131-2-123", formatter.InputDigit('3'))
    assertEquals(t, "+82 131-2-1234", formatter.InputDigit('4'))
}

func TestAYTFShortNumberFormattingFix_MX(t *testing.T) {
    // For Mexico, the national prefix is optional when formatting.
    formatter := NewAsYouTypeFormatter("MX")

    // 911
    assertEquals(t, "9", formatter.InputDigit('9'))
    assertEquals(t, "91", formatter.InputDigit('1'))
    assertEquals(t, "911", formatter.InputDigit('1'))

    // 800 123 4567 - Test a toll-free number, which should have a formatting rule applied to it
    // even though it doesn't begin with the national prefix.
    formatter.Clear()
    assertEquals(t, "8", formatter.InputDigit('8'))
    assertEquals(t, "80", formatter.InputDigit('0'))
    assertEquals(t, "800", formatter.InputDigit('0'))
    assertEquals(t, "800 1", formatter.InputDigit('1'))
    assertEquals(t, "800 12", formatter.InputDigit('2'))
    assertEquals(t, "800 123", formatter.InputDigit('3'))
    assertEquals(t, "800 123 4", formatter.InputDigit('4'))
    assertEquals(t, "800 123 45", formatter.InputDigit('5'))
    assertEquals(t, "800 123 456", formatter.InputDigit('6'))
    assertEquals(t, "800 123 4567", formatter.InputDigit('7'))

    // +52 800 123 4567 - Test the same number, but with the country code.
    formatter.Clear()
    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+5", formatter.InputDigit('5'))
    assertEquals(t, "+52 ", formatter.InputDigit('2'))
    assertEquals(t, "+52 8", formatter.InputDigit('8'))
    assertEquals(t, "+52 80", formatter.InputDigit('0'))
    assertEquals(t, "+52 800", formatter.InputDigit('0'))
    assertEquals(t, "+52 800 1", formatter.InputDigit('1'))
    assertEquals(t, "+52 800 12", formatter.InputDigit('2'))
    assertEquals(t, "+52 800 123", formatter.InputDigit('3'))
    assertEquals(t, "+52 800 123 4", formatter.InputDigit('4'))
    assertEquals(t, "+52 800 123 45", formatter.InputDigit('5'))
    assertEquals(t, "+52 800 123 456", formatter.InputDigit('6'))
    assertEquals(t, "+52 800 123 4567", formatter.InputDigit('7'))
}

func TestAYTFNoNationalPrefix(t *testing.T) {
    formatter := NewAsYouTypeFormatter("IT")

    assertEquals(t, "3", formatter.InputDigit('3'))
    assertEquals(t, "33", formatter.InputDigit('3'))
    assertEquals(t, "333", formatter.InputDigit('3'))
    assertEquals(t, "333 3", formatter.InputDigit('3'))
    assertEquals(t, "333 33", formatter.InputDigit('3'))
    assertEquals(t, "333 333", formatter.InputDigit('3'))
}

func TestAYTFNoNationalPrefixFormattingRule(t *testing.T) {
    formatter := NewAsYouTypeFormatter("AO")

    assertEquals(t, "3", formatter.InputDigit('3'))
    assertEquals(t, "33", formatter.InputDigit('3'))
    assertEquals(t, "333", formatter.InputDigit('3'))
    assertEquals(t, "333 3", formatter.InputDigit('3'))
    assertEquals(t, "333 33", formatter.InputDigit('3'))
    assertEquals(t, "333 333", formatter.InputDigit('3'))
}

func TestAYTFShortNumberFormattingFix_US(t *testing.T) {
    // For the US, an initial 1 is treated specially.
    formatter := NewAsYouTypeFormatter("US")

    // 101 - Test that the initial 1 is not treated as a national prefix.
    assertEquals(t, "1", formatter.InputDigit('1'))
    assertEquals(t, "10", formatter.InputDigit('0'))
    assertEquals(t, "101", formatter.InputDigit('1'))

    // 112 - Test that the initial 1 is not treated as a national prefix.
    formatter.Clear()
    assertEquals(t, "1", formatter.InputDigit('1'))
    assertEquals(t, "11", formatter.InputDigit('1'))
    assertEquals(t, "112", formatter.InputDigit('2'))

    // 122 - Test that the initial 1 is treated as a national prefix.
    formatter.Clear()
    assertEquals(t, "1", formatter.InputDigit('1'))
    assertEquals(t, "12", formatter.InputDigit('2'))
    assertEquals(t, "1 22", formatter.InputDigit('2'))
}

func TestAYTFClearNDDAfterIDDExtraction(t *testing.T) {
    formatter := NewAsYouTypeFormatter("KR")

    // Check that when we have successfully extracted an IDD, the previously extracted NDD is
    // cleared since it is no longer valid.
    assertEquals(t, "0", formatter.InputDigit('0'))
    assertEquals(t, "00", formatter.InputDigit('0'))
    assertEquals(t, "007", formatter.InputDigit('7'))
    assertEquals(t, "0070", formatter.InputDigit('0'))
    assertEquals(t, "00700", formatter.InputDigit('0'))
    assertEquals(t, "0", formatter.GetExtractedNationalPrefix())

    // Once the IDD "00700" has been extracted, it no longer makes sense for the initial "0" to be
    // treated as an NDD.
    assertEquals(t, "00700 1 ", formatter.InputDigit('1'))
    assertEquals(t, "", formatter.GetExtractedNationalPrefix())

    assertEquals(t, "00700 1 2", formatter.InputDigit('2'))
    assertEquals(t, "00700 1 23", formatter.InputDigit('3'))
    assertEquals(t, "00700 1 234", formatter.InputDigit('4'))
    assertEquals(t, "00700 1 234 5", formatter.InputDigit('5'))
    assertEquals(t, "00700 1 234 56", formatter.InputDigit('6'))
    assertEquals(t, "00700 1 234 567", formatter.InputDigit('7'))
    assertEquals(t, "00700 1 234 567 8", formatter.InputDigit('8'))
    assertEquals(t, "00700 1 234 567 89", formatter.InputDigit('9'))
    assertEquals(t, "00700 1 234 567 890", formatter.InputDigit('0'))
    assertEquals(t, "00700 1 234 567 8901", formatter.InputDigit('1'))
    assertEquals(t, "00700123456789012", formatter.InputDigit('2'))
    assertEquals(t, "007001234567890123", formatter.InputDigit('3'))
    assertEquals(t, "0070012345678901234", formatter.InputDigit('4'))
    assertEquals(t, "00700123456789012345", formatter.InputDigit('5'))
    assertEquals(t, "007001234567890123456", formatter.InputDigit('6'))
    assertEquals(t, "0070012345678901234567", formatter.InputDigit('7'))
}

func TestAYTFNumberPatternsBecomingInvalidShouldNotResultInDigitLoss(t *testing.T) {
    formatter := NewAsYouTypeFormatter("CN")

    assertEquals(t, "+", formatter.InputDigit('+'))
    assertEquals(t, "+8", formatter.InputDigit('8'))
    assertEquals(t, "+86 ", formatter.InputDigit('6'))
    assertEquals(t, "+86 9", formatter.InputDigit('9'))
    assertEquals(t, "+86 98", formatter.InputDigit('8'))
    assertEquals(t, "+86 988", formatter.InputDigit('8'))
    assertEquals(t, "+86 988 1", formatter.InputDigit('1'))
    // Now the number pattern is no longer valid because there are multiple leading digit patterns
    // when we try again to extract a country code we should ensure we use the last leading digit
    // pattern, rather than the first one such that it *thinks* it's found a valid formatting rule
    // again.
    // https://github.com/google/libphonenumber/issues/437
    assertEquals(t, "+8698812", formatter.InputDigit('2'))
    assertEquals(t, "+86988123", formatter.InputDigit('3'))
    assertEquals(t, "+869881234", formatter.InputDigit('4'))
    assertEquals(t, "+8698812345", formatter.InputDigit('5'))
}
