// +build tools

package gorules

import (
	"github.com/quasilyte/go-ruleguard/dsl/fluent"
)

func _(m fluent.Matcher) {
	// See http://golang.org/issue/36225
	m.Match(`json.NewDecoder($_).Decode($_)`).
		Report(`this json.Decoder usage is erroneous`)

	// See https://twitter.com/dvyukov/status/1174698980208513024
	m.Match(`type $x error`).
		Report(`error as underlying type is probably a mistake`).
		Suggest(`type $x struct { error }`)

	// From https://github.com/dominikh/go-tools/issues/582
	m.Match(`var()`).Report(`empty var() block`)
	m.Match(`const()`).Report(`empty const() block`)
	m.Match(`type()`).Report(`empty type() block`)

	m.Match(`fmt.Fprint(os.Stdout, $*args)`).Suggest(`fmt.Print($args)`)
	m.Match(`fmt.Fprintln(os.Stdout, $*args)`).Suggest(`fmt.Println($args)`)
	m.Match(`fmt.Fprintf(os.Stdout, $*args)`).Suggest(`fmt.Printf($args)`)

	m.Match(`strings.Count($s1, $s2) > 0`,
		`strings.Count($s1, $s2) >= 1`).
		Suggest(`strings.Contains($s1, $s2)`)
	m.Match(`strings.Count($s1, $s2) == 0`).
		Suggest(`!strings.Contains($s1, $s2)`)

	m.Match(`sort.Slice($s, func($i, $j int) bool { return $s[$i] < $s[$j] })`).
		Where(m["s"].Type.Is(`[]string`)).
		Suggest(`sort.Strings($s)`)

	m.Match(`sort.Slice($s, func($i, $j int) bool { return $s[$i] < $s[$j] })`).
		Where(m["s"].Type.Is(`[]int`)).
		Suggest(`sort.Ints($s)`)

	m.Match(`time.Duration($x) * time.Second`).
		Where(m["x"].Const).
		Suggest(`$x * time.Second`)

	m.Match(`int64(time.Since($t) / time.Microsecond)`).
		Suggest(`time.Since($t).Microseconds()`)
	m.Match(`int64(time.Since($t) / time.Millisecond)`).
		Suggest(`time.Since($t).Milliseconds()`)

	m.Match(`fmt.Sprint($x)`).
		Where(m["x"].Type.Implements(`fmt.Stringer`)).
		Suggest(`$x.String()`)

	m.Match(`os.Open(path.Join($*_))`,
		`ioutil.ReadFile(path.Join($*_))`,
		`$p := path.Join($*_); $_, $_ := os.Open($p)`,
		`$p := path.Join($*_); $_, $_ := ioutil.ReadFile($p)`).
		Report(`use filepath.Join for file paths`)

	m.Match(`select {case <-$ctx.Done(): return $ctx.Err(); default:}`).
		Where(m["ctx"].Type.Is(`context.Context`)).
		Suggest(`if err := $ctx.Err(); err != nil { return err }`)

	m.Match(`len($s) >= len($x) && $s[:len($x)] == $x`).
		Suggest(`strings.HasPrefix($s, $x)`)
	m.Match(`len($s) >= len($x) && $s[len($s)-len($x):] == $x`).
		Suggest(`strings.HasSuffix($s, $x)`)

	m.Match(`fmt.Sprintf("%s%s", $a, $b)`).
		Where(m["a"].Type.Is(`string`) && m["b"].Type.Is(`string`)).
		Suggest(`$a+$b`)
}

func gocriticWrapperFunc(m fluent.Matcher) {
	m.Match(`strings.SplitN($s, $sep, -1)`).Suggest(`strings.Split($s, $sep)`)
	m.Match(`strings.Replace($s, $old, $new, -1)`).Suggest(`strings.ReplaceAll($s, $old, $new)`)
	m.Match(`strings.TrimFunc($s, unicode.IsSpace)`).Suggest(`strings.TrimSpace($s)`)
	m.Match(`strings.Map(unicode.ToUpper, $s)`).Suggest(`strings.ToUpper($s)`)
	m.Match(`strings.Map(unicode.ToLower, $s)`).Suggest(`strings.ToLower($s)`)
	m.Match(`strings.Map(unicode.ToTitle, $s)`).Suggest(`strings.ToTitle($s)`)

	m.Match(`bytes.SplitN($s, $sep, -1)`).Suggest(`bytes.Split($s, $sep)`)
	m.Match(`bytes.Replace($s, $old, $new, -1)`).Suggest(`bytes.ReplaceAll($s, $old, $new)`)
	m.Match(`bytes.TrimFunc($s, unicode.IsSpace)`).Suggest(`bytes.TrimSpace($s)`)
	m.Match(`bytes.Map(unicode.ToUpper, $s)`).Suggest(`bytes.ToUpper($s)`)
	m.Match(`bytes.Map(unicode.ToLower, $s)`).Suggest(`bytes.ToLower($s)`)
	m.Match(`bytes.Map(unicode.ToTitle, $s)`).Suggest(`bytes.ToTitle($s)`)
}

