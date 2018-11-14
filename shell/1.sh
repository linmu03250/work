filelist=`ls $1`
touch s.log
for file in $filelist
do
	echo $file
	#sed -i "" 's/mtz1.mt-meipai.//g' $file
done
