package editor

import (
	"fmt"
	"slices"
	"strings"
)

type Body struct {
	chunks     []string
	chunkLimit int
}

func NewBody(limit int) Body {
	return Body{
		chunks:     make([]string, 0),
		chunkLimit: limit,
	}
}

func (b *Body) Append(text string) {
	b.chunks = append(b.chunks, text)
}

func (b *Body) Insert(text string, index int) {
	i := -1
	var ch string = ""
	for i, ch = range b.chunks {
		if index >= len(ch) {
			index -= len(ch)
		} else {
			if index == 0 {
				b.chunks = slices.Insert(b.chunks, i, text)
			} else {
				if len(b.chunks[i])+len(text) > b.chunkLimit {
					post := b.chunks[i][index:]
					b.chunks[i] = b.chunks[i][:index]

					b.chunks = slices.Insert(b.chunks, i+1, text)
					b.chunks = slices.Insert(b.chunks, i+2, post)
				} else {
					b.chunks[i] = b.chunks[i][:index] + text + b.chunks[i][index:]
				}
			}
			index = -1
			break
		}
	}

	//end - edge case
	if index == 0 {
		if i < 0 {
			b.Append(text)
			return
		} else {
			//append to the last chunk
			if len(b.chunks[i])+len(text) <= b.chunkLimit {
				b.chunks[i] += text
			} else {
				b.Append(text)
			}
		}

	}
}

func (b *Body) Delete(index, length int) {
	for i, ch := range b.chunks {
		if index >= len(ch) {
			index -= len(ch)
		} else {
			// if the delete doesnt end with this chunk
			if len(ch)-index <= length {
				b.chunks[i] = ch[:index]
				length -= len(ch) - index
				index = 0
			} else {
				b.chunks[i] = ch[:index] + ch[index+length:]
				break
			}
		}
	}

	//delete empty strings from the slice
	for i := 0; i < len(b.chunks); i++ {
		if len(b.chunks[i]) == 0 {
			b.chunks = slices.Concat(b.chunks[:i], b.chunks[i+1:])
			i--
		}
	}
}

func (b Body) String() string {
	return fmt.Sprint(len(b.chunks)) + "[\n" + strings.Join(b.chunks, "\n") + "\n]"
}
