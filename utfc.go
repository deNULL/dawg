package dawg

type UtfcFollower interface {
	Follow(byte) bool
}

type byteBuffer struct {
	buf []byte
}

func (bb *byteBuffer) Follow(b byte) bool {
	bb.buf = append(bb.buf, b)
	return true
}

// All characters below this code point are considered Latin, so within this range the state of `offs` stays equal to 0
const maxLatinCp = 0x02FF

// All characters starting from this code encoded in long (21-bit) mode
const min21BitCp = 0x2600

// Offs always includes top 6 bits of the codepoint (it identifies the currently selected "alphabet")
const offsMask13Bit = 0xFFFFFF80 // Characters encoded using their lowest 7 bits
const offsMask21Bit = 0xFFFF8000 // Characters encoded using their lowest 15 bits

const markerAux = 0xc0   // => 1 byte encoding, auxiliary alphabet
const marker13Bit = 0x80 // => 2 byte encoding
const marker21Bit = 0xa0 // => 3 byte encoding
const markerExtra = 0xb0 // => 2 byte encoding, extra ranges

const marker0 = 0b10111011  // zero marker #0 => simple 00 byte replacement
const marker00 = 0b10111100 // zero marker #00 => double 00 byte replacement
const marker10 = 0b10111110 // zero marker #10 => read next value, then add 0, then read byte
const marker01 = 0b10111101 // zero marker #01 => read next value, then another byte, then add 0
const marker11 = 0b10111111 // zero marker #11 => read next value, then add two 0

//const offsInitAux = 0x00C0
const offsInitAux = 0x0410

// The subrange of the previous (auxiliary) alphabet is coded via 0b11000000.
// Unfortunately, a lot of alphabets are not aligned to 64-byte chunks in a good way,
// so we select different portions here to cover most frequently used characters.
var auxOffset = map[int]int{
	// 0x0000, Latin is a special case, it merges A-Z, a-z, 0-9, "-" and " " characters.
	0x0080: 0x00C0, // Latin-1 Supplement
	0x0380: 0x0391, // Greek
	0x0400: 0x0410, // Cyrillic
	0x0580: 0x05BE, // Hebrew
	0x0530: 0x0531, // Armenian
	0x0600: 0x060B, // Arabic
	0x0900: 0x090D, // Devangari
	0x0980: 0x098F, // Bengali
	0x0A00: 0x0A02, // Gurmukhi
	0x0A80: 0x0A8F, // Gujarati
	0x0B00: 0x0B0F, // Oriya
	0x0B80: 0x0B8E, // Tamil
	0x0C80: 0x0C8E, // Kannada
	0x0D00: 0x0D0E, // Malayalam
	0x0D80: 0x0D9B, // Sinhala
	0x0E00: 0x0E01, // Thai
	0x0E80: 0x0E81, // Lao
	0x0F00: 0x0F40, // Tibetan
	0x0F80: 0x0F90, // Tibetan
	0x1080: 0x10B0, // Georgian
	0x3000: 0x3040, // Hiragana
}

// Hiragana and Katakana
var rangeHK = []int{0x3000, 0x3100}

var rangesLatin = [][]int{
	{0x41, 0x5B}, {0x61, 0x7B}, {0x30, 0x3A},
	{0x20, 0x21}, {0x2D, 0x2E},
}
var rangesExtra = [][]int{
	{0x2000, 0x2600}, rangeHK, {0xFE00, 0xFE10}, {0x1F300, 0x1F5F0},
}

func inRanges(cp int, ranges [][]int) bool {
	for _, rng := range ranges {
		if rng[0] <= cp && cp < rng[1] {
			return true
		}
	}
	return false
}

func encodeRanges(cp int, ranges [][]int) int {
	v := 0
	for _, rng := range ranges {
		if rng[0] <= cp && cp < rng[1] {
			return v + (cp - rng[0])
		}
		v += rng[1] - rng[0]
	}
	return -1
}

func decodeRanges(v int, ranges [][]int) int {
	for _, rng := range ranges {
		if v < rng[1]-rng[0] {
			return rng[0] + v
		}
		v -= rng[1] - rng[0]
	}
	return -1
}

