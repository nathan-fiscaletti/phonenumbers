package phonenumbers

import(
    "fmt"
    "regexp"
    "math"
    "strings"
    "unicode"
)

const(
    // This is the minimum length of national number accrued that is required to trigger the
    // formatter. The first element of the leadingDigitsPattern of each numberFormat contains a
    // regular expression that matches up to this number of digits.
    MIN_LEADING_DIGITS_LENGTH = 3

    // The digits that have not been entered yet will be represented by a \u2008, the punctuation
    // space.
    DIGIT_PLACEHOLDER = "\u2008"

    // Character used when appropriate to separate a prefix, such as a long NDD or a country calling
    // code, from the national number.
    SEPARATOR_BEFORE_NATIONAL_NUMBER = ' '
)

// A pattern that is used to determine if a numberFormat under availableFormats is eligible to be
// used by the AYTF. It is eligible when the format element under numberFormat contains groups of
// the dollar sign followed by a single digit, separated by valid phone number punctuation. This
// prevents invalid punctuation (such as the star sign in Israeli star numbers) getting into the
// output of the AYTF.
var ELIGIBLE_FORMAT_PATTERN = regexp.MustCompile(
    "[" + VALID_PUNCTUATION + "]*" + "(\\$\\d" + "[" + VALID_PUNCTUATION + "]*)+",
)

// A set of characters that, if found in a national prefix formatting rules, are an indicator to
// us that we should separate the national prefix from the number when formatting.
var NATIONAL_PREFIX_SEPARATORS_PATTERN = regexp.MustCompile("[- ]")

var DIGIT_PATTERN = regexp.MustCompile(DIGIT_PLACEHOLDER)

type AsYouTypeFormatter struct {
    currentOutput           string
    defaultCountry          string
    extractedNationalPrefix string
    currentFormattingPattern string

    lastMatchPosition  int
    origingalPosition  int
    positionToRemember int

    isCompleteNumber                  bool
    ableToFormat                      bool
    isExpectingCountryCallingCode     bool
    inputHasFormatting                bool
    shouldAddSpaceAfterNationalPrefix bool

    formattingTemplate            *Builder
    nationalNumber                *Builder
    accruedInput                  *Builder
    accruedInputWithoutFormatting *Builder
    prefixBeforeNationalNumber    *Builder

    currentMetadata *PhoneMetadata
    defaultMetadata *PhoneMetadata

    numberFormat *NumberFormat
    
    possibleFormats []*NumberFormat
}

func emptyMetaData() *PhoneMetadata {
    var ignored string = "<ignored>"
    var ipfx string = "NA"
    return &PhoneMetadata{
        Id: &ignored,
        InternationalPrefix: &ipfx,
    }
}

func getMetadataForRegionCode(regionCode string) *PhoneMetadata {
    countryCallingCode := GetCountryCodeForRegion(regionCode)
    mainCountry := GetRegionCodeForCountryCode(countryCallingCode)
    metadata := getMetadataForRegion(mainCountry)
    if metadata == nil {
        return emptyMetaData()
    }
    return metadata
}

func NewAsYouTypeFormatter(regionCode string) *AsYouTypeFormatter {
    metadata := getMetadataForRegionCode(regionCode)

    return &AsYouTypeFormatter{
        defaultCountry: regionCode,
        currentMetadata: metadata,
        defaultMetadata: metadata,
        numberFormat: &NumberFormat{},
        possibleFormats: []*NumberFormat{},
        currentFormattingPattern: "",
        ableToFormat: true,
        formattingTemplate: NewBuilderString(""),
        nationalNumber: NewBuilderString(""),
        accruedInput: NewBuilderString(""),
        accruedInputWithoutFormatting: NewBuilderString(""),
        prefixBeforeNationalNumber: NewBuilderString(""),
    }
}

func (f *AsYouTypeFormatter) InputDigit(nextChar rune) string {
    f.currentOutput = f.inputDigitWithOptionToRememberPosition(nextChar, false)
    return f.currentOutput
}

func (f *AsYouTypeFormatter) InputDigitAndRememberPosition(nextChar rune) string {
    f.currentOutput = f.inputDigitWithOptionToRememberPosition(nextChar, true)
    return f.currentOutput
}

