#include <stdio.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <poll.h>
#include <sys/eventfd.h>
#include <errno.h>

#include "ancillary.h"
#include "ssgo.h"


// X will be replaced by '\0' post-snprintf


int send_fd(int fd)
{
    int sock;
    struct sockaddr_un addr;

    if ((sock = socket(AF_UNIX, SOCK_STREAM, 0)) == -1) {
       
        return -2;
    }

    // Set timeout to 3s
    struct timeval tv;
    tv.tv_sec  = 3;
    tv.tv_usec = 0;
    setsockopt(sock, SOL_SOCKET, SO_RCVTIMEO, (char *)&tv, sizeof(struct timeval));
    setsockopt(sock, SOL_SOCKET, SO_SNDTIMEO, (char *)&tv, sizeof(struct timeval));

    memset(&addr, 0, sizeof(addr));
    addr.sun_family = AF_UNIX;
    strncpy(addr.sun_path, "protect_path", sizeof(addr.sun_path) - 1);

    if (connect(sock, (struct sockaddr *)&addr, sizeof(addr)) == -1) {
       
        close(sock);
        return -3;
    }

    if (ancil_send_fd(sock, fd)) {
     
        close(sock);
        return -4;
    }

    char ret = 0;

    if (recv(sock, &ret, 1, 0) == -1) {
      
        close(sock);
        return -5;
    }

    close(sock);
    return ret;
}
