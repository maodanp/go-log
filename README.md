# go-log
It is my log source for other project, it's easy to use in your go project.

## Usage
~~~go
func test() {
	log := NewLogger(os.Stdout, Config{highlighting: true, DispFuncCall: true})
	log.SetLogLevel(LOG_DEBUG)
	log.Info("hello go-log")
}
~~~
## Features
* log can be written in different ways, which implements io.Writer(file, console, connection etc.)
* You can disply log with different colors which only can be effective to stdout
* log can be rotated automatically based on daily

## License
kingshard is under the Apache 2.0 license.

