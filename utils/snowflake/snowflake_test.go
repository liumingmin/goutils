package snowflake

import "testing"

func TestSnowflake(t *testing.T) {
	n, err := NewNode(-1)
	if err == nil {
		t.FailNow()
	}

	n, err = NewNode(1024)
	if err == nil {
		t.FailNow()
	}

	n, err = NewNode(2)
	if err != nil {
		t.Fatal(err)
	}

	id1 := n.Generate()
	id2 := n.Generate()
	if id1 == id2 {
		t.FailNow()
	}

	if ParseInt64(id1.Int64()) != id1 {
		t.FailNow()
	}

	idtemp, err := ParseString(id1.String())
	if err != nil {
		t.Fatal(err)
	}
	if idtemp != id1 {
		t.FailNow()
	}

	idtemp, err = ParseBase2(id1.Base2())
	if err != nil {
		t.Fatal(err)
	}
	if idtemp != id1 {
		t.FailNow()
	}

	idtemp, err = ParseBase32([]byte(id1.Base32()))
	if err != nil {
		t.Fatal(err)
	}
	if idtemp != id1 {
		t.FailNow()
	}

	idtemp, err = ParseBase36(id1.Base36())
	if err != nil {
		t.Fatal(err)
	}
	if idtemp != id1 {
		t.FailNow()
	}

	idtemp, err = ParseBase58([]byte(id1.Base58()))
	if err != nil {
		t.Fatal(err)
	}
	if idtemp != id1 {
		t.FailNow()
	}

	idtemp, err = ParseBase64(id1.Base64())
	if err != nil {
		t.Fatal(err)
	}
	if idtemp != id1 {
		t.FailNow()
	}

	idtemp, err = ParseBytes(id1.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if idtemp != id1 {
		t.FailNow()
	}

	idtemp = ParseIntBytes(id1.IntBytes())
	if idtemp != id1 {
		t.FailNow()
	}

	bs, err := id1.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	idtemp = ID(0)
	err = idtemp.UnmarshalJSON(bs)
	if err != nil {
		t.Fatal(err)
	}
	if idtemp != id1 {
		t.FailNow()
	}
}