func (f *AsYouTypeFormatter) GetExtractedNationalPrefix() string {
    return f.extractedNationalPrefix
}

func (f *AsYouTypeFormatter) GetRememberedPosition() int {
    if !f.ableToFormat {
        return f.origingalPosition
    }

    accruedInputIndex := 0
    currentOutputIndex := 0
    for accruedInputIndex < f.positionToRemember && currentOutputIndex < len(f.currentOutput) {
       if f.accruedInputWithoutFormatting.String()[accruedInputIndex] == f.currentOutput[currentOutputIndex] {
           accruedInputIndex++
       }
       currentOutputIndex++
    }
    return currentOutputIndex
}

func (f *AsYouTypeFormatter) inputDigitWithOptionToRememberPosition(nextChar rune, rememberPosition bool) string {
    f.accruedInput.WriteRune(nextChar)
    if rememberPosition {
        f.origingalPosition = f.accruedInput.Len()
    }
    
    // We do formatting on-the-fly only when each character entered is either a digit, or a plus
    // sign (accepted at the start of the number only).
    if !f.isDigitOrLeadingPlusSign(nextChar) {
        f.ableToFormat = false
        f.inputHasFormatting = true
    } else {
        nextChar = f.normalizeAndAccrueDigitsAndPlusSign(nextChar, rememberPosition)
    }

    if !f.ableToFormat {
        // When we are unable to format because of reasons other than that formatting chars have been
        // entered, it can be due to really long IDDs or NDDs. If that is the case, we might be able
        // to do formatting again after extracting them.
        if f.inputHasFormatting {
            return f.accruedInput.String()
        } else if f.attemptToExtractIdd() {
            if f.attemptToExtractCountryCallingCode() {
                return f.attemptToChoosePatternWithPrefixExtracted()
            }
        } else if f.ableToExtractLongerNdd() {
            // Add an additional space to separate long NDD and national significant number for
            // readability. We don't set shouldAddSpaceAfterNationalPrefix to true, since we don't want
            // this to change later when we choose formatting templates.
            f.prefixBeforeNationalNumber.WriteRune(SEPARATOR_BEFORE_NATIONAL_NUMBER)
            return f.attemptToChoosePatternWithPrefixExtracted()
        }
        return f.accruedInput.String()
    }

    switch f.accruedInputWithoutFormatting.Len() {
    case 0:
        fallthrough
    case 1:
        fallthrough
    case 2:
        return f.accruedInput.String()
    case 3:
        if f.attemptToExtractIdd() {
            f.isExpectingCountryCallingCode = true
        } else { // No IDD or plus sign is found, might be entering in national format.
            f.extractedNationalPrefix = f.removeNationalPrefixFromNationalNumber()
            return f.attemptToChooseFormattingPattern()
        }
        fallthrough
    default:
        if f.isExpectingCountryCallingCode {
            if f.attemptToExtractCountryCallingCode() {
                f.isExpectingCountryCallingCode = false
            }
            return fmt.Sprintf("%s%s", f.prefixBeforeNationalNumber, f.nationalNumber.String())
        }
        if len(f.possibleFormats) > 0 { // The formatting patterns are already chosen.
            tempNationalNumber := f.inputDigitHelper(nextChar)
            // See if the accrued digits can be formatted properly already. If not, use the results
            // from inputDigitHelper, which does formatting based on the formatting pattern chosen.
            formattedNumber := f.attemptToFormatAccruedDigits()
            if len(formattedNumber) > 0 {
                return formattedNumber
            }
            f.narrowDownPossibleFormats(f.nationalNumber.String())
            if f.maybeCreateNewTemplate() {
                return f.inputAccruedNationalNumber()
            }
            if f.ableToFormat {
                return f.appendNationalNumber(tempNationalNumber)
            } else {
                return f.accruedInput.String()
            }
        } else {
            return f.attemptToChooseFormattingPattern()
        }
    }
  }

