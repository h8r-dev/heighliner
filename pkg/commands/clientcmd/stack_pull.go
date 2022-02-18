package clientcmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var (
	stackPullCmd = &cobra.Command{
		Use:   "pull [stack url]",
		Short: "Pull a stack",
		Long:  "",
		RunE:  pullStack,
	}
	hlnStore string
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	hlnStore = userHomeDir + "/.hln"
	os.MkdirAll(hlnStore, 0777)
	handleErr(err)
}

func pullStack(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		err := errors.New("please specify stack url")
		return err
	}

	wg := sync.WaitGroup{}
	for i, v := range args {
		wg.Add(1)
		go downloadStack(i, v, &wg)
	}
	wg.Wait()

	return nil
}

func downloadStack(i int, url string, wg *sync.WaitGroup) {
	file, err := os.Create(hlnStore + "/stack" + fmt.Sprintf("%d", i) + ".tar.gz")
	handleErr(err)
	defer file.Close()

	rsp, err := http.Get(url)
	handleErr(err)
	defer rsp.Body.Close()

	io.Copy(file, rsp.Body)
	decompressStack(i, wg)
}

func decompressStack(i int, wg *sync.WaitGroup) {
	command := exec.Command("tar",
		"-zxvf", fmt.Sprintf("%s/stack%d.tar.gz", hlnStore, i),
		"-C", hlnStore)
	command.Stdout = &bytes.Buffer{}
	command.Stderr = &bytes.Buffer{}

	err := command.Run()
	if err != nil {
		fmt.Println(command.Stderr.(*bytes.Buffer).String())
	}

	fmt.Println(command.Stdout.(*bytes.Buffer).String())

	err = os.Remove(hlnStore + "/stack" + fmt.Sprintf("%d", i) + ".tar.gz")
	handleErr(err)

	wg.Done()
}

func handleErr(err error) {
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
