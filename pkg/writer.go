package pkg

import (
	"fmt"
	"github.com/secsy/goftp"
	"os"
)

func Write(c *goftp.Client, list []string, path, output string) {
	for _, item := range list {
		out, err := os.Create(output + item)
		if err != nil {
			fmt.Printf("err: %s", err)
		}
		_ = c.Retrieve(path+item, out)
	}

}
