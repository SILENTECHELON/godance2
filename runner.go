package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/stacktitan/smb/smb"
)

type Runner struct {
	counter     int
	running     bool
	conf        *Config
	startTime   time.Time
	currentUser string
	currentPass string
	currentHost string
	workingPass []byte
}

func NewRunner(conf *Config) Runner {
	var r Runner
	r.conf = conf
	r.running = false
	r.counter = 0
	r.currentUser = ""
	r.currentPass = ""
	r.currentHost = ""
	r.workingPass = nil
	return r
}

func (r *Runner) Start() {
	r.running = true
	defer r.Stop()
	fmt.Println(BANNER)
	fmt.Println(SEP)
	fmt.Printf(" [*] Number of usernames: %d\n", r.conf.users.Total())
	fmt.Printf(" [*] Number of passwords: %d\n", r.conf.passwds.Total())
	fmt.Printf(" [*] Test cases: %d\n", r.conf.users.Total()*r.conf.passwds.Total())
	fmt.Printf(" [*] Number of threads: %d\n", r.conf.threads)
	fmt.Println(SEP)

	var wg sync.WaitGroup
	var workingPass []byte
	totals := r.conf.users.Total() * r.conf.passwds.Total() * r.conf.host.Total()
	wg.Add(1)

	go r.runProgress(&wg, totals)

	result, ferr := os.Create("results.csv")

	if ferr != nil {
		return
	}
	defer result.Close()

	limiter := make(chan bool, r.conf.threads)
	for r.conf.passwds.Next() {
		wg.Add(1)
		nextPassword, passPos := r.conf.passwds.Value()

		if passPos < 0 {
			return
		}
		for r.conf.users.Next() {
			wg.Add(1)
			nextUser, userPos := r.conf.users.Value()
			r.counter += r.conf.host.Total()
			if userPos < 0 {
				return
			}
			for r.conf.host.Next() {
				limiter <- true
				wg.Add(1)
				theHost, thePos := r.conf.host.Value()
				if string(nextPassword) == "!!user!!" {
					workingPass = nextUser
				} else {
					workingPass = nextPassword
				}

				go func() {
					// release a slot in queue when exiting
					defer func() { <-limiter }()

					r.currentHost = string(theHost)
					r.currentUser = string(nextUser)
					r.currentPass = string(workingPass)
					w := bufio.NewWriterSize(result, 512)

					taskRes := r.RunTask(nextUser, workingPass, theHost, &wg, w)
					w.Flush()

					if taskRes != nil {
						r.conf.host.Remove(thePos)
						fmt.Printf("\n [%d] Success: %s // Username: %s // Password: %s\n", thePos, theHost, nextUser, workingPass)
					}

				}()

			}

			// Reset the pwd inputlist position
			wg.Done()
			r.conf.host.position = -1

		}

		wg.Done()
		r.conf.users.position = -1
	}

	wg.Wait()
}

func (r *Runner) runProgress(wg *sync.WaitGroup, total int) {
	defer wg.Done()
	r.startTime = time.Now()
	totalProgress := total
	for r.counter <= totalProgress {
		r.updateProgress()
		if r.counter == totalProgress {
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (r *Runner) updateProgress() {
	//TODO: refactor to use a defined progress struct for future output modules
	runningSecs := int((time.Now().Sub(r.startTime)) / time.Second)
	var reqRate int
	if runningSecs > 0 {
		reqRate = int(r.counter / runningSecs)
	} else {
		reqRate = 0
	}
	dur := time.Now().Sub(r.startTime)
	hours := dur / time.Hour
	dur -= hours * time.Hour
	mins := dur / time.Minute
	dur -= mins * time.Minute
	secs := dur / time.Second

	progString := fmt.Sprintf(":: Progress: [:: %6d tries/sec :: Duration: [%02d:%02d:%02d] :: [%15s:%15s]", int(reqRate), hours, mins, secs, r.currentUser, r.currentPass)
	fmt.Fprintf(os.Stderr, "%s%s", TERMINAL_CLEAR_LINE, progString)
}

func (r *Runner) RunTask(username []byte, password []byte, host []byte, work *sync.WaitGroup, txt *bufio.Writer) []byte {
	options := smb.Options{
		Host:     string(host),
		Port:     445,
		User:     string(username),
		Password: string(password),
		Domain:   r.conf.domain,
	}

	session, err := smb.NewSession(options, r.conf.debug)
	defer work.Done()
	//	result := fmt.Sprint("%s %s %s", username, password, host)
	//	os.WriteFile("dummyhash", []byte(result), fs.ModeAppend)
	if err != nil {

		err = nil
		return nil
	}

	if session.IsAuthenticated {

		fmt.Fprintln(txt, string(host), string(username), string(password), "\r")
		session.Close()

		return host
	} else {

		session.Close()

		return nil
	}
}

func (r *Runner) Stop() {
	fmt.Printf("\n")
	r.running = false
}
