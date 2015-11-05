'''
Created on 5 Nov 2015

@author: Lauri, Joona
'''
import csv
import sys
import os

def main():
    
    text_directory = "textLogs"+os.sep
    csv_directory = "csvFiles"+os.sep
    for filename in os.listdir(text_directory):
            data = read_text(text_directory+filename)
            save_to_csv(csv_directory+filename[:-4]+".csv", data)
            os.remove(text_directory+filename)
    

    
    
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
    
#saves the data obtained from read_text into a csv file
def save_to_csv(path, data):
    try:
        file = open(path, 'w')
        writer = csv.writer(file, dialect = "excel", lineterminator='\n')
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