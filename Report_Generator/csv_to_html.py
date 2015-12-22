'''
Created on 18 Dec 2015

@author: Lauri
'''
import csv
import codecs
import traceback
import sys




HIGHLIGHT_HIGHER = (1,3,5,6,7,8)
HIGHLIGHT_LOWER = (2,4,9)

def main():
    args = sys.argv
    if len(args) < 2:
        print("U DOIN IT WRONG SON")
    data = []
    for csv in args[1:]:
        data.append(read_csv_row(csv, 1))
    titles = read_csv_row(args[1], 0)
    html = buildHTML(data, titles)
    save_to_html(html, 'pythonTest.html')


def is_neutral(i):
    return not (i in HIGHLIGHT_HIGHER or i in HIGHLIGHT_LOWER)

def save_to_html(html, to_path):
    try:
        with codecs.open(to_path, 'w', 'utf-8') as output:
            output.write(html)
        
    except:
        print("sth")

def build_row(row_data, title, i):
    column_html = '<td>%s</td>\n' % title
    if not is_neutral(i) and all(row_data[0] == data for data in row_data):
        for d in row_data:
            column_html += ('<td class="even">%s</td>\n' % d)
    elif not is_neutral(i):
        val = max(row_data)
        if i in HIGHLIGHT_LOWER:
            val = min(row_data)
        for d in row_data:
            if d == val:
                column_html += ('<td class="best">%s</td>\n' % d)
            else:
                column_html += ('<td class="neutral">%s</td>\n' % d)
    else:
        for d in row_data:
            column_html += ('<td>%s</td>\n' % d)
    return ('<tr>%s</tr>\n' % column_html) 

            


    
def buildHTML(data, titles):
    html_table = '''<table>\n
                  <tr>\n
                  <th></th>\n
                  <th>RIC</th>\n
                  <th>Tumbor</th>\n
                  </tr>\n'''
    for i in range(len(titles)):
        row_data = []
        for d in data:
            row_data.append(d[i])
        html_table += build_row(row_data, titles[i], i)
    
    html =   '''<!DOCTYPE html>\n
                <html>\n
                <head>\n
                <link rel="stylesheet" 
                type="text/css" 
                href="tablestyle.css">
                </head>\n
                <body>\n
                %s
                </table>
                </body>\n
                </html> ''' % html_table
    return html

def read_csv_row(from_path, row_number):
    try:
        with codecs.open(from_path, 'r', 'utf-8') as inp:
            reader = csv.reader(inp, dialect = 'excel', lineterminator='\n')
            i = 0
            for l in reader:
                if i == row_number and i != 0:
                    data = [x.strip() for x in l]
                    data = [l[0],int(l[1]),
                            float(l[2]),int(l[3]),
                            float(l[4]),float(l[5]),
                            float(l[6]),float(l[7]),
                            int(l[8]),int(l[9])]
                    return data 
                elif i == row_number:
                    return l
                i += 1
    
    except Exception:
        print(traceback.format_exc())
        
main()