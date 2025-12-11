#ifndef __JHASH_H__
#define __JHASH_H__

#include <vmlinux.h>
//#include <linux/types.h>

#define JHASH_INITVAL 0xdeadbeef

static __always_inline __u32 jhash_2words(__u32 a, __u32 b, __u32 initval)
{
    __u32 c = initval, len = 8;
    __u32 x = a, y = b, z = JHASH_INITVAL + len + initval;

    x += y; y += z; z += x;
    x -= y; x ^= (y << 4) | (y >> 28); y += z;
    y -= z; y ^= (z << 6) | (z >> 26); z += x;
    z -= x; z ^= (x << 8) | (x >> 24); x += y;
    x -= y; x ^= (y << 16) | (y >> 16); y += z;
    y -= z; y ^= (z << 19) | (z >> 13); z += x;
    z -= x; z ^= (x << 4) | (x >> 28);

    return z;
}

#endif /* __JHASH_H__ */

