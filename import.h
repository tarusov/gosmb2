#include <errno.h>
#include <fcntl.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>

#include <smb2/smb2.h>
#include <smb2/libsmb2.h>
#include <smb2/libsmb2-raw.h>
#include <libsmb2-private.h>

typedef struct smb2_context *contextPtr;
typedef struct smb2fh       *fileHandlerPtr;
typedef struct smb2dir      *dirHandlerPtr;
typedef struct smb2_stat_64 fileInfo;
typedef struct smb2_statvfs vfsInfo;