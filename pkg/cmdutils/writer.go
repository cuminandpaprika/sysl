package cmdutils

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type SequenceDiagramWriter struct {
	Ind            int
	AtBeginOfLine  bool
	AutogenWarning bool
	Active         map[string]int
	Properties     []string
	Head           bytes.Buffer
	Body           bytes.Buffer
}

func MakeSequenceDiagramWriter(autogenWarning bool, properties ...string) *SequenceDiagramWriter {
	p := make([]string, 0, len(properties))
	p = append(p, properties...)
	return &SequenceDiagramWriter{
		AtBeginOfLine:  true,
		AutogenWarning: autogenWarning,
		Properties:     p,
		Active:         map[string]int{},
	}
}

func (s *SequenceDiagramWriter) WriteTo(w io.Writer) (n int64, err error) {
	i, err := fmt.Fprint(w, s.String())

	return int64(i), err
}

func (s *SequenceDiagramWriter) Write(p []byte) (n int, err error) {
	newline := []byte("\n")
	newlines := bytes.Count(p, newline)
	if newlines == 0 {
		if s.AtBeginOfLine {
			s.WriteIndent()
		}
		n, err = s.Body.Write(p)
		if err == nil {
			s.AtBeginOfLine = false
		}
		return n, err
	}

	frags := bytes.SplitN(p, newline, newlines+1)

	for i, frag := range frags {
		if s.AtBeginOfLine && len(frag) != 0 {
			s.WriteIndent()
		}
		nn, err := s.Body.Write(frag)
		if err != nil {
			return n, err
		}
		s.AtBeginOfLine = false
		n += nn
		if i+1 < len(frags) {
			if err := s.WriteByte('\n'); err != nil {
				return n, err
			}
			n++
		}
	}
	s.AtBeginOfLine = len(frags[len(frags)-1]) == 0
	return n, nil
}

func (s *SequenceDiagramWriter) WriteString(v string) (n int, err error) {
	return s.Write([]byte(v))
}

func (s *SequenceDiagramWriter) WriteByte(c byte) error {
	if s.AtBeginOfLine {
		s.WriteIndent()
	}

	err := s.Body.WriteByte(c)
	if err == nil {
		s.AtBeginOfLine = c == '\n'
	}

	return err
}

func (s *SequenceDiagramWriter) WriteHead(v string) (int, error) {
	return fmt.Fprintln(&s.Head, v)
}

func (s *SequenceDiagramWriter) Indent() {
	s.Ind++
}

func (s *SequenceDiagramWriter) Unindent() {
	if s.Ind == 0 {
		panic("SequenceDiagramWriter unindent too far")
	}
	s.Ind--
}

func (s *SequenceDiagramWriter) Activate(agent string) {
	s.Active[agent]++

	fmt.Fprintf(s, "activate %s\n", agent)
}

func (s *SequenceDiagramWriter) Activated(agent string, suppressed bool) func() {
	if !suppressed {
		s.Activate(agent)
	}

	active := !suppressed

	return func() {
		if active {
			active = false
			s.Deactivate(agent)
		}
	}
}

func (s *SequenceDiagramWriter) Deactivate(agent string) {
	if v, ok := s.Active[agent]; !ok || v == 0 {
		return
	}

	s.Active[agent]--
	fmt.Fprintf(s, "deactivate %s\n", agent)

	if s.Active[agent] == 0 {
		delete(s.Active, agent)
	}
}

func (s *SequenceDiagramWriter) WriteIndent() {
	if !s.AtBeginOfLine {
		return
	}

	s.Body.WriteString(strings.Repeat(" ", s.Ind))
	s.AtBeginOfLine = false
}

func (s *SequenceDiagramWriter) String() string {
	if s.Body.Len() == 0 || s.Head.Len() == 0 {
		return ""
	}

	var sb strings.Builder
	if s.AutogenWarning {
		fmt.Fprintln(&sb, "''''''''''''''''''''''''''''''''''''''''''")
		fmt.Fprintln(&sb, "''                                      ''")
		fmt.Fprintln(&sb, "''  AUTOGENERATED CODE -- DO NOT EDIT!  ''")
		fmt.Fprintln(&sb, "''                                      ''")
		fmt.Fprintln(&sb, "''''''''''''''''''''''''''''''''''''''''''")
		fmt.Fprintln(&sb)
	}
	fmt.Fprintln(&sb, "@startuml")
	sb.WriteString(s.Head.String())
	for _, p := range s.Properties {
		fmt.Fprintln(&sb, p)
	}
	sb.WriteString(s.Body.String())
	fmt.Fprintln(&sb, "@enduml")

	return sb.String()
}