func gocriticNilValReturn(m fluent.Matcher) {
	m.Match(`if $*_; $v == nil { return $v }`).
		Report(`returned expr is always nil; replace $v with nil`)
}

func gocriticBoolExprSimplify(m fluent.Matcher) {
	m.Match(`!!$x`).Suggest(`$x`)
	m.Match(`!($x != $y)`).Suggest(`$x == $y`)
	m.Match(`!($x == $y)`).Suggest(`$x != $y`)
}

func gocriticOffBy1(m fluent.Matcher) {
	m.Match(`$s[len($s)]`).
		Where(m["s"].Type.Is(`[]$elem`) && m["s"].Pure).
		Report(`index expr always panics; maybe you wanted $s[len($s)-1]?`)
}

func gocriticStringXBytes(m fluent.Matcher) {
	m.Match(`copy($b, []byte($s))`).
		Where(m["s"].Type.Is(`string`)).
		Suggest(`copy($b, $s)`)
}

func gocriticArgOrder(m fluent.Matcher) {
	m.Match(`strings.HasPrefix($s1, $s2)`).
		Where(m["s1"].Const && !m["s2"].Const).
		Suggest(`strings.HasPrefix($s2, $s1)`)

	m.Match(`strings.HasSuffix($s1, $s2)`).
		Where(m["s1"].Const && !m["s2"].Const).
		Suggest(`strings.HasPrefix($s2, $s1)`)

	m.Match(`strings.Contains($s1, $s2)`).
		Where(m["s1"].Const && !m["s2"].Const).
		Suggest(`strings.Contains($s2, $s1)`)
}

func gocriticBadCall(m fluent.Matcher) {
	m.Match(`strings.Replace($_, $_, $_, 0)`,
		`bytes.Replace($_, $_, $_, 0)`,
		`strings.SplitN($_, $_, 0)`,
		`bytes.SplitN($_, $_, 0)`).
		Report(`n=0 argument does nothing, maybe n=-1 is indended?`)

	m.Match(`append($_)`).
		Report(`append called with 1 argument does nothing`)
}

func gocriticDupArg(m fluent.Matcher) {
	m.Match(`math.Max($x, $x)`,
		`math.Min($x, $x)`,
		`strings.Contains($x, $x)`,
		`strings.Compare($x, $x)`,
		`strings.EqualFold($x, $x)`,
		`strings.HasPrefix($x, $x)`,
		`strings.HasSuffix($x, $x)`,
		`strings.Index($x, $x)`,
		`strings.LastIndex($x, $x)`,
		`strings.Split($x, $x)`,
		`strings.SplitAfter($x, $x)`,
		`strings.SplitAfterN($x, $x, $_)`,
		`strings.SplitN($x, $x, $_)`,
		`strings.ReplaceAll($_, $x, $x)`,
		`strings.Replace($_, $x, $x, $_)`,
		`bytes.Contains($x, $x)`,
		`bytes.Compare($x, $x)`,
		`bytes.Equal($x, $x)`,
		`bytes.EqualFold($x, $x)`,
		`bytes.HasPrefix($x, $x)`,
		`bytes.HasSuffix($x, $x)`,
		`bytes.Index($x, $x)`,
		`bytes.LastIndex($x, $x)`,
		`bytes.Split($x, $x)`,
		`bytes.SplitAfter($x, $x)`,
		`bytes.SplitAfterN($x, $x, $_)`,
		`bytes.SplitN($x, $x, $_)`,
		`bytes.ReplaceAll($_, $x, $x)`,
		`bytes.Replace($_, $x, $x, $_)`,
		`reflect.Copy($x, $x)`,
		`reflect.DeepEqual($x, $x)`,
		`types.Identical($x, $y)`,
		`io.Copy($x, $x)`,
		`copy($x, $x)`).
		Report(`suspicious duplicated args in $$`)
}

func gocriticDupSubExpr(m fluent.Matcher) {
	m.Match(`$x || $x`,
		`$x && $x`,
		`$x | $x`,
		`$x & $x`,
		`$x ^ $x`,
		`$x < $x`,
		`$x > $x`,
		`$x &^ $x`,
		`$x % $x`,
		`$x == $x`,
		`$x != $x`,
		`$x <= $x`,
		`$x >= $x`,
		`$x / $x`,
		`$x - $x`).
		Where(m["x"].Pure).
		Report(`suspicious identical LHS and RHS`)
}

func gocriticValSwap(m fluent.Matcher) {
	m.Match(`$tmp := $x; $x = $y; $y = $tmp`).Suggest(`$x, $y = $y, $x`)
}

