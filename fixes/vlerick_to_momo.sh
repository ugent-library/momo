catmandu convert \
  --fix vlerick_to_momo.fix \
  to --line-delimited 1 \
  | ./split_recs.sh
