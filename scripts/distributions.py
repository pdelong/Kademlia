import functools
import random
import bisect
import math

class Zipf:
    def __init__(self, arr, alpha):
        n = len(arr)
        tmp = [1. / math.pow(float(i), alpha) for i in range(1, n+1)]
        zeta = functools.reduce(lambda sums, x: sums + [sums[-1] + x], tmp, [0])

        self.arr = arr
        self.distMap = [x / zeta[-1] for x in zeta]

    def next(self):
        u = random.random()
        index = bisect.bisect(self.distMap, u)

        return self.arr[index-1]

class Uniform:
    def __init__(self, arr):
        self.arr = arr

    def next(self):
        return random.choice(self.arr)

class Linear:
    def __init__(self, arr, m):
        n = len(arr)
        distMap = [1] * len(n)
        distMap[-1] = 1/m
        change = (distMap[0] - distMap[-1]) / m

        distMap = [(x - change) / distMap[-1] for x in distMap]

        self.arr = arr
        self.distMap = distMap

    def next(self):
        u = random.random()
        index = bisect.bisect(self.distMap, u)

        return self.arr[index-1]
