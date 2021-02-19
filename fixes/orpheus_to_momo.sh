catmandu convert MARC --type XML \
  --fix orpheus_to_momo.fix \
  to --line-delimited 1 \
  | ./split_recs.sh
