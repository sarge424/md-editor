package editor

import (
	"fmt"
	"slices"
)

type content struct {
	chunks    []string
	chunkSize int
}

func newContent(limit int) content {
	return content{
		chunks:    make([]string, 0),
		chunkSize: limit,
	}
}

func (c *content) append(text string) {
	c.chunks = append(c.chunks, text)
}

func (c *content) Insert(text string, index int) {
	i := -1
	var ch string = ""
	for i, ch = range c.chunks {
		if index >= len(ch) {
			index -= len(ch)
		} else {
			if index == 0 {
				c.chunks = slices.Insert(c.chunks, i, text)
			} else {
				if len(c.chunks[i])+len(text) > c.chunkSize {
					post := c.chunks[i][index:]
					c.chunks[i] = c.chunks[i][:index]

					c.chunks = slices.Insert(c.chunks, i+1, text)
					c.chunks = slices.Insert(c.chunks, i+2, post)
				} else {
					c.chunks[i] = c.chunks[i][:index] + text + c.chunks[i][index:]
				}
			}
			index = -1
			break
		}
	}

	//end - edge case
	if index == 0 {
		if i < 0 {
			c.append(text)
			return
		} else {
			//append to the last chunk
			if len(c.chunks[i])+len(text) <= c.chunkSize {
				c.chunks[i] += text
			} else {
				c.append(text)
			}
		}

	}
}

func (c *content) Delete(index, length int) {
	for i, ch := range c.chunks {
		if index >= len(ch) {
			index -= len(ch)
		} else {
			// if the delete doesnt end with this chunk
			if len(ch)-index <= length {
				c.chunks[i] = ch[:index]
				length -= len(ch) - index
				index = 0
			} else {
				c.chunks[i] = ch[:index] + ch[index+length:]
				break
			}
		}
	}

	//delete empty strings from the slice
	for i := 0; i < len(c.chunks); i++ {
		if len(c.chunks[i]) == 0 {
			c.chunks = slices.Concat(c.chunks[:i], c.chunks[i+1:])
			i--
		}
	}
}

func (c content) Get(index, length int) string {
	s := ""

	chunkStart := 0
	for _, ch := range c.chunks {
		if index < chunkStart+len(ch) {
			s += ch[max(0, index-chunkStart):min(index+length-chunkStart, len(ch))]
		}

		if index+length <= chunkStart {
			break
		}
	}

	return s
}

func (c content) String() string {
	s := ""
	s += fmt.Sprintln(len(c.chunks), "[")

	for _, ch := range c.chunks {
		s += fmt.Sprintf("--- -%d- ---\n", len(ch))
		s += fmt.Sprintln(ch)
	}

	return s + "]"
}
