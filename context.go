package smb

/*
#cgo CFLAGS:  -I./include
#cgo amd64   LDFLAGS: -L./lib/amd64 -lsmb2 -lgssapi_krb5 -lkrb5 -lk5crypto -lkrb5support -lcom_err -ldl -lresolv -lpthread
#cgo 386     LDFLAGS: -L./lib/386   -lsmb2 -lgssapi_krb5 -lkrb5 -lk5crypto -lkrb5support -lcom_err -ldl -lresolv -lpthread
#include "import.h"
*/
import "C"
import (
	"errors"
	"fmt"
)

// context is handler for smb2_context type in Go.
//
// struct smb2_context {
// 	t_socket fd;
// 	t_socket *connecting_fds;
// 	size_t connecting_fds_count;
// 	struct addrinfo *addrinfos;
// 	const struct addrinfo *next_addrinfo;
// 	int timeout;
// 	enum smb2_sec sec;
// 	uint16_t security_mode;
// 	int use_cached_creds;
// 	enum smb2_negotiate_version version;
// 	const char *server;
// 	const char *share;
// 	const char *user;
// 	/* -- Only used with --without-libkrb5 -- */
// 	const char *password;
// 	const char *domain;
// 	const char *workstation;
// 	char client_challenge[8];
// 	smb2_command_cb connect_cb;
// 	void *connect_data;
// 	int credits;
// 	char client_guid[16];
// 	uint32_t tree_id;
// 	uint64_t message_id;
// 	uint64_t session_id;
// 	uint8_t *session_key;
// 	uint8_t session_key_size;
// 	uint8_t seal:1;
// 	uint8_t sign:1;
// 	uint8_t signing_key[SMB2_KEY_SIZE];
// 	uint8_t serverin_key[SMB2_KEY_SIZE];
// 	uint8_t serverout_key[SMB2_KEY_SIZE];
// 	uint8_t salt[SMB2_SALT_SIZE];
// 	uint16_t cypher;
// 	uint8_t preauthhash[SMB2_PREAUTH_HASH_SIZE];
// 	/* -- For handling received smb3 encrypted blobs -- */
// 	unsigned char *enc;
// 	size_t enc_len;
// 	int enc_pos;
// 	/* -- For sending PDUs -- */
// struct smb2_pdu *outqueue;
// struct smb2_pdu *waitqueue;
// 	/* -- For receiving PDUs -- */
// 	struct smb2_io_vectors in;
// 	enum smb2_recv_state recv_state;
// 	/* SPL for the (compound) command we are currently reading */
// 	uint32_t spl;
// 	/* buffer to avoid having to malloc the header */
// 	uint8_t header[SMB2_HEADER_SIZE];
// 	struct smb2_header hdr;
// 	/* Offset into smb2->in where the payload for the current PDU starts */
// 	size_t payload_offset;
// 	/* Pointer to the current PDU that we are receiving the reply for.
// 	 * Only valid once the full smb2 header has been received.
// 	 */
// 	struct smb2_pdu *pdu;
// 	/* Server capabilities */
// 	uint8_t supports_multi_credit;
// 	uint32_t max_transact_size;
// 	uint32_t max_read_size;
// 	uint32_t max_write_size;
// 	uint16_t dialect;
// 	char error_string[MAX_ERROR_SIZE];
// 	/* Open filehandles */
// 	struct smb2fh *fhs;
// 	/* Open dirhandles */
// 	struct smb2dir *dirs;
// 	/* callbacks for the eventsystem */
// 	int events;
// 	smb2_change_fd_cb change_fd;
// 	smb2_change_events_cb change_events;
// 	/* dcerpc settings */
// 	int ndr;
// 	int endianess;
// };
//
type context struct {
	ptr C.contextPtr
}

// mkContext creates new smb2_context instance.
func mkContext() (*context, error) {

	fmt.Println("!!! mkContext")

	result := C.smb2_init_context()
	if result == nil {
		return nil, errors.New("failed to init smb2 context")
	}

	return &context{
		ptr: result,
	}, nil
}

// ok check ptr state for instance.
func (c *context) ok() bool {
	return c != nil && c.ptr != nil
}

// free current smb2_context.
func (c *context) free() {
	if c.ok() {
		C.smb2_destroy_context(c.ptr)
	}
	c.ptr = nil
}
