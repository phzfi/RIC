"""
Created on 5 Nov 2015

@author: Lauri, Joona

"""
import unittest
import text_to_csv

class TestGenerator(unittest.TestCase):


    def test_read_text_correct(self):

        data = text_to_csv.read_text('testfiles/testtextcorrect.txt')

        self.assertEqual(len(data), 10, 'data was not correct length')

        for i in range(9):
            self.assertEqual(data[i], [float(i), i * 10, 
                        '/images/14857596084_b0105b8e88_o.jpg', 'HTTP/1.1 200'],
                        'wrong line value ' + str(i))

    def test_read_text_error(self):

        data = text_to_csv.read_text('testfiles/testtexterror.txt')

        self.assertEqual(len(data), 10, 'data was not correct length')

        for i in range(9):
            self.assertEqual(data[i], [float(i), i * 10,
                            '/images/14857596084_b0105b8e88_o.jpg', 'HTTP/1.1 200'],
                            'wrong line value ' + str(i))

    def test_read_text_empty(self):

        data = text_to_csv.read_text('testfiles/testtextempty.txt')

        self.assertEqual(len(data), 0, 'data was not empty')


    def test_read_log(self):

        data = text_to_csv.read_log('testfiles/testlog.log')

        self.assertEqual(len(data), 3, 'data was not correct length')

        self.assertEqual(data['2015-10-29 16.34.29'],
                         ['2015-10-29 16:34:29', '400', '24.50', '833',
                          '0.54', '16.33', '34.00', '8.88', '400', '0'],
                          'wrong data value 1')
        self.assertEqual(data['2015-10-29 16.35.45'],
                        ['2015-10-29 16:35:45', '400', '25.39', '833', 
                         '0.60', '15.75', '32.81', '9.40', '400', '0'],
                         'wrong data value 2')
        self.assertEqual(data['2015-10-29 16.53.23'], 
                        ['2015-10-29 16:53:23', '400', '23.98', '833',
                         '0.45', '16.68', '34.74', '7.55', '400', '0'], 
                         'wrong data value 3')

    def test_read_log_empty(self):

        data = text_to_csv.read_log('testfiles/testlogempty.log')

        self.assertEqual(len(data), 0, 'data was not correct length')


if __name__ == "__main__":
    #import sys;sys.argv = ['', 'Test.testName']
    unittest.main()
