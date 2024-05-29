package art

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/googlecloudplatform/gcsfuse/internal/fs"
)

// TODO: this is a hacky skeleton function; real ART code goes here that would listen to pubsub
func ArtMainLoop(kek *[fs.KeySize]byte, kekMutex *sync.Mutex, kekPath string) {
	fmt.Println("hello from art thread")
	var kek1, kek2 [fs.KeySize]byte

	for {
		// TODO: event based instead of crappy polling
		//fmt.Println("art thread sleeping")
		time.Sleep(2 * time.Second)

		// lock and grab copy of current KEK
		kekMutex.Lock()
		copy(kek1[:], kek[:])
		kekMutex.Unlock()

		kekFile, err := os.Open(kekPath)
		if err != nil {
			fmt.Printf("art thread failed to open kek file: %v\n", err)
			continue
		}

		n, err := kekFile.Read(kek2[:])
		if err != nil {
			fmt.Printf("art thread failed to read kek file: %v\n", err)
			continue
		}
		if n != len(kek2) {
			fmt.Printf("art thread failed to read kek file: %d != %d\n", n, len(kek2))
			continue
		}
//fmt.Printf("cmp w/ k1 %d %x k2 %d %x\n", len(kek1), kek1, len(kek2), kek2)
		if bytes.Equal(kek1[:], kek2[:]) {
			continue
		}

		println("copying new key")
		kekMutex.Lock()
		copy(kek[:], kek2[:])
		kekMutex.Unlock()
	}
}
