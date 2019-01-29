package main

import (
	"fmt"

	. "github.com/ArtemT/books2vk"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("books2vk")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/book2vk")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file is not loaded")
	}
	viper.SetEnvPrefix("books2vk")

	pflag.String("file", "", "XLSX file")
	pflag.Parse()
	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		println(err)
	}
}

func main() {
	f := OpenFile()
	defer func() {
		f.Save()
		f.Close()
	}()

	vk := NewService("")

	in := f.Proceed()
	out := vk.Send(in)
	done := f.Update(out)

	<-done
}
