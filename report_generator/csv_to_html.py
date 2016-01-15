'''
Created on 18 Dec 2015

@author: Lauri
'''
import csv
import codecs
import traceback
import sys
import logging
import os

"""
Different formatted csvs made with the csv_formatter are used to
make html tables. The different formatted csv paths are given as parameters.

"""


HIGHLIGHT_HIGHER = (1, 3, 5, 6, 7, 8)
HIGHLIGHT_LOWER = (2, 4, 9)

logging.basicConfig(filename='log/error.log',
                    level=logging.ERROR,
                    format='%(asctime)s: %(levelname)s: %(message)s')


def main():
    args = sys.argv
    if len(args) < 2:
        logging.critical('Wrong number of arguments.' +
                         'Usage: "python csv_to_html.py testData*"')
        sys.exit()
    data = []
    software = []
    for csv in args[1:]:
        csv_name = os.path.basename(csv)
        software.append(csv_name.split("_")[0])
        data.append(read_csv_row(csv, 1))
    titles = read_csv_row(args[1], 0)
    html = buildHTML(software, data, titles)
    save_to_html(html, 'siegeresults.html')


def is_neutral(i):
    """
    Checks if the current row in the table should be highlited or not.

    """

    return not (i in HIGHLIGHT_HIGHER or i in HIGHLIGHT_LOWER)


def save_to_html(html, to_path):
    try:
        with codecs.open(to_path, 'w', 'utf-8') as output:
            output.write(html)

    except:
        logging.critical(traceback.format_exc() +
                         ' exception while saving html to: ' +
                         str(to_path))
        sys.exit


def build_row(row_data, title, i):
    """
    Returns row html built with the given row_data, title of row and
    row index.

    """

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


def buildHTML(software, data, titles):
    """
    Builds the html to display the given data in a table.

    """

    html_table = '''<table>\n
                  <tr>\n
                  <th></th>\n'''
    for s in software:
        html_table += '<th>%s</th>\n''' % s
    html_table += '</tr>\n'

    for i in range(len(titles)):
        row_data = []
        for d in data:
            row_data.append(d[i])
        html_table += build_row(row_data, titles[i], i)


    html =   '''<!DOCTYPE html>\n
                <html>\n
                <head>\n
                <title>Results</title>\n
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
    """
    Reads a row in a csv file formatted with the csv_formatter script.

    """

    try:
        with codecs.open(from_path, 'r', 'utf-8') as inp:
            reader = csv.reader(inp, dialect='excel', lineterminator='\n')
            i = 0
            for l in reader:
                if i == row_number and i != 0:
                    data = [x.strip() for x in l]
                    data = [l[0], int(l[1]),
                            float(l[2]), int(l[3]),
                            float(l[4]), float(l[5]),
                            float(l[6]), float(l[7]),
                            int(l[8]), int(l[9])]
                    return data
                elif i == row_number:
                    return l
                i += 1

    except Exception:
        logging.critical(traceback.format_exc() +
                         ' exception while reading csv from: ' +
                          str(from_path))
        sys.exit

main()
