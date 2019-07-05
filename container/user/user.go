package user

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const (
	unixPasswdPath = "/etc/passwd"
)

func SetUser(user string) error {
	userName, _, err := getUserGroupName(user)
	if err != nil {
		return err
	}

	uid, gid, err := lookupUidGid(userName)
	if err != nil {
		return err
	}

	if err := setgid(uid); err != nil {
		return err
	}
	if err := setuid(gid); err != nil {
		return err
	}

	return nil
}

func getUserGroupName(user string) (string, string, error) {
	//check input user format
	userGroup := strings.Split(user, ":")
	if len(userGroup) != 2 {
		return "", "", fmt.Errorf("user format should be username:groupname")
	}

	return userGroup[0], userGroup[1], nil
}

func lookupUidGid(userName string) (int, int, error) {
	passwd, err := os.Open(unixPasswdPath)
	if err != nil {
		return 0, 0, err
	}
	defer passwd.Close()

	s := bufio.NewScanner(passwd)
	for s.Scan() {
		if err := s.Err(); err != nil {
			return 0, 0, err
		}

		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}

		// see: man 5 passwd
		//  name:password:UID:GID:GECOS:directory:shell
		// Name:Pass:Uid:Gid:Gecos:Home:Shell
		//  root:x:0:0:root:/root:/bin/bash
		//  adm:x:3:4:adm:/var/adm:/bin/false
		var name, pass, gecos, home, shell string
		var uid, gid int
		parseLine(line, &name, &pass, &uid, &gid, &gecos, &home, &shell)

		if name == userName {
			return uid, gid, nil
		}
	}
	return 0, 0, nil
}

// setuid sets the uid of the calling thread to the specified uid.
func setuid(uid int) (err error) {
	_, _, e1 := syscall.RawSyscall(syscall.SYS_SETUID, uintptr(uid), 0, 0)
	if e1 != 0 {
		err = e1
	}
	return nil
}

// setgid sets the gid of the calling thread to the specified gid.
func setgid(gid int) (err error) {
	_, _, e1 := syscall.RawSyscall(syscall.SYS_SETGID, uintptr(gid), 0, 0)
	if e1 != 0 {
		err = e1
	}
	return nil
}

func parseLine(line string, v ...interface{}) {
	parseParts(strings.Split(line, ":"), v...)
}

func parseParts(parts []string, v ...interface{}) {
	if len(parts) == 0 {
		return
	}

	for i, p := range parts {
		// Ignore cases where we don't have enough fields to populate the arguments.
		// Some configuration files like to misbehave.
		if len(v) <= i {
			break
		}

		// Use the type of the argument to figure out how to parse it, scanf() style.
		// This is legit.
		switch e := v[i].(type) {
		case *string:
			*e = p
		case *int:
			// "numbers", with conversion errors ignored because of some misbehaving configuration files.
			*e, _ = strconv.Atoi(p)
		case *int64:
			*e, _ = strconv.ParseInt(p, 10, 64)
		case *[]string:
			// Comma-separated lists.
			if p != "" {
				*e = strings.Split(p, ",")
			} else {
				*e = []string{}
			}
		default:
			// Someone goof'd when writing code using this function. Scream so they can hear us.
			panic(fmt.Sprintf("parseLine only accepts {*string, *int, *int64, *[]string} as arguments! %#v is not a pointer!", e))
		}
	}
}
