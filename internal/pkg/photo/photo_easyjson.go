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

func easyjson49c0357aDecodeGithubComGoParkMailRu20191EscapadeInternalPkgPhoto(in *jlexer.Lexer, out *AwsPrivateConfig) {
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
func easyjson49c0357aEncodeGithubComGoParkMailRu20191EscapadeInternalPkgPhoto(out *jwriter.Writer, in AwsPrivateConfig) {
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
	easyjson49c0357aEncodeGithubComGoParkMailRu20191EscapadeInternalPkgPhoto(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AwsPrivateConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson49c0357aEncodeGithubComGoParkMailRu20191EscapadeInternalPkgPhoto(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AwsPrivateConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson49c0357aDecodeGithubComGoParkMailRu20191EscapadeInternalPkgPhoto(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AwsPrivateConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson49c0357aDecodeGithubComGoParkMailRu20191EscapadeInternalPkgPhoto(l, v)
}
func easyjson49c0357aDecodeGithubComGoParkMailRu20191EscapadeInternalPkgPhoto1(in *jlexer.Lexer, out *AwsPublicConfig) {
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
		case "maxFileSize":
			out.MaxFileSize = int64(in.Int64())
		case "allowedFileTypes":
			if in.IsNull() {
				in.Skip()
				out.AllowedFileTypes = nil
			} else {
				in.Delim('[')
				if out.AllowedFileTypes == nil {
					if !in.IsDelim(']') {
						out.AllowedFileTypes = make([]string, 0, 4)
					} else {
						out.AllowedFileTypes = []string{}
					}
				} else {
					out.AllowedFileTypes = (out.AllowedFileTypes)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.AllowedFileTypes = append(out.AllowedFileTypes, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson49c0357aEncodeGithubComGoParkMailRu20191EscapadeInternalPkgPhoto1(out *jwriter.Writer, in AwsPublicConfig) {
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
	{
		const prefix string = ",\"maxFileSize\":"
		out.RawString(prefix)
		out.Int64(int64(in.MaxFileSize))
	}
	{
		const prefix string = ",\"allowedFileTypes\":"
		out.RawString(prefix)
		if in.AllowedFileTypes == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.AllowedFileTypes {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AwsPublicConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson49c0357aEncodeGithubComGoParkMailRu20191EscapadeInternalPkgPhoto1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AwsPublicConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson49c0357aEncodeGithubComGoParkMailRu20191EscapadeInternalPkgPhoto1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AwsPublicConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson49c0357aDecodeGithubComGoParkMailRu20191EscapadeInternalPkgPhoto1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AwsPublicConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson49c0357aDecodeGithubComGoParkMailRu20191EscapadeInternalPkgPhoto1(l, v)
}