func (f *AsYouTypeFormatter) getAvailableFormats(leadingDigits string) {
    isInternationalNumber := f.isCompleteNumber && len(f.extractedNationalPrefix) == 0
    var formatList []*NumberFormat

    intlNumberFormatList := f.currentMetadata.GetIntlNumberFormat()
    if isInternationalNumber && len(intlNumberFormatList) > 0 {
        formatList = f.currentMetadata.GetIntlNumberFormat()
    } else {
        formatList = f.currentMetadata.GetNumberFormat()
    }

    for _,format := range formatList {
        if len(f.extractedNationalPrefix) > 0 &&
           format.NationalPrefixFormattingRule != nil && formattingRuleHasFirstGroupOnly(*format.NationalPrefixFormattingRule) &&
           format.NationalPrefixOptionalWhenFormatting != nil && !*format.NationalPrefixOptionalWhenFormatting &&
           format.DomesticCarrierCodeFormattingRule == nil {
            // If it is a national number that had a national prefix, any rules that aren't valid with a
            // national prefix should be excluded. A rule that has a carrier-code formatting rule is
            // kept since the national prefix might actually be an extracted carrier code - we don't
            // distinguish between these when extracting it in the AYTF.
            continue
        } else if len(f.extractedNationalPrefix) == 0 &&
                     !f.isCompleteNumber &&
                     format.NationalPrefixFormattingRule != nil && !formattingRuleHasFirstGroupOnly(*format.NationalPrefixFormattingRule) &&
                     !*format.NationalPrefixOptionalWhenFormatting {
            // This number was entered without a national prefix, and this formatting rule requires one,
            // so we discard it.
            continue
        }

        if ELIGIBLE_FORMAT_PATTERN.MatchString(format.GetFormat()) {
            f.possibleFormats = append(f.possibleFormats, format)
        }
    }

    f.narrowDownPossibleFormats(leadingDigits)
}

func (f *AsYouTypeFormatter) narrowDownPossibleFormats(leadingDigits string) {
    indexOfLeadingDigitsPattern := len(leadingDigits) - MIN_LEADING_DIGITS_LENGTH

    output_index := 0
    for _,format := range f.possibleFormats {
        if len(format.GetLeadingDigitsPattern()) == 0 {
            // Keep everything that isn't restricted by leading digits.
            continue
        }

        lastLeadingDigitsPattern := int(math.Min(float64(indexOfLeadingDigitsPattern), float64(len(format.GetLeadingDigitsPattern()) - 1)))
        r := regexFor(format.GetLeadingDigitsPattern()[lastLeadingDigitsPattern])
        matched, _, _ := lookingAt(r, leadingDigits)
        if !matched {
            f.possibleFormats[output_index] = format
            output_index++
        }
    }
    if output_index != 0 {
        f.possibleFormats = f.possibleFormats[:output_index]
    }
}

// Returns true if a new template is created as opposed to reusing the existing template.
func (f *AsYouTypeFormatter) maybeCreateNewTemplate() bool {
    defer func() {
        // delete all formats marked for deletion
        newPossibleFormats := []*NumberFormat{}
        for _,format := range f.possibleFormats {
            if format != nil {
                newPossibleFormats = append(newPossibleFormats, format)
            }
        }
        f.possibleFormats = newPossibleFormats
    }()

    // When there are multiple available formats, the formatter uses the first format where a
    // formatting template could be created.
    for idx,format := range f.possibleFormats {
        pattern := format.GetPattern()
        if f.currentFormattingPattern == pattern {
            return false
        }
        if f.createFormattingTemplate(format) {
            f.currentFormattingPattern = pattern
            f.shouldAddSpaceAfterNationalPrefix = NATIONAL_PREFIX_SEPARATORS_PATTERN.MatchString(
                format.GetNationalPrefixFormattingRule(),
            )
            // With a new formatting template, the matched position using the old template needs to be
            // reset.
            f.lastMatchPosition = 0
            return true
        } else {  // set this index to nil to mark it for deletion
            f.possibleFormats[idx] = nil
        }
    }

    f.ableToFormat = false
    return false
}

func (f *AsYouTypeFormatter) createFormattingTemplate(format *NumberFormat) bool {
    numberPattern := format.GetPattern()
    f.formattingTemplate.Reset()

    tempTemplate := f.getFormattingTemplate(numberPattern, format.GetFormat())
    if len(tempTemplate) > 0 {
        f.formattingTemplate.WriteString(tempTemplate)
        return true
    }
    
    return false
}

