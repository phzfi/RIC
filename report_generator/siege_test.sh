#!/bin/bash

# Checks if user started program properly, if not then starts it
if [ -z $BASH_VERSION ]; then
	exec bash "$0" "$@"
	exit $?
fi

echo "Use the -h flag for help."

SEED=5
CONCURRENT="-c20"
REQUESTS_OR_TIME="-r40"
LOCAL=false
AUTO=false
DELAY="-d1"
RIC=false
THUMBOR=false
CIB=false

while getopts ":RCTlhaus:c:r:d:t:" OPTION; do
  case $OPTION in
		h)
		  printf "\nHelp: Flag usage -x -y -z OR -xyz (if they do not take values)\nValue for x flag is give given as -yz -xVALUE\n"
			printf "Available flags are:\n
-h		Help Text\n
-R		Use RIC Software (default if no software given)\n
-C		Use CIB Software\n
-T		Use Thumbor Software\n
-l		Run tests locally (localhost). Default is off (runs on staging)\n
-a		Run tests automatically with cache clearing and program restarts. Default is off\n
-s 		Takes a seed value. This randomises the order URLS are accessed. Defaults to 5\n
-d		Takes a delay value between 0 - NUM. Each siege simulated user is delayed for a random number of seconds between one and NUM. Setting delay to 0 causes there to be no delay. Default is 1\n
-r		Takes a request per user value. Every concurrent user will do this many requests. DO NOT SET IF USING TIME SETTING. Default is 40\n
-t		Takes a time value in seconds on how long to siege the server. DO NOT SET IF USING REQUESTS PER USER. If both time and requests per user are given, the last one is used. Default is to use requests, not time.\n
-c		Takes a concurrent value. This tells how many concurrent users are using the server. Forexample, if concurrent is set to 40 and requests per user is set to 40, then a total of 1600 requests are made. Default is 20.\n\n"
			exit
			;;
    s)
      echo "The value of seed is $OPTARG"
			SEED=$OPTARG
      ;;
		c)
			echo "The value of concurrent is $OPTARG"
			CONCURRENT="-c"$OPTARG
			;;
		t)
			echo "The value of time is "$OPTARG"s"
			REQUESTS_OR_TIME="-t"$OPTARG"s"
			;;
		l)
			echo "Tests run locally"
			LOCAL=true
			;;
		r)
			echo "The value of requests per user is $OPTARG"
			REQUESTS_PER_USER="-r"$OPTARG
			TIME=false
			;;
		a)
		  echo "Tests run automatically with cache clearing and program restarts"
			AUTO=true
			;;
		d)
			if [ $OPTARG -eq 0 ]; then
				echo "Delay set to 0, benchmarking setting will be used"
				DELAY="-b"
			else
				echo "The value of delay is $OPTARG"
				DELAY="-d"$OPTARG
			fi
			;;
		R)
			echo "RIC will be tested"
			RIC=true
			;;
		T)
			echo "THUMBOR will be tested"
			THUMBOR=true
			;;
		C)
			echo "CIB will be tested"
			CIB=true
			;;
    \?)
      echo "Invalid option: -$OPTARG"
			exit 1
      ;;
		:)
			echo "Option -$OPTARG requires an argument."
			exit 1
    	;;
  esac
done

