import sys
import logging
import codecs
import random
import traceback


"""
urls_randomizer is used to randomize the urls.txt and turls.txt with a given
seed.

"""

logging.basicConfig(filename='log/error.log',
                    level=logging.ERROR,
                    format='%(asctime)s: %(levelname)s: %(message)s')


def main():
    if len(sys.argv) != 3:
        logging.critical('Wrong number of arguments.',
                         'Usage: "python urls_randomizer URLS_FILE SEED"')
        sys.exit(1)
    urls_file = sys.argv[1]
    seed = sys.argv[2]
    try:
        seed = int(seed)
        urls = read_urls(urls_file)
        save_random_urls(urls, urls_file+'_temp.txt', seed)
    except ValueError:
        logging.critical('Seed was not an integer')


def read_urls(urls_file):
    data = []
    try:
        with codecs.open(urls_file, 'r', 'utf-8') as inp:
            for line in inp:
                data.append(line)
        return data
    except OSError as o:
        logging.critical(str(o) + ' while reading urls file from: '+urls_file)
        sys.exit(1)


def save_random_urls(data, to_file, seed):
    try:
        with codecs.open(to_file, 'w', 'utf-8') as output:
            output.write(data[0])
            data = data[1:]
            random.seed(seed)
            for line in data:
                output.write(line)
            return
    except OSError as o:
        logging.critical(str(o) + ' while saving urls file to: '+to_file)
    except Exception:
        logging.critical(traceback.format_exc() +
                         ' exception while saving urls file to: '+to_file)
    sys.exit(1)

main()
