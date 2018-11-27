This project forked from [alecthomas/log4go](http://github.com/alecthomas/log4go) and adapted with golang v1.5.3.

I hava change some features and done with some bugfixes. Include:

1. Reuse the exist log file (only if the size and maxline not reached) rather than create a new file every boot time.

2. Rotate log filename bug fixes with filelog

3. log4go_test.go now can run on go1.5.3

4. More accurate log time

5. Exclude specified log by adding the exclude pattern, which will match the *source* field. See &lt;exclude&gt; attribute in [example.xml](https://github.com/kimiazhu/log4go/blob/master/example.xml)

6. Auto load configuration in init() method, when the `log4go.xml` placed in exec dir or `{exec_dir}/conf` dir.

7. Add access log. example:

		<filter enabled="true">
			<tag>access</tag> <!-- the tag of accesslog MUST be access -->
			<type>file</type>
			<level>ACCESS</level> <!-- the level of accesslog MUST be access -->
			<property name="filename">log/access.log</property>
			<property name="format">[%D %T] [%L] %M</property>
			<property name="rotate">true</property>
			<property name="maxsize">100M</property>
			<property name="maxlines">100K</property>
			<property name="daily">true</property>
		</filter>

8. When you call log4go.Critical(), it will print the trace stack. I Also added a Recover() method to print stack only when there is a panic, otherwise it will print log as Error()

	`usage:`
	
		defer log4go.Recover("this is a msg: %v", "msg")
		// or
		defer log4go.Recover(func(err interface{}) string {
            // ... put your code here, construct the error message and return
            return fmt.Sprintf("recover..v1=%v;v2=%v;err=%v", 1, 2, err)
        })

9. Add Setup() to load config by passing a xml config string

### Installation:
- Run `go get github.com/kimiazhu/log4go`

### Usage:
- Add the following import:
import log "github.com/kimiazhu/log4go"

### TODO:

### Acknowledgements:
- pomack
  For providing awesome patches to bring log4go up to the latest Go spec
