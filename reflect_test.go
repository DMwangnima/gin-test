package test

import "testing"

func TestJointUrl(t *testing.T) {
	basicUrl := "hey"
	test1Struct := []struct {
		Name    string `uri:"name"`
		Address string `uri:"address"`
	}{
		{Name: "hello", Address: "china"},
		{Name: "", Address: "America"},
		{Name: "Jordan", Address: ""},
		{Name: "", Address: ""},
	}

	expect1 := []string {
		basicUrl + "?name=hello&address=china",
		basicUrl + "?address=America",
		basicUrl + "?name=Jordan",
		basicUrl,
	}

	test2Struct := []struct {
		Age1 int   `uri:"age1"`
		Age2 int8  `uri:"age2"`
		Age3 int16 `uri:"age3"`
		Age4 int32 `uri:"age4"`
		Age5 int64 `uri:"age5"`
	}{
		{Age1: 12, Age2: 12, Age3: 12, Age4: 12, Age5: 12},
		{Age1: 0, Age2: 0, Age3: 0, Age4: 0, Age5: 0},
	}

	expect2 := []string {
		basicUrl + "?age1=12&age2=12&age3=12&age4=12&age5=12",
		basicUrl,
	}

	test3Struct := []struct {
		Name string `uri:"name"`
		Age  int    `uri:"age"`
	}{
		{Name: "hello", Age: 5},
		{Name: "", Age: 5},
		{Name: "hello", Age: 0},
		{Name: "", Age: 0},
	}
	test3StructPtr := []*struct{
		Name string `uri:"name"`
		Age  int    `uri:"age"`
	}{
		&test3Struct[0],
		&test3Struct[1],
		&test3Struct[2],
		&test3Struct[3],
	}

	expect3 := []string{
		basicUrl + "?name=hello&age=5",
		basicUrl + "?age=5",
		basicUrl + "?name=hello",
		basicUrl,
	}

	for i, ts := range test1Struct {
		url := JointUrl(basicUrl, ts)
		if url != expect1[i] {
			t.Errorf("JointUrl failed, expect %s, get %s", expect1[i], url)
		}
	}

	for i, ts := range test2Struct {
		url := JointUrl(basicUrl, ts)
		if url != expect2[i] {
			t.Errorf("JointUrl failed, expect %s, get %s", expect2[i], url)
		}
	}

	for i, tsp := range test3StructPtr {
		url := JointUrl(basicUrl, tsp)
		if url != expect3[i] {
			t.Errorf("JointUrl failed, expect %s, get %s", expect3[i], url)
		}
	}
}