func getAuxOffset(offs int) int {
	if remappedOffs, ok := auxOffset[offs]; ok {
		return remappedOffs
	}
	return offs
}

// Calls Follow(byte) on the follower for each byte of the encoded string
// Aborts if Follow returns false
// Returns true if process was completed successfully, false otherwise
// Allows to abort encoding early and prevent unneeded memory allocations
func UtfcFollow(str string, f UtfcFollower) bool {
	// `offs`, `auxOffs` and `is21Bit` describe the current state.
	// `offs` is the start of the currently active window of Unicode codepoints.
	// `auxOffs` allows encoding 64 codepoints of the auxiliary alphabet.
	// `is21Bit` is true if we're in 21-bit mode (2-3 bytes per character).
	offs := 0
	auxOffs := offsInitAux
	is21Bit := false
	for _, ch := range str {
		cp := int(ch)
		// First, check if we can use 1-byte encoding via small 6-bit auxiliary alphabet
		if auxOffs == 0 && inRanges(cp, rangesLatin) {
			// 1 byte: auxiliary alphabet is Latin, rearrange it to fit 0xC0-0xFF range
			if !f.Follow(byte(markerAux | encodeRanges(cp, rangesLatin))) {
				return false
			}
		} else if auxOffs != 0 && cp >= auxOffs && cp <= auxOffs+0x3F {
			// 1 byte: code point is within the auxiliary alphabet (non-Latin)
			if !f.Follow(byte(markerAux | (cp - auxOffs))) {
				return false
			}
		} else
		// Second, there're 6 extra ranges (Hiragana, Katakana, and Emojis) that normally would require 3 bytes/character,
		// but are encoded with 2 (using range of codepoints 0x10FFFF-0x1FFFFF, which are not covered by Unicode)
		if inRanges(cp, rangesExtra) {
			newOffs := cp & offsMask13Bit
			if !is21Bit && newOffs == offs { // 1 byte: code point is within the current alphabet
				lo := byte(cp & 0x7F)
				if lo == 0 {
					if !f.Follow(marker0) {
						return false
					}
				} else {
					if !f.Follow(lo) {
						return false
					}
				}
			} else {
				// Reindex 6 ranges into a single contiguous one
				extra := encodeRanges(cp, rangesExtra)
				lo := byte(extra)
				if lo == 0 {
					if !f.Follow(marker10) || !f.Follow(byte(markerExtra|(1+(extra>>8)))) {
						return false
					}
				} else {
					if !f.Follow(byte(markerExtra|(1+(extra>>8)))) || !f.Follow(lo) {
						return false
					}
				}
				if cp >= rangeHK[0] && cp < rangeHK[1] { // Only Hiragana and Katakana change the current alphabet
					auxOffs = getAuxOffset(offs)
					offs = newOffs
					is21Bit = false
				}
			}
		} else
		// Lastly, check codepoint size to determine if it needs short (13-bit) or long (21-bit) mode
		if cp >= min21BitCp {
			// This code point requires 21 bit to encode
			// Characters up to 0x2800 can be encoded in shorter forms, so we start from 0
			cp -= min21BitCp
			newOffs := cp & offsMask21Bit
			if is21Bit && newOffs == offs { // 2 bytes: code point is within the current alphabet
				hi := byte((cp >> 8) & 0x7F)
				lo := byte(cp)
				if hi == 0 && lo == 0 {
					if !f.Follow(marker0) {
						return false
					}
				} else if hi == 0 {
					if !f.Follow(marker00) || !f.Follow(lo) {
						return false
					}
				} else if lo == 0 {
					if !f.Follow(marker10) || !f.Follow(hi) {
						return false
					}
				} else {
					if !f.Follow(hi) || !f.Follow(lo) {
						return false
					}
				}
			} else { // 3 bytes: we need to switch to the new alphabet
				hi := byte(cp >> 8)
				lo := byte(cp)
				if hi == 0 && lo == 0 {
					if !f.Follow(marker11) || !f.Follow(byte(marker21Bit|(cp>>16))) {
						return false
					}
				} else if hi == 0 {
					if !f.Follow(marker10) || !f.Follow(byte(marker21Bit|(cp>>16))) || !f.Follow(lo) {
						return false
					}
				} else if lo == 0 {
					if !f.Follow(marker01) || !f.Follow(byte(marker21Bit|(cp>>16))) || !f.Follow(hi) {
						return false
					}
				} else {
					if !f.Follow(byte(marker21Bit|(cp>>16))) || !f.Follow(hi) || !f.Follow(lo) {
						return false
					}
				}
				auxOffs = offs
				offs = newOffs
				is21Bit = true
			}
		} else { // This code point requires max 13 bits to encode
			newOffs := cp & offsMask13Bit
			if !is21Bit && newOffs == offs { // 1 byte: code point is within the current alphabet
				lo := byte(cp & 0x7F)
				if lo == 0 {
					if !f.Follow(marker0) {
						return false
					}
				} else {
					if !f.Follow(lo) {
						return false
					}
				}
			} else { // Final case: we need 2 bytes for this character
				lo := byte(cp & 0xFF)
				if lo == 0 {
					if !f.Follow(marker10) || !f.Follow(byte(marker13Bit|(cp>>8))) {
						return false
					}
				} else {
					if !f.Follow(byte(marker13Bit|(cp>>8))) || !f.Follow(lo) {
						return false
					}
				}
				auxOffs = getAuxOffset(offs)
				if cp <= maxLatinCp {
					offs = 0
				} else {
					offs = newOffs
				}
				is21Bit = false
			}
		}
	}
	return true
}

