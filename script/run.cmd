cd ..

set NUMWORKER=3
set NUMPOINT=100
set NUMCLUSTER=30

echo %NUMWORKER%
echo %NUMPOINT%
echo %NUMCLUSTER%

CHOICE /C 123 /M "Select [1]: LLyod, [2]: standard kmeans, [3]: keans plus plus"
set choice=%errorlevel%

if %choice%==1 (
 set ALGO=llyod
)
if %choice%==2 (
 set ALGO=standardKMeans
)
if %choice%==3 (
 set ALGO=kmeansAlgo
)

echo NUMWORKER=%NUMWORKER% > .env
echo NUMPOINT=%NUMPOINT%>> .env
echo NUMCLUSTER=%NUMCLUSTER%>> .env
echo ALGO=%ALGO%>> .env

docker-compose build

timeout 10

docker-compose up --scale worker_s=%NUMWORKER%