// Gets a formatting template which can be used to efficiently format a partial number where
// digits are added one by one.
func (f *AsYouTypeFormatter) getFormattingTemplate(numberPattern string, numberFormat string) string {
    // Creates a phone number consisting only of the digit 9 that matches the
    // numberPattern by applying the pattern to the longestPhoneNumber string.
    longestPhoneNumber := "999999999999999"
    
    pattern := regexFor(numberPattern)
    _, start, end := regexFind(pattern, longestPhoneNumber, 0)
    aPhoneNumber := longestPhoneNumber[start:end]

    // No formatting template can be created if the number of digits entered so far is longer than
    // the maximum the current formatting rule can accommodate.
    if len(aPhoneNumber) < f.nationalNumber.Len() {
        return ""
    }

    // Formats the number according to numberFormat
    template := strings.Replace(aPhoneNumber, numberPattern, numberFormat, -1)
    template = strings.Replace(template, "9", DIGIT_PLACEHOLDER, -1)
    return template
}

func (f *AsYouTypeFormatter) isDigitOrLeadingPlusSign(nextChar rune) bool {
    return unicode.IsDigit(nextChar) ||
           (f.accruedInput.Len() == 1 &&
               PLUS_CHARS_PATTERN.MatchString(string(nextChar)))
}

// Accrues digits and the plus sign to accruedInputWithoutFormatting for later use. If nextChar
// contains a digit in non-ASCII format (e.g. the full-width version of digits), it is first
// normalized to the ASCII version. The return value is nextChar itself, or its normalized
// version, if nextChar is a digit in non-ASCII format. This method assumes its input is either a
// digit or the plus sign.
func (f *AsYouTypeFormatter) normalizeAndAccrueDigitsAndPlusSign(nextChar rune, rememberPosition bool) rune {
    var normalizedChar rune
    if nextChar == PLUS_SIGN {
        normalizedChar = nextChar
        f.accruedInputWithoutFormatting.WriteRune(nextChar)
    } else {
        normalizedChar = nextChar
        // radix := 10
        // normalizedChar = Character.forDigit(Character.digit(nextChar, radix), radix);
        // accruedInputWithoutFormatting.append(normalizedChar);
        f.accruedInputWithoutFormatting.WriteRune(normalizedChar)
        f.nationalNumber.WriteRune(normalizedChar)
    }

    if rememberPosition {
        f.positionToRemember = f.accruedInputWithoutFormatting.Len()
    }

    return normalizedChar
}

/**
* Extracts IDD and plus sign to prefixBeforeNationalNumber when they are available, and places
* the remaining input into nationalNumber.
*
* @return  true when accruedInputWithoutFormatting begins with the plus sign or valid IDD for
*     defaultCountry.
*/
func (f *AsYouTypeFormatter) attemptToExtractIdd() bool {
    internationalPrefix := regexFor(fmt.Sprintf("\\%c|%s", PLUS_SIGN, f.currentMetadata.GetInternationalPrefix()))
    matched, _, end := lookingAt(internationalPrefix, f.accruedInputWithoutFormatting.String())
    if matched {
        f.isCompleteNumber = true
        startOfCountryCallingCode := end
        f.nationalNumber.Reset()
        f.nationalNumber.WriteString(f.accruedInputWithoutFormatting.String()[startOfCountryCallingCode:])
        f.prefixBeforeNationalNumber.Reset()
        f.prefixBeforeNationalNumber.WriteString(f.accruedInputWithoutFormatting.String()[:startOfCountryCallingCode])
        if f.accruedInputWithoutFormatting.String()[0] != PLUS_SIGN {
            f.prefixBeforeNationalNumber.WriteString(string(SEPARATOR_BEFORE_NATIONAL_NUMBER))
        }
        return true
    }
    return false
}

