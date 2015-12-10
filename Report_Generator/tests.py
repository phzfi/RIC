"""
Created on 5 Nov 2015

@author: Lauri, Joona

"""
import unittest
import csv_formatter

class TestGenerator(unittest.TestCase):


    def test_read_csv_correct(self):

        data = csv_formatter.read_csv('testfiles/correct.csv')

        correct_data = [['Date & Time',  'Trans',  'Elap Time',  'Data Trans',
                         'Resp Time', 'Trans Rate',  'Throughput',  'Concurrent',
                         'OKAY', 'Failed'],
                        ['10/12/2015 20:00','49','39.02','26','0.1','1.26','0.67',
                         '0.13','49','0'],
                        [''],[''],
                        ['Round trip time', 'Size', 'Path', 'Response'],
                        [0.1, 1, 'a', 200], [0.2, 2, 'b', 200],
                        [0.3, 3, 'c', 200], [0.4, 4, 'd', 200],
                        [0.5, 5, 'e', 200], [0.6, 6, 'f', 200],
                        [0.7, 7, 'g', 200], [0.8, 8, 'h', 200]]

        self.assertEqual(len(data), 13, 'data was not correct length')

        for i in range(13):
            self.assertEqual(data[i], correct_data[i])

    def test_read_csv_incorrect(self):

        data = csv_formatter.read_csv('testfiles/incorrect.csv')

        self.assertEqual(len(data), 0, 'data was not empty')

if __name__ == "__main__":
    #import sys;sys.argv = ['', 'Test.testName']
    unittest.main()
