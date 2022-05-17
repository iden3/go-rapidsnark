// This file uses a preprocessor macro defined by the various build_*.go
// files to determine whether to import the bundled librapidsnark header, or
// the system one.
// This is needed because cgo will automatically add -I. to the include
// path, so <prover.h> would find a bundled header instead of
// the system one

#ifdef USE_VENDORED_RAPIDSNARK
#include "rapidsnark_vendor/prover.h"
#else
#include <prover.h>
#endif