package unicode_name

import "unicode"

// AnalyzeName считает байты, Unicode-кодовые точки, латинские и кириллические
// буквы и возвращает байтовые позиции начала каждой Unicode-кодовой точки.
// Предполагается, что входная строка содержит корректный UTF-8.
func AnalyzeName(name string) (byteCount, runeCount, latinCount, cyrillicCount int, positions []int) {

	byteCount = len(name)

	if byteCount == 0 {
		return 0, 0, 0, 0, nil
	}

	for i, r := range name {
		runeCount += 1
		positions = append(positions, i)

		if unicode.IsLetter(r) {
			if unicode.Is(unicode.Latin, r) {
				latinCount += 1
			} else if unicode.Is(unicode.Cyrillic, r) {
				cyrillicCount += 1
			}
		}
	}

	return byteCount, runeCount, latinCount, cyrillicCount, positions

}
