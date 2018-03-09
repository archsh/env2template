package main

import (
    "os"
    "fmt"
    "strings"
    "flag"
    "html/template"
    "io"
)

func err_println(s... interface{}) {
    os.Stderr.Write([]byte(fmt.Sprintln(s...)))
}

func err_printf(f string, s...interface{}) {
    os.Stderr.Write([]byte(fmt.Sprintf(f, s...)))
}

const VERSION = "develop"

func main() {
    var show_version bool
    flag.BoolVar(&show_version, "v", false, "Display version")
    var template_filename, output_filename string
    //Template filename
    flag.StringVar(&template_filename, "t", "", "Template filename.")
    //Output filename
    flag.StringVar(&output_filename, "o", "", "Output filename. Empty means output to stdout.")
    flag.Parse()
    if show_version {
        err_println(VERSION)
        os.Exit(0)
    }
    if template_filename == "" {
        err_println("Template filename required!")
        os.Exit(1)
    }
    tmpl, e := template.ParseFiles(template_filename)
    if nil != e {
        err_printf("Load template file '%s' failed! <%s>\n", template_filename, e)
        os.Exit(1)
    }
    var env_map = make(map[string]string)
    for _, x := range os.Environ() {
        //err_println(">", x)
        pair := strings.SplitN(x, "=", 2)
        env_map[pair[0]] = pair[1]
    }
    var output io.Writer
    if output_filename == "" {
        output = os.Stdout
    }else if f, e := os.OpenFile(output_filename, os.O_WRONLY|os.O_CREATE, 0666); nil != e {
        err_printf("Open or create output file '%s' failed: %s \n", output_filename, e)
    }else{
        output = f
        defer f.Close()
    }
    err := tmpl.Execute(output, env_map)
    if nil != err {
        err_println("Execute template failed:>", err)
    }
}