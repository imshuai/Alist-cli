package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/imshuai/alistsdk-go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var (
	client *alistsdk.Client
	logger *logrus.Logger
	conf   *Config
)

const (
	VERSION = "v0.1.1"
)

type Config struct {
	Endpoint   string `json:"endpoint" yaml:"endpoint"`
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	SkipVerify bool   `json:"skip-verify" yaml:"skip-verify"`
	Proxy      string `json:"proxy" yaml:"proxy"`
}

func main() {
	logger = logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	app := cli.NewApp()
	app.Name = "alist-cli"
	app.Description = "alist command line interface"
	app.Version = VERSION
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{}
	app.HideVersion = true
	app.Commands = []*cli.Command{
		{
			Name:  "version",
			Usage: "show alist-cli version",
			Action: func(c *cli.Context) error {
				fmt.Printf("version: %s\n", VERSION)
				return nil
			},
		},
		{ // move
			Name:      "move",
			Aliases:   []string{"mv", "rename"},
			Usage:     "move file from src to dst, cannot cross mount point",
			UsageText: fmt.Sprintf(`%s move [-v] src-path dst-folder`, app.Name),
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "verbose",
					Aliases: []string{"v"},
					Usage:   "verbose mode, show more infomation",
					Value:   false,
					Action: func(c *cli.Context, v bool) error {
						if v {
							logger.SetLevel(logrus.InfoLevel)
							logger.Info("verbose mode enabled")
						}
						return nil
					},
				},
			},
			Action: func(c *cli.Context) error {
				//TODO: move file from src to dst
				err := ReadConfig()
				if err != nil {
					return fmt.Errorf("read config failed: %v", err)
				}
				client = alistsdk.NewClient(conf.Endpoint, conf.Username, conf.Password, conf.SkipVerify, 30)
				u, err := client.Login()
				if err != nil {
					return fmt.Errorf("login failed: %v", err)
				}
				logger.Infof("login success: %v", u.Username)
				src := c.Args().Get(0)
				dst := c.Args().Get(1)
				if !CheckPathIsExist(src) {
					return fmt.Errorf("path %s not exist", src)
				}
				if !CheckPathIsFolder(dst) {
					return fmt.Errorf("path %s is not a folder", dst)
				}
				//split src path to folder and filename
				pathParts := strings.Split(src, "/")
				folder := "/" + strings.Join(pathParts[:len(pathParts)-1], "/")
				filename := pathParts[len(pathParts)-1]
				err = client.Move(folder, dst, []string{filename})
				if err != nil {
					return fmt.Errorf("move %s to %s failed: %v", src, dst, err)
				}
				logger.Infof("move %s to %s success", src, dst)
				return nil
			},
		},
		{ // copy
			Name:      "copy",
			Aliases:   []string{"cp"},
			Usage:     "copy file from src to dst, can cross mount point",
			UsageText: fmt.Sprintf(`%s copy [-v] src-path dst-folder`, app.Name),
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "verbose",
					Aliases: []string{"v"},
					Usage:   "verbose mode, show more infomation",
					Value:   false,
					Action: func(c *cli.Context, v bool) error {
						if v {
							logger.SetLevel(logrus.InfoLevel)
						}
						return nil
					},
				},
			},
			Action: func(c *cli.Context) error {
				//TODO: copy file from src to dst
				err := ReadConfig()
				if err != nil {
					return fmt.Errorf("read config failed: %v", err)
				}
				client = alistsdk.NewClient(conf.Endpoint, conf.Username, conf.Password, conf.SkipVerify, 30)
				u, err := client.Login()
				if err != nil {
					return fmt.Errorf("login failed: %v", err)
				}
				logger.Infof("login success: %v", u.Username)
				src := c.Args().Get(0)
				dst := c.Args().Get(1)
				dst = strings.TrimSuffix(dst, "/")
				if !CheckPathIsExist(src) {
					return fmt.Errorf("path %s not exist", src)
				}
				pathParts := strings.Split(src, "/")
				folder := "/" + strings.Join(pathParts[:len(pathParts)-1], "/")
				filename := pathParts[len(pathParts)-1]
				if CheckPathIsExist(dst + "/" + filename) {
					return fmt.Errorf("dst path %s exist", dst+"/"+filename)
				}
				if !CheckPathIsFolder(dst) {
					return fmt.Errorf("path %s is not a folder", dst)
				}
				err = client.Copy(folder, dst, []string{filename})
				if err != nil {
					return fmt.Errorf("copy %s to %s failed: %v", src, dst, err)
				}
				logger.Infof("copy %s to %s mission submit success", src, dst)
				return nil
			},
		},
		{ // list
			Name:      "list",
			Aliases:   []string{"ls"},
			Usage:     "list all files in path",
			UsageText: fmt.Sprintf(`%s list [-v] path`, app.Name),
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "verbose",
					Aliases: []string{"v"},
					Usage:   "verbose mode, show more infomation",
					Value:   false,
					Action: func(c *cli.Context, v bool) error {
						if v {
							logger.SetLevel(logrus.InfoLevel)
						}
						return nil
					},
				},
			},
			Action: func(c *cli.Context) error {
				//TODO: list all files in path
				err := ReadConfig()
				if err != nil {
					return fmt.Errorf("read config failed: %v", err)
				}
				client = alistsdk.NewClient(conf.Endpoint, conf.Username, conf.Password, conf.SkipVerify, 30)
				u, err := client.Login()
				if err != nil {
					return fmt.Errorf("login failed: %v", err)
				}
				logger.Infof("login success: %v", u.Username)
				path := c.Args().Get(0)
				files, err := client.List(path, "", 0, 0, true)
				if err != nil {
					return fmt.Errorf("list %s failed: %v", path, err)
				}
				fmt.Printf("Files in %s:\n", path)
				for _, f := range files {
					fmt.Printf("%-30.30s\t%6.6s\n", f.Name, func() string {
						if f.IsDir {
							return "folder"
						}
						return "file"
					}())
				}
				logger.Infof("list %s success, get %d files", path, len(files))
				return nil
			},
		},
		{ // delete
			Name:      "delete",
			Aliases:   []string{"rm"},
			Usage:     "delete file from path",
			UsageText: fmt.Sprintf(`%s delete [-v] path`, app.Name),
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "verbose",
					Aliases: []string{"v"},
					Usage:   "verbose mode, show more infomation",
					Value:   false,
					Action: func(c *cli.Context, v bool) error {
						if v {
							logger.SetLevel(logrus.InfoLevel)
						}
						return nil
					},
				},
			},
			Action: func(c *cli.Context) error {
				//TODO: delete file from path
				err := ReadConfig()
				if err != nil {
					return fmt.Errorf("read config failed: %v", err)
				}
				client = alistsdk.NewClient(conf.Endpoint, conf.Username, conf.Password, conf.SkipVerify, 30)
				u, err := client.Login()
				if err != nil {
					return fmt.Errorf("login failed: %v", err)
				}
				logger.Infof("login success: %v", u.Username)
				path := c.Args().Get(0)
				if !CheckPathIsExist(path) {
					return fmt.Errorf("path %s not exist", path)
				}
				pathParts := strings.Split(path, "/")
				folder := "/" + strings.Join(pathParts[:len(pathParts)-1], "/")
				filename := pathParts[len(pathParts)-1]
				err = client.Remove(folder, []string{filename})
				if err != nil {
					return fmt.Errorf("delete %s failed: %v", path, err)
				}
				logger.Infof("delete %s success", path)
				return nil
			},
		},
		{ // upload
			Name:      "upload",
			Usage:     "upload file to path, imcomplete not supported yet",
			UsageText: fmt.Sprintf(`%s upload [-v] src-file dst-folder`, app.Name),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "path",
					Aliases: []string{"p"},
					Usage:   "path to upload",
				},
				&cli.BoolFlag{
					Name:    "verbose",
					Aliases: []string{"v"},
					Usage:   "verbose mode, show more infomation",
					Value:   false,
					Action: func(c *cli.Context, v bool) error {
						if v {
							logger.SetLevel(logrus.InfoLevel)
						}
						return nil
					},
				},
			},
			Action: func(c *cli.Context) error {
				//TODO: upload file to path
				return nil
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.Error(err)
	}
}

func ReadConfig() error {
	conf = &Config{}
	// read config file from config.json or config.yaml, prefer config.yaml
	if _, err := os.Stat("config.yaml"); err == nil {
		logger.Info("read config from config.yaml")
		byts, err := os.ReadFile("config.yaml")
		if err != nil {
			return err
		}
		return yaml.Unmarshal(byts, conf)
	} else if os.IsNotExist(err) {
		logger.Info("read config from config.json")
		byts, err := os.ReadFile("config.json")
		if err != nil {
			return err
		}
		return json.Unmarshal(byts, conf)
	} else {
		return err
	}
}

func CheckPathIsFolder(path string) bool {
	f, e := client.Get(path, "")
	if e != nil {
		logger.Errorf("get %s info failed with error: %v", path, e)
		return false
	}
	if !f.IsDir {
		logger.Errorf("%s is not a directory", path)
		return false
	}
	return true
}

func CheckPathIsExist(path string) bool {
	_, e := client.Get(path, "")
	if e != nil {
		logger.Errorf("get %s info failed with error: %v", path, e)
		return false
	}
	return true
}
