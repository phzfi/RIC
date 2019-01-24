HOW TO USE:

Quick test:
1. Run run_siege_tests.sh by...
sh run_siege_tests.sh SEED_VALUE CONCURRENCY

2. This updates siegeresults.html automatically

OR

1. Run the different run_siege_xxx_xxx.sh files you want
2. Resulting csv files go into results folder
3. Update the siegeresults.html table with the tests you want with...

python csv_to_html.py results/<csvFileName1> results/<csvFileName2>
results/<csvFileName3> ...
