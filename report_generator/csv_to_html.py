"""
Created on 18 Dec 2015

@author: Lauri
@modified-by: Kristian

"""
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


HIGHLIGHT_HIGHER = (1, 3, 5, 6, 8)
HIGHLIGHT_LOWER = (2, 4, 9)
HTML_DOC = '''
<!DOCTYPE html>
<html>
<head>
<title>Results</title>
<link rel="stylesheet" type="text/css" href="tablestyle.css">
</head>
<body>
<table>
{0}
</table>
</body>
</html>'''

logging.basicConfig(
    filename='log/error.log',
    level=logging.ERROR,
    format='%(asctime)s: %(levelname)s: %(message)s'
)


def main(args):
    if len(args) < 3:
        logging.critical('Wrong number of arguments.\n' +
                         'Usage: "python csv_to_html.py htmlTable testData*"')
        sys.exit(1)
    software = []
    titles = None
    for csv in args[2:]:
        csv_name = os.path.basename(csv)
        name = csv_name.split("_")[0]
        if titles is None:
            titles = read_csv_row(csv, 0)
        data = read_csv_row(csv, 1)
        software.append({'name': name, 'data': data})
    html = buildHTML(software, titles)
    sys.exit(save_to_html(html, args[1]))


def get_picker(i):
    """
    Checks if the current row in the table should be highlited or not.

    """
    if i in HIGHLIGHT_HIGHER:
        return max
    if i in HIGHLIGHT_LOWER:
        return min
    return None


def save_to_html(html, to_path):
    try:
        with codecs.open(to_path, 'w', 'utf-8') as output:
            output.write(html)
    except Exception:
        msg = '{0} exception while saving html to: {1}'
        logging.critical(msg.format(traceback.format_exc(), to_path))
        return 1
    return 0


def build_row(row_data, picker):
    """
    Returns row html built with the given row_data, title of row and
    row index.

    """
    column = '<td>{0}</td>'
    column_cls = '<td class="{0}">{1}</td>'
    if picker is None:
        return ''.join([column.format(d) for d in row_data])
    if min(row_data) == max(row_data):
        return ''.join([column_cls.format('even', d) for d in row_data])
    val = picker(row_data)
    return ''.join([
        column_cls.format('best' if d == val else 'neutral', d)
        for d in row_data
    ])


def buildHTML(software, titles):
    """
    Builds the html to display the given data in a table.

    """
    row = '<tr>{0}</tr>'
    head = '<th>{0}</th>'
    column = '<td>{0}</td>'
    headers = [head.format('')]
    headers.extend([head.format(s['name']) for s in software])
    html_table = row.format(''.join(headers))

    for i, title in enumerate(titles):
        row_data = [app['data'][i] for app in software]
        picker = get_picker(i)
        row_head = column.format(title)
        row_content = build_row(row_data, picker)
        html_table += row.format(row_head + row_content)

    return HTML_DOC.format(html_table)


def read_csv_row(from_path, row_number=None):
    """
    Reads a row in a csv file formatted with the csv_formatter script.

    """
    forms = [int, float, int, float, float, float, float, int, int]
    try:
        with codecs.open(from_path, 'r', 'utf-8') as inp:
            reader = csv.reader(inp, dialect='excel', lineterminator='\n')
            if row_number is None:
                return [row for row in reader]
            if row_number == 0:
                return next(reader)
            f = lambda x: x[1] if (x[0] == row_number) else None
            iterated = map(f, enumerate(reader))
            row = filter(lambda x: x is not None, iterated)[0]
            return [a(b.strip()) for (a, b) in zip(forms, row)]
    except Exception as excp:
        print(excp)
        msg = '{0} exception while reading csv from: {1}'
        logging.critical(msg.format(traceback.format_exc(), from_path))
        sys.exit(1)

if __name__ == '__main__':
    main(sys.argv)