/**
* Extracts the country calling code from the beginning of nationalNumber to
* prefixBeforeNationalNumber when they are available, and places the remaining input into
* nationalNumber.
*
* @return  true when a valid country calling code can be found.
*/
func (f *AsYouTypeFormatter) attemptToExtractCountryCallingCode() bool {
    if f.nationalNumber.Len() == 0 {
        return false
    }

    var numberWithoutCountryCallingCode *Builder = NewBuilderString("")
    
    countryCode := extractCountryCode(f.nationalNumber, numberWithoutCountryCallingCode)
    if countryCode == 0 {
      return false
    }

    f.nationalNumber.Reset()
    f.nationalNumber.WriteString(numberWithoutCountryCallingCode.String())

    newRegionCode := GetRegionCodeForCountryCode(countryCode)
    if REGION_CODE_FOR_NON_GEO_ENTITY == newRegionCode {
        f.currentMetadata = getMetadataForNonGeographicalRegion(countryCode)
    } else if newRegionCode != f.defaultCountry {
        f.currentMetadata = getMetadataForRegionCode(newRegionCode)
    }

    countryCodeString := fmt.Sprintf("%v", countryCode)
    f.prefixBeforeNationalNumber.WriteString(countryCodeString)
    f.prefixBeforeNationalNumber.WriteRune(SEPARATOR_BEFORE_NATIONAL_NUMBER)

    // When we have successfully extracted the IDD, the previously extracted NDD should be cleared
    // because it is no longer valid.
    f.extractedNationalPrefix = ""
    
    return true
}

func (f *AsYouTypeFormatter) attemptToChoosePatternWithPrefixExtracted() string {
    f.ableToFormat = true
    f.isExpectingCountryCallingCode = false
    f.possibleFormats = []*NumberFormat{}
    f.lastMatchPosition = 0
    f.formattingTemplate.Reset()
    f.currentFormattingPattern = ""

    return f.attemptToChooseFormattingPattern()
}

/**
* Attempts to set the formatting template and returns a string which contains the formatted
* version of the digits entered so far.
*/
func (f *AsYouTypeFormatter) attemptToChooseFormattingPattern() string {
    // We start to attempt to format only when at least MIN_LEADING_DIGITS_LENGTH digits of national
    // number (excluding national prefix) have been entered.
    if f.nationalNumber.Len() >= MIN_LEADING_DIGITS_LENGTH {
        f.getAvailableFormats(f.nationalNumber.String())
        // See if the accrued digits can be formatted properly already.
        formattedNumber := f.attemptToFormatAccruedDigits()
        if len(formattedNumber) > 0 {
            return formattedNumber
        }

        if f.maybeCreateNewTemplate() {
            return f.inputAccruedNationalNumber()
        }

        return f.accruedInput.String()
    }

    return f.appendNationalNumber(f.nationalNumber.String())
}

/**
* Checks to see if there is an exact pattern match for these digits. If so, we should use this
* instead of any other formatting template whose leadingDigitsPattern also matches the input.
*/
func (f *AsYouTypeFormatter) attemptToFormatAccruedDigits() string {
    for _,numberFormat := range f.possibleFormats {
        m := regexFor(numberFormat.GetPattern())
        if m.MatchString(f.nationalNumber.String()) {
            f.shouldAddSpaceAfterNationalPrefix, _, _ = regexFind(NATIONAL_PREFIX_SEPARATORS_PATTERN, numberFormat.GetNationalPrefixFormattingRule(), 0)
            formattedNumber := m.ReplaceAllString(f.nationalNumber.String(), numberFormat.GetFormat())
            // Check that we did not remove nor add any extra digits when we matched
            // this formatting pattern. This usually happens after we entered the last
            // digit during AYTF. Eg: In case of MX, we swallow mobile token (1) when
            // formatted but AYTF should retain all the number entered and not change
            // in order to match a format (of same leading digits and length) display
            // in that way.
            fullOutput := f.appendNationalNumber(formattedNumber)
            formattedNumberDigitsOnly := normalizeDiallableCharsOnly(fullOutput)
            if strings.Contains(formattedNumberDigitsOnly, f.accruedInputWithoutFormatting.String()) {
                return fullOutput
            }
        }
    }
    return ""
}

