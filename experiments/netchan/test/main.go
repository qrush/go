package main

import "fmt"
import "netchan"

type value struct {
	i int
	s string
}

const count = 10

func exportSend(exp *netchan.Exporter) {
	ch := make(chan value)
	exp.Export("exportedSend", ch, netchan.Send, new(value))
	for i := 0; i < count; i++ {
		ch <- value{23 + i, "hello"}
	}
}

func importReceive(imp *netchan.Importer) {
	ch := make(chan value)
	imp.ImportNValues("exportedSend", ch, netchan.Recv, new(value), count)
	for i := 0; i < count; i++ {
		v := <-ch
		fmt.Printf("%v\n", v)
		if v.i != 23+i || v.s != "hello" {
			fmt.Printf("importReceive: bad value: expected %d, hello; got %+v", 23+i, v)
		}
	}
}

func main() {
	exp, _ := netchan.NewExporter("tcp", ":9292")
	imp, _ := netchan.NewImporter("tcp", ":9292")
	go exportSend(exp)
	importReceive(imp)
}
