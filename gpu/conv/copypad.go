package conv

import (
	"github.com/barnex/cuda5/cu"
	"github.com/barnex/cuda5/safe"
	"nimble-cube/core"
	"nimble-cube/gpu"
	"nimble-cube/gpu/ptx"
	"unsafe"
)

var copyPadKern cu.Function

// Copies src into dst (which is larger or smaller), at offset position.
func copyPad(dst, src safe.Float32s, dstsize, srcsize, offset [3]int, stream cu.Stream) {
	//panic("need to loop for 3D, in kernel")
	core.Assert(dst.Len() == core.Prod(dstsize))
	core.Assert(src.Len() == core.Prod(srcsize))
	// TODO: either remove offset or check offset

	if copyPadKern == 0 {
		mod := cu.ModuleLoadData(ptx.COPYPAD)
		copyPadKern = mod.GetFunction("copypad")
	}

	dstptr := dst.Pointer()
	srcptr := src.Pointer()

	block := 16
	gridJ := gpu.DivUp(gpu.Min(dstsize[1], srcsize[1]), block)
	gridK := gpu.DivUp(gpu.Min(dstsize[2], srcsize[2]), block)
	shmem := 0
	args := []unsafe.Pointer{
		unsafe.Pointer(&dstptr),
		unsafe.Pointer(&dstsize[0]),
		unsafe.Pointer(&dstsize[1]),
		unsafe.Pointer(&dstsize[2]),
		unsafe.Pointer(&srcptr),
		unsafe.Pointer(&srcsize[0]),
		unsafe.Pointer(&srcsize[1]),
		unsafe.Pointer(&srcsize[2]),
		unsafe.Pointer(&offset[0]),
		unsafe.Pointer(&offset[1]),
		unsafe.Pointer(&offset[2])}

	// launch config is called x, y, z, but for us this is z, y, x
	cu.LaunchKernel(copyPadKern, gridK, gridJ, 1, block, block, 1, shmem, stream, args)
}
