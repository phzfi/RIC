'''
Created on 5 Nov 2015

@author: Lauri, Joona
'''
import unittest
import text_to_csv


class TestGenerator(unittest.TestCase):


    def test_read(self):
        data = text_to_csv.read_text("testtext.txt")
        self.assertEqual(len(data), 10, "data was not correct length")
        time_sum = 0
        for line in data:
            time_sum += line[0]
        self.assertEqual(27.11, time_sum, "transfer time was wrong" )


if __name__ == "__main__":
    #import sys;sys.argv = ['', 'Test.testName']
    unittest.main()