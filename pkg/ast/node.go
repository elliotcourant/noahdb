package ast

import (
	"crypto/sha1"
	"hash"
	"io"
	"reflect"
	"strings"
)

type FingerprintContext interface {
	WriteString(string)
}

type FingerprintHashContext struct {
	hash hash.Hash
}

func NewFingerprintHashContext() *FingerprintHashContext {
	return &FingerprintHashContext{hash: sha1.New()}
}

func (ctx FingerprintHashContext) WriteString(str string) {
	io.WriteString(ctx.hash, str)
}

func (ctx FingerprintHashContext) Sum() []byte {
	return ctx.hash.Sum(nil)
}

// ...

type FingerprintSubContext struct {
	parts []string
}

type FingerprintSubContextSlice []FingerprintSubContext

func NewFingerprintSubContext() *FingerprintSubContext {
	return &FingerprintSubContext{parts: []string{}}
}

func (ctx *FingerprintSubContext) WriteString(str string) {
	ctx.parts = append(ctx.parts, str)
}

func (ctx FingerprintSubContext) Sum() []string {
	return ctx.parts
}

func (ctx FingerprintSubContext) Equal(other FingerprintSubContext) bool {
	return reflect.DeepEqual(ctx, other)
}

func (p FingerprintSubContextSlice) Len() int {
	return len(p)
}

func (p FingerprintSubContextSlice) Less(i, j int) bool {
	left := strings.Join(p[i].parts, "")
	right := strings.Join(p[j].parts, "")
	return left < right
}

func (p FingerprintSubContextSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *FingerprintSubContextSlice) AddIfUnique(ctx FingerprintSubContext) {
	for _, existing := range *p {
		if ctx.Equal(existing) {
			return
		}
	}

	*p = append(*p, ctx)
}

// ...

type Node interface {
	Deparse(ctx Context) (string, error)
	Fingerprint(FingerprintContext, Node, string)
}

func Deparse(node Node) (string, error) {
	return node.Deparse(Context_None)
}

func deparseNodeList(nodes []Node, ctx Context) ([]string, error) {
	out := make([]string, len(nodes))
	for i, node := range nodes {
		if str, err := node.Deparse(ctx); err != nil {
			return nil, err
		} else {
			out[i] = str
		}
	}
	return out, nil
}
