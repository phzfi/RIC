"""
Created on 5 Nov 2015

@author: Lauri, Joona

"""
import csv
import sys
import os
import logging
import codecs
import traceback

"""
The testResults are in the from_dir.
The created csv files are saved into the to_dir.
The main function goes throught test results and
saves them as an edited csv file better used for
data analysis.

"""

logging.basicConfig(filename='log/error.log',
                    level=logging.ERROR,
                    format='%(asctime)s: %(levelname)s: %(message)s')

def main():
    if len(sys.argv) != 3:
        logging.critical('Wrong number of arguments. Usage: "python csv_formatter.py RAW_FILE OUT_FILE"')
        sys.exit(1)
    raw_file = sys.argv[1]
    out_file = sys.argv[2]

    data = read_csv(raw_file)
    if len(data) != 0:
        save_csv(out_file, data)
    else:
        logging.error('Either empty file or unable to open: ' + raw_file)


def read_csv(from_path):
    """
    Reads the data from a test result csv file.
    Returns the data as a list containing lists.

    """

    data = []

    try:
        with codecs.open(from_path, 'r', 'utf-8') as inp:
            reader = csv.reader(inp, dialect = 'excel', lineterminator='\n')

            for l in reader:
                data.append(l)

            data[0] = [i.strip() for i in data[0]]
            data[1] = [i.strip() for i in data[1]]

            data.insert(2, [''])
            data.insert(3, [''])
            data.insert(4, ['Round trip time', 'Size', 'Path', 'Response'])

            for i in range(5, len(data)):
                orig = data[i]
                data[i] = [float(orig[2]), int(orig[3]), orig[4], int(orig[1])]
            return data

    except OSError as o:
        logging.critical(str(o) + ' while reading csv file from: ' + from_path)

    except Exception:
        logging.critical(traceback.format_exc() + 
                         ' exception while reading csv from: ' + from_path)
    return []

def save_csv(to_path, data):
    """
    Uses the data given and saves the data into a csv file in the given path.

    """

    try:
        with codecs.open(to_path, 'w', 'utf-8') as output:
            writer = csv.writer(output, dialect = 'excel', lineterminator='\n')
            writer.writerows(data)

    except OSError as o:
        logging.critical(str(o) + ' while saving to csv file from file: ' + to_path)
        sys.exit()

if __name__ == '__main__':
    main()
