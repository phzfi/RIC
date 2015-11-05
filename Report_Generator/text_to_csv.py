'''
Created on 5 Nov 2015

@author: Lauri, Joona
'''

def main():
    "lol"
    
#The following line is an example line from the text Document
#HTTP/1.1 200   0.04 secs:  653743 bytes ==> GET  /images/383504_9b66b4a1f2_o.jpg
#Returns a list which contains lists which contain round trip time taken,
#bytes transferred and get Path
def read_text(filename):
    list = []
    try:
        
        textfile = open(filename, "r")
        for line in textfile:
            try:
                splitted = line.split("  ")
                time = float(splitted[1].split(" ")[0])
                bytes = int(splitted[2].split(" ")[0])
                path = splitted[3]
                list.append([time, bytes, path])
            except ValueError:
                print(line)
    except OSError:
        print("Problem reading text file: ", filename)
    return list
    
def save_to_csv(filename):
    "something"
    
main()