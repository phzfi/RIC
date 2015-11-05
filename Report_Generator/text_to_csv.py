'''
Created on 5 Nov 2015

@author: Lauri, Joona
'''
import csv
import sys
import os

#TODO: Change directories and filenames after we have chosen the directory
#structure and the naming policy for the log files.



#The test text logs are saved in the text_directory
#The created csv files are saved into the csv directory
#All tests are saved into one log, here log.txt
#The main function goes throught the text log directory
#And checks that the test is found in the logs.txt file.
#If it is, it will be converted to a csv format and saved
#in the csv_directory.
#All test text logs are deleted. The log.txt file is untouched.
def main():    
    text_directory = "textLogs" + os.sep
    csv_directory = "csvFiles" + os.sep
    logs = read_log("log.txt")
    for filename in os.listdir(text_directory):
        filename = filename[:-4]
        if filename in logs:    
            data = read_text(text_directory + filename + ".txt")
            save_to_csv(csv_directory + filename + ".csv", data, logs[filename])
        os.remove(text_directory + filename + ".txt")
    
#Reads the log files and returns a map of the log where the timestamp is a key
def read_log(path):
    timestamp_map = {}
    try: 
        log = open(path, "r")
        for line in log:
            splitted = line.split(",")
            if splitted[0].strip(" ") != "Date & Time":
                #Windows doesn't allow ":" in filenames, so for testing purposes they
                #are replaced with ".".
                stripped = splitted[0].replace(':',".").strip(" ")
                timestamp_map[stripped] = splitted 
        log.close()    
    except OSError:
        print("OSError while reading log at:", path)    
    return timestamp_map
    
#The following line is an example line from the text Document
#HTTP/1.1 200   0.04 secs:  653743 bytes ==> GET  /images/383504_9b66b4a1f2_o.jpg
#Returns a list which contains lists which contain round trip time taken,
#bytes transferred and get Path
def read_text(path):
    values = []
    try:      
        textfile = open(path, "r")
        for line in textfile:
            try:
                splitted = line.split("  ")
                transfer_data = splitted[1].strip()
                splitted_transfer = transfer_data.split(" ")
                time = float(splitted_transfer[0])
                byte_data = int(splitted_transfer[2])
                path = splitted[2]
                values.append([time, byte_data, path])
            except ValueError:
                print("There was a ValueError in the following line:")
                print(line)
        textfile.close()
    except OSError:
        print("Problem reading text file: ", path)
    return values
    
#saves the data obtained from read_text and read_log into a csv file
def save_to_csv(path, data, end_data):
    try:
        first_line = ["Date & Time","Trans","Elap Time", "Data Trans", "Resp Time",
                       "Trans Rate", "Throughput", "Concurrent", "OKAY", "Failed"]
        file = open(path, 'w')
        writer = csv.writer(file, dialect = "excel", lineterminator='\n')
        writer.writerow(first_line)
        writer.writerow(end_data)
        writer.writerow([""])
        writer.writerow([""])
        writer.writerow(["Roundtrip time", "Bytes", "Path"])
        for line in data:
            writer.writerow(line)
            
        file.close()
    except OSError:
        print("OSError while saving to csv file from file:",path)
        sys.exit
    except:
        print("Other error occured while saving to csv:", path)
        sys.exit()
main()