/**
* Combines the national number with any prefix (IDD/+ and country code or national prefix) that
* was collected. A space will be inserted between them if the current formatting template
* indicates this to be suitable.
*/
func (f *AsYouTypeFormatter) appendNationalNumber(nationalNumber string ) string {
    prefixBeforeNationalNumberLength := f.prefixBeforeNationalNumber.Len();

    if f.shouldAddSpaceAfterNationalPrefix && prefixBeforeNationalNumberLength > 0 && 
       f.prefixBeforeNationalNumber.String()[prefixBeforeNationalNumberLength - 1] != SEPARATOR_BEFORE_NATIONAL_NUMBER {
        // We want to add a space after the national prefix if the national prefix formatting rule
        // indicates that this would normally be done, with the exception of the case where we already
        // appended a space because the NDD was surprisingly long.
        return fmt.Sprintf("%s%c%s", f.prefixBeforeNationalNumber.String(), SEPARATOR_BEFORE_NATIONAL_NUMBER, nationalNumber)
    } else {
        return fmt.Sprintf("%s%s", f.prefixBeforeNationalNumber.String(), nationalNumber)
    }
}

/**
* Invokes inputDigitHelper on each digit of the national number accrued, and returns a formatted
* string in the end.
*/
func (f *AsYouTypeFormatter) inputAccruedNationalNumber() string {
    lengthOfNationalNumber := f.nationalNumber.Len()
    
    if lengthOfNationalNumber > 0 {
        tempNationalNumber := ""
        for i := 0; i < lengthOfNationalNumber; i++ {
            tempNationalNumber = f.inputDigitHelper(rune(f.nationalNumber.String()[i]))
        }
        if f.ableToFormat {
            return f.appendNationalNumber(tempNationalNumber)
        } else {
            return f.accruedInput.String()
        }
    }

    return f.prefixBeforeNationalNumber.String()
}

func (f *AsYouTypeFormatter) Clear() {
    f.currentOutput = ""
    f.accruedInput.Reset()
    f.accruedInputWithoutFormatting.Reset()
    f.formattingTemplate.Reset()
    f.lastMatchPosition = 0
    f.currentFormattingPattern = ""
    f.prefixBeforeNationalNumber.Reset()
    f.extractedNationalPrefix = ""
    f.nationalNumber.Reset()
    f.ableToFormat = true
    f.inputHasFormatting = false
    f.positionToRemember = 0
    f.origingalPosition = 0
    f.isCompleteNumber = false
    f.isExpectingCountryCallingCode = false
    f.possibleFormats = []*NumberFormat{}
    f.shouldAddSpaceAfterNationalPrefix = false
    if f.currentMetadata == f.defaultMetadata {
        f.currentMetadata = getMetadataForRegionCode(f.defaultCountry)
    }
}

func (f *AsYouTypeFormatter) inputDigitHelper(nextChar rune) string {
    // Note that formattingTemplate is not guaranteed to have a value, it could be empty, e.g.
    // when the next digit is entered after extracting an IDD or NDD.
    found, _, _ := regexFind(DIGIT_PATTERN, f.formattingTemplate.String(), f.lastMatchPosition)
    if found {
        tempTemplate := f.formattingTemplate.String()
        start := DIGIT_PATTERN.FindStringIndex(tempTemplate)[0]
        found := DIGIT_PATTERN.FindString(tempTemplate)
        if found != "" {
            tempTemplate = strings.Replace(tempTemplate, found, string(nextChar), 1)
        }

        formattingTemplate := f.formattingTemplate.String()
        f.formattingTemplate.Reset()
        f.formattingTemplate.WriteString(tempTemplate)
        f.formattingTemplate.WriteString(formattingTemplate[len(tempTemplate):])
        f.lastMatchPosition = start

        return f.formattingTemplate.String()[:f.lastMatchPosition + 1]
    } else {
        if len(f.possibleFormats) == 1 {
            // More digits are entered than we could handle, and there are no other valid patterns to
            // try.
            f.ableToFormat = false
        } // else, we just reset the formatting pattern.
        f.currentFormattingPattern = ""
        return f.accruedInput.String()
    }
}

