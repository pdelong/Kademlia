import functools
import random
import bisect
import math

class Zipf:
    def __init__(self, n, alpha):
        tmp = [1. / math.pow(float(i), alpha) for i in range(1, n+1)]
        zeta = functools.reduce(lambda sums, x: sums + [sums[-1] + x], tmp, [0])

        self.distMap = [x / zeta[-1] for x in zeta]

    def next(self):
        u = random.random()

        return bisect.bisect(self.distMap, u)
