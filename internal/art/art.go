package art

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)


const KeySize = 32

// TODO: these probably shouldn't be global'
var KekMutex *sync.Mutex
var Kek [KeySize]byte

// TODO: this is a hacky skeleton function; real ART code goes here that would listen to pubsub
func ArtMainLoop(kekPath string) {
	fmt.Println("hello from art thread")
	var kek1, kek2 [KeySize]byte

	for {
		// TODO: event based instead of crappy polling
		//fmt.Println("art thread sleeping")
		time.Sleep(2 * time.Second)

		// lock and grab copy of current KEK
		KekMutex.Lock()
		copy(kek1[:], Kek[:])
		KekMutex.Unlock()

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
		KekMutex.Lock()
		copy(Kek[:], kek2[:])
		KekMutex.Unlock()
	}
}
