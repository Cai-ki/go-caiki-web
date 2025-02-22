package ccache

import "ccache/ccachepb"

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	// Get(group string, key string) ([]byte, error)
	Get(in *ccachepb.Request, out *ccachepb.Response) error
}
