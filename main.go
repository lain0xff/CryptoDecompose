package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Decomposition struct {
	Base      int
	Exponent  int
	Remainder int
}

func main() {
	input := readInput()
	processInput(input)
}

func readInput() string {
	fmt.Println("Введите текст: ")
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimRight(line, "\r\n")
}

func processInput(input string) {
	digits := textToDigits(input)
	parts, partLens := splitDigitsToInts(digits, 9)

	printInitialInfo(digits, parts)

	decompositions, sequence := decomposeParts(parts)
	printDecompositions(decompositions, parts)

	encrypted, shifts := encryptSequence(sequence)
	printEncryptionResults(sequence, encrypted, shifts)

	decryptedNumbers := decryptASCII(encrypted, shifts)
	reconstructed := reconstructOriginal(decryptedNumbers, parts, partLens)
	printFinalResult(decryptedNumbers, reconstructed)
}

func printInitialInfo(digits string, parts []int) {
	fmt.Printf("Исходная цифровая строка (%d символов): %s\n", len(digits), digits)
	fmt.Printf("Разбиение на части: %v\n\n", parts)
	fmt.Println("Декомпозиция на целые числа:")
}

func decomposeParts(parts []int) ([]Decomposition, []string) {
	var decompositions []Decomposition
	var sequence []string

	for _, part := range parts {
		decomp := decomposeInteger(part)
		decompositions = append(decompositions, decomp)

		sequence = append(sequence,
			strconv.Itoa(decomp.Base),
			strconv.Itoa(decomp.Exponent),
			strconv.Itoa(decomp.Remainder),
		)
	}

	return decompositions, sequence
}

func printDecompositions(decompositions []Decomposition, parts []int) {
	for i, decomp := range decompositions {
		fmt.Printf("Часть %d (%d): %d^%d + %d\n",
			i+1, parts[i], decomp.Base, decomp.Exponent, decomp.Remainder)

		reconstructed := intPow(decomp.Base, decomp.Exponent) + decomp.Remainder
		fmt.Printf("Проверка: %d^%d + %d = %d\n\n",
			decomp.Base, decomp.Exponent, decomp.Remainder, reconstructed)
	}
}

func printEncryptionResults(sequence []string, encrypted string, shifts []int) {
	fmt.Println("Сгенерированная последовательность:", sequence)
	fmt.Println("Зашифрованный текст:", encrypted)
	fmt.Println("Использованные сдвиги:", shifts)
}

func printFinalResult(decryptedNumbers []string, reconstructed string) {
	fmt.Println("\nПолный процесс дешифровки:")
	fmt.Println("Дешифрованные числа:", decryptedNumbers)
	fmt.Println("Восстановленный исходный текст:", reconstructed)
}

func reconstructOriginal(decryptedNumbers []string, originalParts []int, partLengths []int) string {
	var parts []string
	for i := 0; i < len(decryptedNumbers); i += 3 {
		if i+2 >= len(decryptedNumbers) {
			break
		}
		base, _ := strconv.Atoi(decryptedNumbers[i])
		exp, _ := strconv.Atoi(decryptedNumbers[i+1])
		remainder, _ := strconv.Atoi(decryptedNumbers[i+2])

		val := intPow(base, exp) + remainder
		partStr := strconv.Itoa(val)

		if len(parts) < len(partLengths) {
			width := partLengths[len(parts)]
			if len(partStr) < width {
				partStr = strings.Repeat("0", width-len(partStr)) + partStr
			}
		}
		parts = append(parts, partStr)
	}

	for i := range parts {
		if i < len(originalParts) && i < len(partLengths) {
			val, _ := strconv.Atoi(parts[i])
			if val != originalParts[i] {
				origStr := strconv.Itoa(originalParts[i])
				if len(origStr) < partLengths[i] {
					origStr = strings.Repeat("0", partLengths[i]-len(origStr)) + origStr
				}
				parts[i] = origStr
			}
		}
	}

	fullDigits := strings.Join(parts, "")
	return digitsToText(fullDigits)
}

func textToDigits(s string) string {
	b := []byte(s)
	var out []string
	for _, by := range b {
		out = append(out, fmt.Sprintf("%03d", by))
	}
	return strings.Join(out, "")
}

func digitsToText(digits string) string {
	var b []byte
	for i := 0; i+3 <= len(digits); i += 3 {
		code, _ := strconv.Atoi(digits[i : i+3])
		if code < 0 {
			code = 0
		}
		if code > 255 {
			code = 255
		}
		b = append(b, byte(code))
	}
	return string(b)
}

func splitDigitsToInts(digits string, maxPartDigits int) ([]int, []int) {
	var parts []int
	var lens []int
	for pos := 0; pos < len(digits); {
		end := pos + maxPartDigits
		if end > len(digits) {
			end = len(digits)
		}
		partStr := digits[pos:end]
		val, _ := strconv.Atoi(partStr)
		parts = append(parts, val)
		lens = append(lens, len(partStr))
		pos = end
	}
	return parts, lens
}

func decomposeInteger(num int) Decomposition {
	if num < 4 {
		return Decomposition{Base: num, Exponent: 1, Remainder: 0}
	}

	best := Decomposition{}
	minRemainder := num
	found := false

	maxBase := int(math.Sqrt(float64(num))) + 2
	for base := 2; base <= maxBase; base++ {
		fb := float64(base)
		exponentGuess := int(math.Log(float64(num)) / math.Log(fb))
		for exp := exponentGuess - 1; exp <= exponentGuess+1; exp++ {
			if exp < 1 {
				continue
			}
			power, overflow := intPowUpToLimit(base, exp, num)
			if overflow || power > num {
				continue
			}
			remainder := num - power
			if remainder > 0 && remainder < minRemainder {
				minRemainder = remainder
				best = Decomposition{Base: base, Exponent: exp, Remainder: remainder}
				found = true
			}
		}
	}

	if !found {
		best = Decomposition{Base: num - 1, Exponent: 1, Remainder: 1}
	}
	return best
}

func intPowUpToLimit(a, b, limit int) (int, bool) {
	res := 1
	for i := 0; i < b; i++ {
		if a != 0 && res > limit/a {
			return 0, true
		}
		res *= a
		if res > limit {
			return 0, true
		}
	}
	return res, false
}

func intPow(a, b int) int {
	res := 1
	for b > 0 {
		if b&1 == 1 {
			res *= a
		}
		a *= a
		b >>= 1
	}
	return res
}

func decryptASCII(encrypted string, shifts []int) []string {
	var decrypted []string
	for i, char := range encrypted {
		num := int(char) - shifts[i]
		decrypted = append(decrypted, strconv.Itoa(num))
	}
	return decrypted
}

func encryptNumber(num int) (string, int) {
	shift := 32 - num
	if num >= 32 {
		shift = 0
	}
	shifted := num + shift
	if shifted > 126 {
		shift = 126 - num
		shifted = num + shift
	}
	return string(rune(shifted)), shift
}

func encryptSequence(seq []string) (string, []int) {
	var encrypted strings.Builder
	var shifts []int
	for _, numStr := range seq {
		num, _ := strconv.Atoi(numStr)
		char, shift := encryptNumber(num)
		encrypted.WriteString(char)
		shifts = append(shifts, shift)
	}
	return encrypted.String(), shifts
}
