import sys
import logging
import codecs
import random
import traceback
from os import path

"""
urls_randomizer is used to randomize the urls.txt and turls.txt with a given
seed.

"""

logging.basicConfig(filename='log/error.log',
                    level=logging.ERROR,
                    format='%(asctime)s: %(levelname)s: %(message)s')


def main():
    if len(sys.argv) != 3:
        mess = ('Wrong number of arguments.\n' +
                'Usage: "python urls_randomizer URLS_FILE SEED')
        logging.critical(mess)
        print(mess)
        sys.exit(1)
    urls_file = sys.argv[1]
    temp_file = '{:s}_temp.txt'.format(path.splitext(urls_file)[0])
    seed = sys.argv[2]
    try:
        seed = int(seed)
    except ValueError:
        logging.critical('Seed was not an integer')
    urls = read_urls(urls_file)
    save_random_urls(urls, temp_file, seed)


def read_urls(urls_file):
    try:
        with codecs.open(urls_file, 'r', 'utf-8') as inp:
            return inp.readlines()
    except OSError as err:
        logging.critical(str(err) + ' while reading urls file from: ' + urls_file)
        sys.exit(1)


def save_random_urls(data, to_file, seed):
    try:
        with codecs.open(to_file, 'w', 'utf-8') as output:
            output.write(data[0])
            data = data[1:]
            random.seed(seed)
            random.shuffle(data)
            output.writelines(data)
            return
    except OSError as err:
        logging.critical(str(err) + ' while saving urls file to: '+to_file)
    except Exception:
        logging.critical(traceback.format_exc() +
                         ' exception while saving urls file to: '+to_file)
    sys.exit(1)

if __name__ == '__main__':
    main()
