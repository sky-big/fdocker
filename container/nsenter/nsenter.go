package nsenter

/*
#define _GNU_SOURCE
#include <unistd.h>
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>

__attribute__((constructor)) void enter_namespace(void) {
	char *fdocker_pid;
	fdocker_pid = getenv("fdocker_pid");
	if (fdocker_pid) {
		//fprintf(stdout, "got fdocker_pid=%s\n", fdocker_pid);
	} else {
		//fprintf(stdout, "missing fdocker_pid env skip nsenter");
		return;
	}
	char *fdocker_cmd;
	fdocker_cmd = getenv("fdocker_cmd");
	if (fdocker_cmd) {
		//fprintf(stdout, "got fdocker_cmd=%s\n", fdocker_cmd);
	} else {
		//fprintf(stdout, "missing fdocker_cmd env skip nsenter");
		return;
	}
	int i;
	char nspath[1024];
	char *namespaces[] = { "ipc", "uts", "net", "pid", "mnt" };

	for (i=0; i<5; i++) {
		sprintf(nspath, "/proc/%s/ns/%s", fdocker_pid, namespaces[i]);
		int fd = open(nspath, O_RDONLY);

		if (setns(fd, 0) == -1) {
			//fprintf(stderr, "setns on %s namespace failed: %s\n", namespaces[i], strerror(errno));
		} else {
			//fprintf(stdout, "setns on %s namespace succeeded\n", namespaces[i]);
		}
		close(fd);
	}
	int res = system(fdocker_cmd);
	exit(0);
	return;
}
*/
import "C"
