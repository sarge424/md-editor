package buffer

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Buffer struct {
	chunks []string
	Length int
	size   int

	newLines []int
}

type Parser struct {
	index int
	buf   Buffer

	chunkStart int
	chunkNo    int

	RowNo int
	Data  string
}

func New(size int) Buffer {
	b := Buffer{size: size}
	return b
}

func (b *Buffer) AddNewLines(t string, index int) {
	// add newlines
	for i, ch := range t {
		if ch == '\n' {
			b.newLines = append(b.newLines, i+index)
		}
	}

	slices.Sort(b.newLines)
}

func (b *Buffer) FixNewLines(index, length int) {
	// remove newlines in the deleted area
	b.newLines = slices.DeleteFunc(b.newLines, func(x int) bool {
		return index <= x && x < index+length
	})

	// offset newlines that come after the delete
	for i, ix := range b.newLines {
		if ix > index {
			b.newLines[i] -= length
		}
	}
}

func (b *Buffer) Insert(t string, index int) bool {
	t = strings.ReplaceAll(t, "\t", "    ")

	// if the buffer is empty, create the first chunk
	if b.Length == 0 {
		b.Length += len(t)

		// add newlines
		b.AddNewLines(t, index)

		// make sure chunks are not larger than max size
		for len(t) > 0 {
			en := min(b.size, len(t))
			b.chunks = append(b.chunks, t[:en])
			t = t[en:]
		}

		return true
	}

	ix := 0 // start index of current chunk
	for i, ch := range b.chunks {
		if ix+len(ch) >= index {
			// insert text into chunk
			b.chunks[i] = ch[:index-ix] + t + ch[index-ix:]
			b.Length += len(t)

			// add newlines
			b.AddNewLines(t, index)

			// split the chunk if it's too large
			for len(b.chunks[i]) > b.size {
				b.chunks = slices.Concat(
					b.chunks[:i],
					[]string{
						b.chunks[i][:b.size],
						b.chunks[i][b.size:],
					},
					b.chunks[i+1:],
				)
				i++
			}

			return true
		}

		ix += len(ch)
	}

	return false
}

func (b *Buffer) Delete(index, length int) {
	ix := 0
	delLen := 0

	for i, ch := range b.chunks {
		if index < ix+len(ch) {
			en := min(index+length-ix, len(ch))
			b.chunks[i] = ch[:index-ix] + ch[en:]

			diff := en - (index - ix)
			delLen += diff
			length -= diff
			index += diff

			b.Length -= diff

			if length == 0 {
				break
			}
		}

		ix += len(ch)
	}

	// fix newlines
	b.FixNewLines(index, delLen)

	// delete empty chunks
	b.chunks = slices.DeleteFunc(b.chunks, func(x string) bool {
		return len(x) == 0
	})
}

func (b *Buffer) LoadFile(fp string) error {
	dat, err := os.ReadFile(fp)
	if err != nil {
		fmt.Println("ERRORR")
		return err
	}
	b.Insert(string(dat), 0)

	return nil
}

func (b Buffer) String() string {
	return fmt.Sprint(b.chunks) + "\n" + fmt.Sprint(b.newLines)
}

func (b Buffer) Parser() Parser {
	return Parser{
		index: 0,
		buf:   b,

		chunkStart: 0,
		chunkNo:    0,

		RowNo: -1,
		Data:  "",
	}
}

func (p Parser) nextNL() int {
	return p.buf.newLines[p.RowNo]
}

func (p *Parser) Next() bool {
	p.Data = ""
	p.RowNo++
	if p.RowNo > len(p.buf.newLines) {
		return false
	}

	onLastRow := p.RowNo == len(p.buf.newLines)

	for i := p.chunkNo; i < len(p.buf.chunks); i++ {
		ch := p.buf.chunks[i]

		// cut out a substring to return as a row
		st := p.index

		// check for next newline only if NOT on last row
		en := len(ch)
		if !onLastRow {
			en = min(p.nextNL()-p.chunkStart, len(ch))
		}

		// add text to data
		p.Data += ch[st:en]
		// fmt.Println("tx is<", p.Data, ">", onLastRow, st, en)

		// move the start index
		p.index = en
		if !onLastRow && p.index == p.nextNL() {
			p.index++

			// increase the chunk start index
			if p.index > p.chunkStart+len(ch) {
				p.chunkStart += len(ch)
			}

			return true
		}

		// increase the chunk start index
		p.chunkStart += len(ch)
		p.chunkNo++
	}

	return true
}
