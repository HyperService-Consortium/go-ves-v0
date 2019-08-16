#!/bin/bash

host=47.251.2.73:26657
out_file=data.csv
test=s-bench

for i in 1 5 50 100 120 140 160 180 200 210 ; do
    let session_limit=$i
    for j in 100 200 400 800 1000 1500 ; do
    let signature_size=$j
        for k in 1 2 4 8 16 32 ; do
        let action_long=$k
            for isbatch in true false ; do
                echo `./$test -host $host -o $out_file -batch=$isbatch -ses $session_limit -con $signature_size -accs $action_long`
            done
        done
    done
done

# echo $case_count, $predict_count