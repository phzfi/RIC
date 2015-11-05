'''
Created on 5 Nov 2015

@author: Lauri, Joona
'''
import csv

def main():
    filename = "testtext"
    data = read_text(filename+".txt")
    save_to_csv(filename+".csv", data)
    
#The following line is an example line from the text Document
#HTTP/1.1 200   0.04 secs:  653743 bytes ==> GET  /images/383504_9b66b4a1f2_o.jpg
#Returns a list which contains lists which contain round trip time taken,
#bytes transferred and get Path
def read_text(filename):
    values = []
    try:      
        textfile = open(filename, "r")
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
        print("Problem reading text file: ", filename)
    return values
    
def save_to_csv(filename, data):
    try:
        file = open(filename, 'w')
        writer = csv.writer(file, dialect = "excel")
        for line in data:
            writer.writerow(line)
            
        file.close()
    except OSError:
        print("Problem saving to csv file")
    
    
main()