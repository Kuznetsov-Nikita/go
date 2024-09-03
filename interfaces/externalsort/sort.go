//go:build !solution

package externalsort

import (
	"bufio"
	"container/heap"
	"io"
	"os"
	"sort"
)

type lineReader struct {
	reader *bufio.Reader
}

func (lr *lineReader) ReadLine() (string, error) {
	line, err := lr.reader.ReadString('\n')
	if err != nil {
		if len(line) != 0 {
			if line[len(line)-1] == '\n' {
				line = line[:len(line)-1]
			}
			return line, nil
		}
		return "", err
	}
	if line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	return line, nil
}

func NewReader(r io.Reader) LineReader {
	return &lineReader{reader: bufio.NewReader(r)}
}

type lineWriter struct {
	writer io.Writer
}

func (lw *lineWriter) Write(line string) error {
	line = line + "\n"
	_, err := io.WriteString(lw.writer, line)
	return err
}

func NewWriter(w io.Writer) LineWriter {
	return &lineWriter{writer: w}
}

type LineHeap []string

func (h LineHeap) Len() int           { return len(h) }
func (h LineHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h LineHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *LineHeap) Push(x interface{}) {
	*h = append(*h, x.(string))
	heap.Fix(h, h.Len()-1)
}

func (h *LineHeap) Pop() interface{} {
	x := (*h)[h.Len()-1]
	t := (*h)[0]
	*h = (*h)[1:h.Len()]
	if h.Len() == 0 {
		return x
	}
	(*h)[h.Len()-1] = t
	heap.Fix(h, h.Len()-1)
	return x
}

func (h *LineHeap) Top() interface{} {
	return (*h)[0]
}

func Merge(w LineWriter, readers ...LineReader) error {
	readerLines := make([]string, len(readers))

	h := &LineHeap{}
	heap.Init(h)

	for i, reader := range readers {
		line, err := reader.ReadLine()
		if err == io.EOF {
			continue
		}
		if err != nil {
			return err
		}
		readerLines[i] = line
		heap.Push(h, line)
	}

	for h.Len() > 0 {
		minLine := heap.Pop(h).(string)
		if err := w.Write(minLine); err != nil {
			return err
		}

		for i, reader := range readers {
			if readerLines[i] == minLine {
				line, err := reader.ReadLine()
				if err == io.EOF {
					readerLines[i] = ""
					continue
				} else if err != nil {
					return err
				} else {
					readerLines[i] = line
					h.Push(line)
				}
				break
			}
		}
	}

	return nil
}

func Sort(w io.Writer, in ...string) error {
	if len(in) == 0 {
		return nil
	}
	if len(in) == 1 {
		file, err := os.Open(in[0])
		if err != nil {
			return err
		}
		defer file.Close()
		reader := NewReader(file)

		lines := make([]string, 0)
		for {
			line, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			lines = append(lines, line)
		}
		sort.Strings(lines)
		writer := NewWriter(w)
		for _, line := range lines {
			writer.Write(line)
		}
		return nil
	}

	mid := len(in) / 2
	leftFiles := in[:mid]
	rightFiles := in[mid:]

	leftTempFile, err := os.CreateTemp("", "left*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(leftTempFile.Name())

	rightTempFile, err := os.CreateTemp("", "right*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(rightTempFile.Name())

	if err := Sort(leftTempFile, leftFiles...); err != nil {
		return err
	}

	if err := Sort(rightTempFile, rightFiles...); err != nil {
		return err
	}

	leftTempFile.Seek(0, 0)
	rightTempFile.Seek(0, 0)
	leftReader := NewReader(leftTempFile)
	rightReader := NewReader(rightTempFile)

	return Merge(NewWriter(w), leftReader, rightReader)
}
