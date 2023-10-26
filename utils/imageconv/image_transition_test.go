package imageconv_test

import (
	"testing"

	"github.com/erbieio/web2-bridge/utils/imageconv"
	"github.com/noelyahan/impexp"
	"github.com/noelyahan/mergi"
)

func TestInkEffect(t *testing.T) {
	anim, err := imageconv.InkEffect("https://www.erbiescan.io/ipfs/QmS7Pm4CAhU64qJB8Dh9mDq56vZpcCKfYLFAHCy87btGun")
	if err != nil {
		t.Error(err)
		return
	}
	mergi.Export(impexp.NewAnimationExporter(anim, "out.gif"))
}
