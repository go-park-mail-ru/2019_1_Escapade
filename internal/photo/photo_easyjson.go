// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package photo

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

func easyjson49c0357aDecodeGithubComGoParkMailRu20191EscapadeInternalPhoto(in *jlexer.Lexer, out *AwsPrivateConfig) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "accessUrl":
			out.AccessURL = string(in.String())
		case "accessKey":
			out.AccessKey = string(in.String())
		case "secretUrl":
			out.SecretURL = string(in.String())
		case "secretKey":
			out.SecretKey = string(in.String())
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
func easyjson49c0357aEncodeGithubComGoParkMailRu20191EscapadeInternalPhoto(out *jwriter.Writer, in AwsPrivateConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"accessUrl\":"
		out.RawString(prefix[1:])
		out.String(string(in.AccessURL))
	}
	{
		const prefix string = ",\"accessKey\":"
		out.RawString(prefix)
		out.String(string(in.AccessKey))
	}
	{
		const prefix string = ",\"secretUrl\":"
		out.RawString(prefix)
		out.String(string(in.SecretURL))
	}
	{
		const prefix string = ",\"secretKey\":"
		out.RawString(prefix)
		out.String(string(in.SecretKey))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AwsPrivateConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson49c0357aEncodeGithubComGoParkMailRu20191EscapadeInternalPhoto(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AwsPrivateConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson49c0357aEncodeGithubComGoParkMailRu20191EscapadeInternalPhoto(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AwsPrivateConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson49c0357aDecodeGithubComGoParkMailRu20191EscapadeInternalPhoto(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AwsPrivateConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson49c0357aDecodeGithubComGoParkMailRu20191EscapadeInternalPhoto(l, v)
}
func easyjson49c0357aDecodeGithubComGoParkMailRu20191EscapadeInternalPhoto1(in *jlexer.Lexer, out *AwsPublicConfig) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "region":
			out.Region = string(in.String())
		case "endpoint":
			out.Endpoint = string(in.String())
		case "playersAvatarsStorage":
			out.PlayersAvatarsStorage = string(in.String())
		case "defaultAvatar":
			out.DefaultAvatar = string(in.String())
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
func easyjson49c0357aEncodeGithubComGoParkMailRu20191EscapadeInternalPhoto1(out *jwriter.Writer, in AwsPublicConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"region\":"
		out.RawString(prefix[1:])
		out.String(string(in.Region))
	}
	{
		const prefix string = ",\"endpoint\":"
		out.RawString(prefix)
		out.String(string(in.Endpoint))
	}
	{
		const prefix string = ",\"playersAvatarsStorage\":"
		out.RawString(prefix)
		out.String(string(in.PlayersAvatarsStorage))
	}
	{
		const prefix string = ",\"defaultAvatar\":"
		out.RawString(prefix)
		out.String(string(in.DefaultAvatar))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AwsPublicConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson49c0357aEncodeGithubComGoParkMailRu20191EscapadeInternalPhoto1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AwsPublicConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson49c0357aEncodeGithubComGoParkMailRu20191EscapadeInternalPhoto1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AwsPublicConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson49c0357aDecodeGithubComGoParkMailRu20191EscapadeInternalPhoto1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AwsPublicConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson49c0357aDecodeGithubComGoParkMailRu20191EscapadeInternalPhoto1(l, v)
}
