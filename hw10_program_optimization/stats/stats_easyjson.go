// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package stats

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonE3ab7953DecodeGithubComMaksimIschenkoHwOtusGolangHw10ProgramOptimizationStats(in *jlexer.Lexer, out *LightUser) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Email":
			out.Email = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonE3ab7953EncodeGithubComMaksimIschenkoHwOtusGolangHw10ProgramOptimizationStats(out *jwriter.Writer, in LightUser) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Email\":"
		out.RawString(prefix[1:])
		out.String(string(in.Email))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LightUser) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE3ab7953EncodeGithubComMaksimIschenkoHwOtusGolangHw10ProgramOptimizationStats(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LightUser) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE3ab7953EncodeGithubComMaksimIschenkoHwOtusGolangHw10ProgramOptimizationStats(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LightUser) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE3ab7953DecodeGithubComMaksimIschenkoHwOtusGolangHw10ProgramOptimizationStats(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LightUser) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE3ab7953DecodeGithubComMaksimIschenkoHwOtusGolangHw10ProgramOptimizationStats(l, v)
}