// Some national prefixes are a substring of others. If extracting the shorter NDD doesn't result
// in a number we can format, we try to see if we can extract a longer version here.
func (f *AsYouTypeFormatter) ableToExtractLongerNdd() bool {
    if len(f.extractedNationalPrefix) > 0 {
        // Put the extracted NDD back to the national number before attempting to extract a new NDD.
        f.nationalNumber.InsertString(0, f.extractedNationalPrefix)
        // Remove the previously extracted NDD from prefixBeforeNationalNumber. We cannot simply set
        // it to empty string because people sometimes incorrectly enter national prefix after the
        // country code, e.g. +44 (0)20-1234-5678.  
        indexOfPreviousNdd := strings.LastIndex(f.prefixBeforeNationalNumber.String(), f.extractedNationalPrefix)
        prefixBeforeNationalNumber := f.prefixBeforeNationalNumber.String()
        f.prefixBeforeNationalNumber.Reset()
        f.prefixBeforeNationalNumber.WriteString(prefixBeforeNationalNumber[:indexOfPreviousNdd])
    }

    return f.extractedNationalPrefix != f.removeNationalPrefixFromNationalNumber()
}

// Returns the national prefix extracted, or an empty string if it is not present.
func (f *AsYouTypeFormatter) removeNationalPrefixFromNationalNumber() string {
    startOfNationalNumber := 0
    if f.isNanpaNumberWithNationalPrefix() {
        startOfNationalNumber = 1
        f.prefixBeforeNationalNumber.WriteString("1")
        f.prefixBeforeNationalNumber.WriteRune(SEPARATOR_BEFORE_NATIONAL_NUMBER)
        f.isCompleteNumber = true
    } else if f.currentMetadata.NationalPrefixForParsing != nil {
        nationalPrefixForParsing := regexFor(*f.currentMetadata.NationalPrefixForParsing)
        // Since some national prefix patterns are entirely optional, check that a national prefix
        // could actually be extracted.
        lookingAt, _, end := lookingAt(nationalPrefixForParsing, f.nationalNumber.String())
        if lookingAt && end > 0 {
            // When the national prefix is detected, we use international formatting rules instead of
            // national ones, because national formatting rules could contain local formatting rules
            // for numbers entered without area code.
            f.isCompleteNumber = true
            startOfNationalNumber = end
            f.prefixBeforeNationalNumber.WriteString(f.nationalNumber.String()[:startOfNationalNumber])
        }
    }
    nationalPrefix := f.nationalNumber.String()[:startOfNationalNumber]
    nationalNumber := f.nationalNumber.String()
    f.nationalNumber.Reset()
    f.nationalNumber.WriteString(nationalNumber[startOfNationalNumber:])
    return nationalPrefix
}

/**
* Returns true if the current country is a NANPA country and the national number begins with
* the national prefix.
*/
func (f *AsYouTypeFormatter) isNanpaNumberWithNationalPrefix() bool {
    // For NANPA numbers beginning with 1[2-9], treat the 1 as the national prefix. The reason is
    // that national significant numbers in NANPA always start with [2-9] after the national prefix.
    // Numbers beginning with 1[01] can only be short/emergency numbers, which don't need the
    // national prefix.
    return f.currentMetadata.GetCountryCode() == 1 && f.nationalNumber.Len() > 0 && f.nationalNumber.String()[0] == '1' && f.nationalNumber.Len() > 2 && f.nationalNumber.String()[1] != '0' && f.nationalNumber.String()[1] != '1'
}

func regexFind(pattern *regexp.Regexp, stringToFind string, startIdx int) (bool, int, int) {
    idxs := pattern.FindStringIndex(stringToFind[startIdx:])
    if len(idxs) > 0 {
        return true, idxs[0], idxs[1]
    }

    return false, -1, -1
}

// lookingAt Attempts to match the input sequence, starting at the beginning of the region, against the pattern
// without requiring that the entire region be matched and returns a value indicating whether or not the
// pattern was matched and the start and end indexes of where the match was found.
func lookingAt(pattern *regexp.Regexp, stringToMatch string) (bool, int, int) {
    idxs := pattern.FindStringIndex(stringToMatch)
    var matched bool
    var start int = -1
    var end   int = -1

    matched = len(idxs) > 0 && idxs[0] == 0
    if len(idxs) > 0 {
        start = idxs[0]
    }
    if len(idxs) > 1 {
        end = idxs[1]
    }

    return matched, start, end
}
