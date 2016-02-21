#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Single image resize perf-test.

"""
import argparse
import http.client
import io
import logging
import statistics
import sys
import time
from PIL import Image

PLT = None
NP = None
try:
    import matplotlib.pyplot
    import numpy
    PLT = matplotlib.pyplot
    NP = numpy
except ImportError:
    print('Matplotlib or Numpy not found, skipping plotting functionality')

logging.basicConfig(
    stream=sys.stdout,
    level=logging.INFO
)
LOGGER = logging.getLogger(__name__)

def get_opts(args):
    parser = argparse.ArgumentParser()
    parser.add_argument('-d', '--debug', help='Print debug information',
                        default=False, type=bool)
    parser.add_argument('--host', help='The target hostname to query',
                        default='localhost')
    parser.add_argument('-p', '--port', help='The target host port to query',
                        default=8005, type=int)
    parser.add_argument('-i', '--image', help='The image id to query for',
                        default='01.jpg')
    parser.add_argument('-s', '--step', help='Resize in this large steps',
                        default=8, type=int)
    parser.add_argument('--xkcd', help='plot in XKCD style',
                        default=False, type=bool)
    return parser.parse_args(args)

class Query(object):

    def __init__(self, options):
        self.client = http.client.HTTPConnection(options.host, options.port,
                                                 timeout=30)
        if options.debug:
            self.client.set_debuglevel(1)
        self.id = '/' + options.image

    def get_response(self, url):
        response = None
        begin = None
        duration = None
        try:
            self.client.request('GET', url)
            begin = time.time()
            response = self.client.getresponse()
            duration = time.time() - begin
        except http.client.HTTPException as excp:
            LOGGER.exception('HTTPException', exc_info=excp)
            return None
        if response.status != 200:
            LOGGER.warn('STATUS: ' + response.status)
            return None
        return {
            'elapsed': duration,
            'content': response.read()
        }

    def get_original_size(self):
        response = self.get_response(self.id)
        if response is None:
            return None
        vfile = io.BytesIO(response['content'])
        vimage = Image.open(vfile)
        vfile.close()
        return vimage.size

    def get_resized_time(self, width, height):
        url = '{0}?width={1}&height={2}'.format(self.id, width, height)
        response = self.get_response(url)
        if response is None:
            return None
        return response['elapsed']


def paint(raw, plotdata, xkcd=False):
    # pixels -> time
    xsrc = sorted(plotdata.keys())
    N = len(xsrc)
    xdata = NP.zeros((N, 1))
    ydata = NP.zeros((N, 1))
    for i, x in enumerate(xsrc):
        xdata[i, 0] = x
        ydata[i, 0] = statistics.mean(plotdata[x])

    # Ordered lists are nice ;)
    min_pixels = xsrc[0]
    max_pixels = xsrc[-1]
    min_time = raw[0]
    max_time = raw[-1]
    rect = [min_pixels, min_time, max_pixels - min_time, max_time - min_time]

    # Clear
    PLT.cla()
    PLT.clf()
    if xkcd:
        PLT.xkcd()
    #PLT.axes(rect, axisbg='w', frameon=True)
    PLT.xlabel('pixels')
    PLT.ylabel('seconds')
    PLT.grid(True, which='major', axis='both', linestyle='--')

    # Errors
    yerr = NP.zeros((2, N))
    for i in range(N):
        x, y = xdata[i, 0], ydata[i, 0]
        ys = plotdata[x]
        devi = abs(statistics.stdev(ys) - y) if len(ys) > 1 else 0.0
        yerr[0, i] = devi
        yerr[1, i] = devi
    PLT.errorbar(xdata, ydata, yerr)
    PLT.plot(xdata, ydata, 'r-')
    PLT.axis('auto')
    PLT.show()

def main(options):
    LOGGER.info('hello LOGGER')
    logging.info('hello logging')
    query = Query(options)
    size = query.get_original_size()
    if size is None:
        return 1
    LOGGER.info('Original: {0}x{1}'.format(*size))
    width, height = size
    timings = []
    tplot = {}
    for h in range(1, height, options.step):
        LOGGER.info('Query range with height={0}'.format(h))
        for w in range(1, width, options.step):
            elapsed = query.get_resized_time(w, h)
            if elapsed is None:
                continue
            timings.append(elapsed)
            pixels = w * h
            if pixels not in tplot:
                tplot[pixels] = []
            tplot[pixels].append(elapsed)
    count = len(timings)
    ok_set = sorted(list(filter(lambda x: x is not None, timings)))
    count_ok = len(ok_set)
    print('Query count: {0}'.format(count))
    print('Successful transfers: {0}'.format(count_ok))
    if count_ok < 1:
        LOGGER.error('Can not produce statistics because of too many failures')
        return 1

    mintime = min(ok_set)
    maxtime = max(ok_set)
    mean = statistics.mean(ok_set)
    median = statistics.median(ok_set)
    total = sum(ok_set)
    print('min: {0} s'.format(mintime))
    print('max: {0} s'.format(maxtime))
    print('mean: {0} s'.format(mean))
    print('median: {0} s'.format(median))
    print('sum: {0} s'.format(total))
    if count_ok < 2:
        return 1

    deviation = statistics.stdev(ok_set)
    print('standard-deviation: {0} s'.format(deviation))
    if PLT is not None:
        paint(ok_set, tplot)
    return 0
    

if __name__ == '__main__':
    options = get_opts(sys.argv[1:])
    sys.exit(main(options))
