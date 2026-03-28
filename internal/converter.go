package internal

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/samber/lo"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"go.uber.org/zap"
)

type Convert struct {
	in  string
	out string

	logger *zap.SugaredLogger
}

func NewConverter(in, out string) Converter {
	log, _ := zap.NewProduction()
	logger := log.Sugar()

	return &Convert{
		in:     in,
		out:    out,
		logger: logger,
	}
}

func (c Convert) Process() {
	list, err := c.fileList(c.in)
	if err != nil {
		c.logger.Fatalf("fileList: %s", err.Error())
		return
	}

	c.logger.Info(len(list))

	c.transform(list, c.in, c.out)
}

func (c Convert) transform(list []string, in, out string) {
	for i, v := range list {
		c.logger.Infof("%d of %d. %s -> %s", i+1, len(list), fmt.Sprintf("%s/%s", strings.TrimRight(in, "/"), v), fmt.Sprintf("%s/%s.mp3", strings.TrimRight(out, "/"), v[:strings.LastIndex(v, ".")]))
		time.Sleep(time.Second * 2)

		if err := ffmpeg.Input(fmt.Sprintf("%s/%s", strings.TrimRight(in, "/"), v)).
			Output(fmt.Sprintf("%s/%s.mp3", strings.TrimRight(out, "/"), v[:strings.LastIndex(v, ".")])).
			OverWriteOutput().ErrorToStdOut().Run(); err != nil {

			c.logger.Errorf("FFMPEG: %s", err.Error())
			continue
		}

		c.logger.Infof("FILE: %s CONVERTED\n", fmt.Sprintf("%s/%s", strings.TrimRight(in, "/"), v))
		c.logger.Infof("OUT: %s\n", fmt.Sprintf("%s/%s.mp3", strings.TrimRight(out, "/"), v[:strings.LastIndex(v, ".")]))

		if err := c.remove(in, v); err != nil {
			c.logger.Warnf("REMOVE FAILED: %s\n", err.Error())
			continue
		}

		c.logger.Infof("REMOVED FROM %s\n\n", in)
	}

}

func (c Convert) remove(in, name string) error {
	return os.Remove(fmt.Sprintf("%s/%s", strings.TrimRight(in, "/"), name))
}

func (c Convert) fileList(indir string) ([]string, error) {
	dir, err := os.ReadDir(indir)
	if err != nil {
		return nil, fmt.Errorf("ReadDir: %w", err)
	}

	return lo.Map(dir, func(item os.DirEntry, _ int) string {
		return item.Name()
	}), nil
}
