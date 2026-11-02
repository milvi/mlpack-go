//go:build linux && mlpack_memtest

package mlpack

/*
#include <malloc.h>
#include <stdint.h>
static int64_t mlpack_native_bytes(void) {
  struct mallinfo2 m = mallinfo2();
  return (int64_t)m.uordblks + (int64_t)m.hblkhd;
}
static void mlpack_trim_native(void) { malloc_trim(0); }
*/
import "C"

// Ces mesures sont compilées uniquement pour les tests explicites sous glibc.
func nativeBytes() int64 { C.mlpack_trim_native(); return int64(C.mlpack_native_bytes()) }
