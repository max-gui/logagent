package logagent

import (
	"fmt"
	"io"
	"log"

	"github.com/sirupsen/logrus"
)

// NullFormatter formats logs into text
type agentHook struct {
	fireChannel     chan *logrus.Entry
	asyncBufferSize int
	formatter       logrus.Formatter
	writer          io.Writer
	// loghook         logrus.Hook
}

// Levels returns configured log levels.
func (hook *agentHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// func NewHook(buffersize int, realhook logrus.Hook, logwritter io.Writer, logformatter logrus.Formatter) *agentHook {
func NewHook(buffersize int, logwritter io.Writer, logformatter logrus.Formatter) *agentHook {
	hook := &agentHook{
		asyncBufferSize: buffersize,
		// loghook:         realhook,
		writer:    logwritter,
		formatter: logformatter,
		// fireChannel:     make(chan *logrus.Entry, asyncBufferSize),
	}
	hook.makeAsync()
	// hook.fireChannel = make(chan *logrus.Entry, hook.asyncBufferSize)
	return hook
}

func (f *agentHook) makeAsync() {

	f.fireChannel = make(chan *logrus.Entry, f.asyncBufferSize)
	fmt.Printf("file hook will use a async buffer with size %d\n", f.asyncBufferSize)
	go func() {
		for entry := range f.fireChannel {
			// time.Sleep(time.Duration(200) * time.Millisecond)
			if err := f.send(entry); err != nil {
				fmt.Println("Error during sending message to file:", err)
			}
		}
	}()
}

// Fire is called when a log event is fired.
func (f *agentHook) Fire(entry *logrus.Entry) error {
	if f.fireChannel != nil { // Async mode.
		select {
		case f.fireChannel <- entry: // try and put into chan, if fail will to default
		default:
			// if f.asyncBlock {
			// 	fmt.Println("the log buffered chan is full! will block")
			// 	f.fireChannel <- entry // Blocks the goroutine because buffer is full.
			// 	return nil
			// }
			// fmt.Println("the log buffered chan is full! will drop")
			// Drop message by default.
		}
		return nil
	}

	// Sync mode.
	return nil
	// return f.send(entry)
}

//if has loghook
// func (f *agentHook) send(entry *logrus.Entry) error {

// 	return f.loghook.Fire(entry)
// }
// type mutexKV struct {
// 	sync.RWMutex
// 	msgs    [][]byte
// 	lastest time.Time
// }

// var logs = mutexKV{msgs: make([][]byte, 0), lastest: time.Now()}

// func (v *mutexKV) help(tricky func([][]byte) (bool, interface{})) (bool, interface{}) {
// 	v.Lock()
// 	ok, res := tricky(v.msgs)
// 	v.Unlock()
// 	return ok, res
// }

func (f *agentHook) send(entry *logrus.Entry) error {

	msg, err := f.formatter.Format(entry)

	if err != nil {
		log.Println("failed to generate string for entry:", err)
		return err
	}
	// logs.Lock()
	// logs.msgs = append(logs.msgs, msg)
	// if len(logs.msgs) > 100 || time.Duration(100)*time.Millisecond > time.Since(logs.lastest) {
	// 	for _, logmsg := range logs.msgs {

	// 		_, err = f.writer.Write(logmsg)
	// 		// return err
	// 	}
	// 	logs.msgs = make([][]byte, 0)
	// 	logs.lastest = time.Now()
	// }
	// logs.Unlock()
	_, err = f.writer.Write(msg)
	return err

	// return f.loghook.Fire(entry)
}