var sharedBuffer = &byteBuffer{
	buf: make([]byte, 0, 1024),
}

// UtfcEncode converts string to an UTF-C byte array
func UtfcEncode(str string) []byte {
	sharedBuffer.buf = sharedBuffer.buf[:0]
	UtfcFollow(str, sharedBuffer)
	return sharedBuffer.buf
}

// UtfcDecode converts UTF-C byte array to a string
func UtfcDecode(buf []byte) string {
	offs := 0
	auxOffs := offsInitAux
	is21Bit := false
	str := ""
	i := 0
	for i < len(buf) {
		zm := 0
		if buf[i] >= marker0 && buf[i] <= marker11 { // Decode zero-marker
			zm = int(buf[i])
			i++
		}
		cp := 0
		if zm != marker0 && zm != marker00 {
			cp = int(buf[i])
			i++
		}
		if (cp & markerAux) == markerAux {
			if auxOffs == 0 {
				cp = decodeRanges(cp^markerAux, rangesLatin)
			} else {
				cp = auxOffs + (cp ^ markerAux)
			}
		} else if (cp&markerExtra) == markerExtra && (cp^markerExtra) != 0 {
			lo := 0
			if zm != marker10 {
				lo = int(buf[i])
				i++
			}
			cp = decodeRanges(((cp^markerExtra)-1)<<8|lo, rangesExtra)
			if cp >= rangeHK[0] && cp < rangeHK[1] {
				auxOffs = getAuxOffset(offs)
				offs = cp & offsMask13Bit
				is21Bit = false
			}
		} else if (cp & marker21Bit) == marker21Bit {
			hi := 0
			if zm != marker10 && zm != marker11 {
				hi = int(buf[i])
				i++
			}
			lo := 0
			if zm != marker01 && zm != marker11 {
				lo = int(buf[i])
				i++
			}
			cp = (cp^marker21Bit)<<16 | hi<<8 | lo
			auxOffs = offs
			offs = cp & offsMask21Bit
			is21Bit = true
			cp += min21BitCp
		} else if (cp & marker13Bit) == marker13Bit {
			hi := 0
			if zm != marker10 {
				hi = int(buf[i])
				i++
			}
			cp = (cp^marker13Bit)<<8 | hi
			auxOffs = getAuxOffset(offs)
			if cp <= maxLatinCp {
				offs = 0
			} else {
				offs = cp & offsMask13Bit
			}
			is21Bit = false
		} else if is21Bit {
			hi := 0
			if zm != marker10 {
				hi = int(buf[i])
				i++
			}
			cp = min21BitCp + (offs | cp<<8 | hi)
		} else {
			cp = offs | cp
		}
		str += string(rune(cp))
	}
	return str
}