func gocriticAssignOp(m fluent.Matcher) {
	// We need to define ++ and -- rules before the other,
	// so they can take a precedence.
	m.Match(`$x = $x + 1`).Suggest(`$x++`)
	m.Match(`$x = $x - 1`).Suggest(`$x--`)
	m.Match(`$x = $x * $y`).Suggest(`$x *= $y`)
	m.Match(`$x = $x / $y`).Suggest(`$x /= $y`)
	m.Match(`$x = $x % $y`).Suggest(`$x %= $y`)
	m.Match(`$x = $x + $y`).Suggest(`$x += $y`)
	m.Match(`$x = $x - $y`).Suggest(`$x -= $y`)
	m.Match(`$x = $x & $y`).Suggest(`$x &= $y`)
	m.Match(`$x = $x | $y`).Suggest(`$x |= $y`)
	m.Match(`$x = $x ^ $y`).Suggest(`$x ^= $y`)
	m.Match(`$x = $x << $y`).Suggest(`$x <<= $y`)
	m.Match(`$x = $x >> $y`).Suggest(`$x >>= $y`)
	m.Match(`$x = $x &^ $y`).Suggest(`$x &^= $y`)
}

func gocriticRegexpMust(m fluent.Matcher) {
	m.Match(`regexp.Compile($pat)`,
		`regexp.CompilePOSIX($pat)`).
		Where(m["pat"].Const).
		Report(`can use MustCompile for const patterns`)
}

func gocriticMapKey(m fluent.Matcher) {
	m.Match(`map[$_]$_{$*_, $k: $_, $*_, $k: $_, $*_}`).
		Where(m["k"].Pure).
		Report(`suspicious duplicate key $k`).
		At(m["k"])
}

func gocriticAppendCombine(m fluent.Matcher) {
	m.Match(`$dst = append($x, $a); $dst = append($x, $b)`).
		Suggest(`$dst = append($x, $a, $b)`)
}

func gocriticYodaStyleExpr(m fluent.Matcher) {
	m.Match(`nil != $_`,
		`0 != $_`).
		Report(`yoda-style expression`)
}

func gocriticUnderef(m fluent.Matcher) {
	m.Match(`(*$arr)[$i]`).
		Where(m["arr"].Type.Is(`*[$_]$_`)).
		Suggest(`$arr[$i]`)
}

func gocriticEmptyStringTest(m fluent.Matcher) {
	m.Match(`len($s) == 0`).
		Where(m["s"].Type.Is(`string`)).
		Suggest(`$s == ""`)
	m.Match(`len($s) != 0`).
		Where(m["s"].Type.Is(`string`)).
		Suggest(`$s != ""`)
}

func gocriticUnslice(m fluent.Matcher) {
	m.Match(`$s[:]`).Where(m["s"].Type.Is(`string`)).Suggest(`$s`)
	m.Match(`$s[:]`).Where(m["s"].Type.Is(`[]$_`)).Suggest(`$s`)
}

func gocriticSwitchTrue(m fluent.Matcher) {
	m.Match(`switch true {$*_}`).Report(`can omit true in switch`)
}

func gocriticSloppyLen(m fluent.Matcher) {
	m.Match(`len($_) >= 0`).Report(`$$ is always true`)
	m.Match(`len($_) < 0`).Report(`$$ is always false`)
	m.Match(`len($s) <= 0`).Suggest(`len($s) == 0`)
}

func gocriticNewDeref(m fluent.Matcher) {
	// TODO: add missing patterns.
	m.Match(`*new(bool)`).Suggest(`false`)
	m.Match(`*new(string)`).Suggest(`""`)
	m.Match(`*new(int)`).Suggest(`0`)
	m.Match(`*new(int32)`).Suggest(`int32(0)`)
	m.Match(`*new(float64)`).Suggest(`0.0`)
	m.Match(`*new(float32)`).Suggest(`float32(0)`)
}

func gocriticFlagDeref(m fluent.Matcher) {
	m.Match(`*flag.Bool($*_)`,
		`*flag.Float64($*_)`,
		`*flag.Duration($*_)`,
		`*flag.Int($*_)`,
		`*flag.Int64($*_)`,
		`*flag.String($*_)`,
		`*flag.Uint($*_)`,
		`*flag.Uint64($*_)`).
		Report(`immediate deref in $$ is most likely an error`)
}

func reviveBoolLiteralInExpr(m fluent.Matcher) {
	m.Match(`$x == true`,
		`$x != true`,
		`$x == false`,
		`$x != false`).
		Report(`omit bool literal in expression`)
}

func gosimpleS1003(m fluent.Matcher) {
	m.Match(`strings.Index($s1, $s2) != -1`).Suggest(`strings.Contains($s1, $s2)`)
	m.Match(`strings.Index($s1, $s2) == -1`).Suggest(`!strings.Contains($s1, $s2)`)
	m.Match(`strings.IndexAny($s1, $s2) != -1`).Suggest(`strings.ContainsAny($s1, $s2)`)
	m.Match(`strings.IndexAny($s1, $s2) == -1`).Suggest(`!strings.ContainsAny($s1, $s2)`)
}
