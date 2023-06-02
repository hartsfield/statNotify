#!/bin/sh
output_file_name=$1
output_dir_name=".${output_file_name}_upload"
mkdir $output_dir_name
go build -o $output_file_name
mv output_file_name $output_dir_name
cp Dockerfile $output_dir_name

scp -r $output_dir_name hrtsfld@173.230.142.131:~
ssh -t root@173.230.142.131 \
    -- "cd /home/hrtsfld/${output_dir_name} \
    && docker build . --network=container:purple_panties --tag=fresh_panties \
    && docker run -ti --network=container:purple_panties --rm fresh_panties"
