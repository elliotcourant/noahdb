package ast

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/elliotcourant/timber"
	"github.com/stretchr/testify/assert"
	"runtime/debug"
	"testing"
)

type DeparseTest struct {
	Query                string
	Expected             string
	ExpectedParseError   string
	ExpectedCompileError string
}

type parsetreeList struct {
	Statements []Node
	Query      string
}

func (input parsetreeList) MarshalJSON() ([]byte, error) {
	return json.Marshal(input.Statements)
}

func (output *parsetreeList) UnmarshalJSON(input []byte) (err error) {
	var list []json.RawMessage
	err = json.Unmarshal([]byte(input), &list)
	if err != nil {
		return
	}

	for _, nodeJson := range list {
		var node Node
		node, err = UnmarshalNodeJSON(nodeJson)
		if err != nil {
			return
		}
		output.Statements = append(output.Statements, node)
	}

	return
}

func (input parsetreeList) Fingerprint() string {
	const fingerprintVersion uint = 2

	ctx := NewFingerprintHashContext()
	for _, node := range input.Statements {
		node.Fingerprint(ctx, nil, "")
	}

	return fmt.Sprintf("%02x%s", fingerprintVersion, hex.EncodeToString(ctx.Sum()))
}

func parse(input string, log bool) (t *parsetreeList, errr error) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			errr = r.(error)
		}
	}()
	jsonTree, err := parseToJSON(input)
	if err != nil {
		return nil, err
	}
	if log {
		timber.Debugf(" QUERY  | %s", input)
		timber.Debugf(" TREE   | %s", string(jsonTree))
	}

	// JSON unmarshalling can panic in edge cases we don't support yet. This is
	// still a *bug that needs to be fixed*, but this way the caller can expect an
	// error to be returned always, instead of a panic

	tree := &parsetreeList{}
	err = json.Unmarshal([]byte(jsonTree), &tree)
	tree.Query = input
	return tree, err
}

func DoTest(t *testing.T, test DeparseTest) {
	// First we want to parse the provided query.
	ast, err := parse(test.Query, true)
	if test.ExpectedParseError != "" {
		assert.EqualError(t, err, test.ExpectedParseError, "did not receive expected parse error")
	} else {
		if !assert.NoError(t, err, "received an unexpected error while parsing query") {
			t.FailNow()
		}
	}

	j, err := ast.MarshalJSON()
	if !assert.NoError(t, err) {
		panic(err)
	}
	assert.NotEmpty(t, j)

	if finger := ast.Fingerprint(); !assert.NotEmpty(t, finger) {
		panic("fingerprint is empty")
	}

	recompiled, err := ast.Statements[0].Deparse(Context_None)
	if test.ExpectedCompileError != "" {
		assert.EqualError(t, err, test.ExpectedCompileError, "did not receive the expected error when recompiling")
	} else {
		if !assert.NoError(t, err, "received an unexpected error while recompiling query") {
			t.FailNow()
		}
	}

	timber.Debugf("RESULT | %s\n", recompiled)

	_, err = parse(recompiled, false)
	if err != nil {
		t.Errorf("failed to parse recompiled query: %s", err)
		t.FailNow()
	}

	if test.Expected != "" {
		assert.Equal(t, test.Expected, recompiled, "compiled query did not match expected query")
	} else {
		assert.NotNil(t, recompiled, "compiled query is nil")
	}
}
