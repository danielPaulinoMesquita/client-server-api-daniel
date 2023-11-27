package file

import (
	"os"
)

func WriterFile(dollar string) {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	_, err = f.WriteString("Dólar: " + dollar)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}
}
