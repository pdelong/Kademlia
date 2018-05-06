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

if __name__ == '__main__':
    n = []
    for _ in range(1, 100):
        n = n + [random.random()]

    a = 3
    dist = Zipf(n, a)

    for _ in range(1, 1000):
        print(dist.next())
