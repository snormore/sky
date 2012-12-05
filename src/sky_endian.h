#ifndef _sky_endian_h
#define _sky_endian_h

#include <inttypes.h>


//==============================================================================
//
// Byte Order
//
//==============================================================================

#include <sys/types.h>

#ifndef BYTE_ORDER
#if defined(linux) || defined(__linux__)
# include <endian.h>
#else
# include <machine/endian.h>
#endif
#endif

#if !defined(BYTE_ORDER) && !defined(__BYTE_ORDER)
#error "Undefined byte order"
#endif


//==============================================================================
//
// Functions
//
//==============================================================================

uint16_t bswap16(uint16_t value);
uint32_t bswap32(uint32_t value);
uint64_t __bswap64(uint64_t value);

#if (BYTE_ORDER == LITTLE_ENDIAN) || (__BYTE_ORDER == __LITTLE_ENDIAN)
#define htonll(x) __bswap64(x)
#define ntohll(x) __bswap64(x)
#else
#define htonll(x) x
#define ntohll(x) x
#endif

#endif