# cd to report_generator
DIR="$( cd "$(dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"

# Defaults to use RIC if not other software has been given
if [ "$RIC" == false ] && [ "$THUMBOR" == false ] && [ "$CIB" == false ]; then
	echo "No software given, RIC will be used"
	RIC=true
fi


OUTPUT_FILES=()
SETTINGS="s$SEED""$CONCURRENT""$REQUESTS_OR_TIME"-local"$LOCAL"-auto"$AUTO"
SOFT=""

if [ "$RIC" == true ]; then
	SOFT+="ric-"
	if [ "$AUTO" == true ]; then
		RAW_FILE=raw/ric-before_$(date +%Y-%m-%d_%H-%M-%S).txt
		OUT_FILE=results/ric-before_"$SETTINGS"_$(date +%Y-%m-%d_%H-%M-%S).csv
		OUTPUT_FILES+="$OUT_FILE "
		bash ./helper_scripts/siege_runner.sh urls.txt $SEED $REQUESTS_OR_TIME $CONCURRENT $LOCAL $AUTO $DELAY ric $RAW_FILE $OUT_FILE
	fi
	RAW_FILE=raw/ric_$(date +%Y-%m-%d_%H-%M-%S).txt
	OUT_FILE=results/ric_"$SETTINGS"_$(date +%Y-%m-%d_%H-%M-%S).csv
	OUTPUT_FILES+="$OUT_FILE "
	bash ./helper_scripts/siege_runner.sh urls.txt $SEED $REQUESTS_OR_TIME $CONCURRENT $LOCAL false $DELAY ric $RAW_FILE $OUT_FILE
fi

if [ "$CIB" == true ]; then
	SOFT+="cib-"
	if [ "$AUTO" == true ]; then
		RAW_FILE=raw/cib-before_$(date +%Y-%m-%d_%H-%M-%S).txt
		OUT_FILE=results/cib-before_"$SETTINGS"_$(date +%Y-%m-%d_%H-%M-%S).csv
		OUTPUT_FILES+="$OUT_FILE "
		bash ./helper_scripts/siege_runner.sh jurls.txt $SEED $REQUESTS_OR_TIME $CONCURRENT $LOCAL $AUTO $DELAY cib $RAW_FILE $OUT_FILE
	fi
	RAW_FILE=raw/cib_$(date +%Y-%m-%d_%H-%M-%S).txt
	OUT_FILE=results/cib_"$SETTINGS"_$(date +%Y-%m-%d_%H-%M-%S).csv
	OUTPUT_FILES+="$OUT_FILE "
	bash ./helper_scripts/siege_runner.sh jurls.txt $SEED $REQUESTS_OR_TIME $CONCURRENT $LOCAL false $DELAY cib $RAW_FILE $OUT_FILE
fi

if [ "$THUMBOR" == true ]; then
	SOFT+="thumbor-"
	if [ "$AUTO" == true ]; then
		RAW_FILE=raw/thumbor-before_$(date +%Y-%m-%d_%H-%M-%S).txt
		OUT_FILE=results/thumbor-before_"$SETTINGS"_$(date +%Y-%m-%d_%H-%M-%S).csv
		OUTPUT_FILES+="$OUT_FILE "
		bash ./helper_scripts/siege_runner.sh turls.txt $SEED $REQUESTS_OR_TIME $CONCURRENT $LOCAL $AUTO $DELAY thumbor $RAW_FILE $OUT_FILE
	fi
	RAW_FILE=raw/thumbor_$(date +%Y-%m-%d_%H-%M-%S).txt
	OUT_FILE=results/thumbor_"$SETTINGS"_$(date +%Y-%m-%d_%H-%M-%S).csv
	OUTPUT_FILES+="$OUT_FILE "
	bash ./helper_scripts/siege_runner.sh turls.txt $SEED $REQUESTS_OR_TIME $CONCURRENT $LOCAL false $DELAY thumbor $RAW_FILE $OUT_FILE
fi


HTML=$SOFT"$SETTINGS".html
python3 csv_to_html.py html_tables/$HTML $OUTPUT_FILES
echo "Siege Complete and results can be found in file $HTML"


#If auto is true, then reboot RIC
if [ "$AUTO" == true ]; then
	echo "Because AUTO was set to true, RIC is being rebooted"
  if [ "$LOCAL" == true ]; then
      echo "Waiting 10s for ric to boot"
      sh ../scripts/start_ric_stop_rest.sh
      sleep 10s
  else
    sh start_ric_stop_rest.sh
  fi
fi
