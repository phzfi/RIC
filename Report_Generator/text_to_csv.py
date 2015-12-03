"""
Created on 5 Nov 2015

@author: Lauri, Joona

"""
import csv
import sys
import os
import logging
import codecs

"""
TODO: Change directories and filenames after we have chosen the directory
structure and the naming policy for the log files.

The test text logs are saved in the text_directory
The created csv files are saved into the csv directory
All tests are saved into one log, here log.txt
The main function goes throught the text log directory
And checks that the test is found in the logs.txt file.
If it is, it will be converted to a csv format and saved
in the csv_directory.
All test text logs are deleted. The log.txt file is untouched.

"""

logging.basicConfig(filename='Report_Generator_error.log',
                    level=logging.ERROR,
                    format='%(asctime)s: %(levelname)s: %(message)s')

def main():    
    text_directory = 'textLogs' 
    csv_directory = 'csvFiles'
    logs = read_log('siege.log')

    if len(logs) == 0:
        logging.critical('siege.log is empty.')
        sys.exit(1)

    try:
        text_dir_list = os.listdir(text_directory)
    except Exception as e:
        logging.critical(str(e) + ' Critical error while listing text directory')
        sys.exit(1)

    """
    Goes through all files in text directory and attempts to convert them to a csv file. 
    The text file is removed in the process if conversion is successful or if text file
    does not appear in the logs file.

    """

    for filename in text_dir_list:
        basename, _ = os.path.splitext(filename)
        if basename in logs:    
            data = read_text(os.path.join(text_directory, filename))
            save_to_csv(os.path.join(csv_directory, basename) + '.csv', data, logs[basename])
    try:
        os.remove(os.path.join(text_directory, filename))
    except OSError as o:
        logging.error(str(o) + ' while attempting to remove: ' + filename)


def read_log(path):
    """
    Reads the log files and returns a map of the log data where the 
    timestamp is a key.
    Two top lines of an example log file is shown below:

    Date & Time,  Trans,  Elap Time,  Data Trans,  Resp Time,  Trans Rate,  Throughput,  Concurrent,    OKAY,   Failed
    2015-10-29 16:34:29,    400,      24.50,         833,       0.54,       16.33,       34.00,        8.88,     400,       0

    """

    timestamp_map = {}    
    try:
        with codecs.open(path, 'r', 'utf-8') as log:
            for line in log:
                splitted = line.split(',')                    
                """
                Windows doesn't allow ':' in filenames, so 
                they are replaced with '.'.

                """
                timestamp = splitted[0].replace(':','.').strip()
                if timestamp != 'Date & Time':
                    timestamp_map[timestamp] = [x.strip() for x in splitted] 

    except OSError as o:
        logging.critical(str(o) + ' while reading log at: ' + path)
        sys.exit(1)
    except Exception as e:
        logging.critical(str(e) + ' while reading log at: ' + path)
        sys.exit(1)
    return timestamp_map



def read_text(path):
    """    
    The following line is an example line from the text Document:

    HTTP/1.1 200   2.66 secs: 7841615 bytes ==> GET  /images/14857596084_b0105b8e88_o.jpg

    Returns a list which contains lists which contain round trip time taken,
    bytes transferred, GET path and server response.

    """
    values = []
    try:      
        with codecs.open(path, 'r', 'utf-8') as textfile: 
            for line in textfile:
                try:
                    server_response, transfer_data, image_path = line.split('  ')
                    transfer_data = transfer_data.strip()
                    roundtrip_time, _, byte_data, _, _, _ = transfer_data.split(' ')
                    values.append([float(roundtrip_time), int(byte_data), 
                                   image_path.strip(), server_response.strip()])
                except ValueError as v:
                    logging.error(str(v) + ' in the following line: ' + line)

    except OSError as o:
        logging.error(str(o) + ' while reading text at: ' + path)
    except Exception as e:
        logging.critical(str(e) + ' while reading text at: ' + path)
        sys.exit(1)
    return values


def save_to_csv(path, data, summary_data):
    """ Saves the data obtained from read_text and read_log into a csv file. """
    try:
        with codecs.open(path, 'w', 'utf-8') as output:
            writer = csv.writer(output, dialect = 'excel', lineterminator='\n')

            writer.writerow(['Date & Time','Trans','Elap Time', 'Data Trans',
                             'Resp Time','Trans Rate', 'Throughput',
                             'Concurrent', 'OKAY', 'Failed'])
            writer.writerow(summary_data)
            writer.writerow([''])
            writer.writerow([''])
            writer.writerow(['Roundtrip time', 'Bytes', 'Path', 'Response'])
            for line in data:
                writer.writerow(line)
            return
    except OSError as o:
        logging.critical(str(o) + ' while saving to csv file from file: ' + path)

    except Exception as e:
        logging.critical(str(e) + ' exception while saving to csv: ' + path)
    sys.exit(1)

if __name__ == '__main__':
    main()
